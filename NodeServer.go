package main

import (
	"fmt"
	"strings"
	"sync"
)

// NodeServer is the server of the node and it is responsible for communication between nodes
type NodeServer struct {
	node         *Node
	peers        []string
	mutex        *sync.Mutex
	communicator *Communicator
	webServer    *WebServer
	recvChannel  chan *Packet
	sendChannel  chan *Packet
}

const (

	// ListenPort is the port on which the Server is listening to
	ListenPort = 4416
)

func (n *NodeServer) init(node *Node, config *JSONConfig) {
	n.node = node
	n.mutex = &sync.Mutex{}
	n.peers = []string{}
	peerStr := config.Peers
	splat := strings.Split(peerStr, ";")
	for i := 0; i < len(splat); i++ {
		splat[i] = strings.TrimSpace(splat[i])
	}
	if len(splat) > 0 {
		if splat[0] != "" {
			n.peers = append(node.server.peers, splat...)
		}
	}
	n.webServer = &WebServer{server: n}
	n.recvChannel = make(chan *Packet)
	n.sendChannel = make(chan *Packet)
	n.communicator = NewCommunicator(n, config.Addr, n.recvChannel, n.sendChannel, ListenPort)
	go n.communicator.Listen()
	go n.webServer.Start()
	go n.handlePackets()
}

func (n *NodeServer) handlePackets() {
	for {
		p := <-n.recvChannel
		retP := &Packet{requestType: ""}
		switch p.Type() {
		case TPR:
			retP = NewPacket(STPM, n.node.transactionPool.FormatSTPM())
		case BR:
			retP = NewPacket(SCM, FormatSCM(n.node.blockchain.GetLatestIndex(), n.node.blockchain.GetLatestHash()))
		case PR:
			retP = NewPacket(PA, FormatPA(n.peers))
		case FT:
			if !n.node.blockchain.IsUpdating() {
				num, err := UnformatFT(p.data)
				if err != nil {
					break
				}
				retP = NewPacket(BP, FormatBP(n.node.blockchain.GetBlocksFromTop(num)))
			}
		case IS:
			if !n.node.blockchain.IsUpdating() {
				index, err := UnformatIS(p.data)
				if err != nil {
					break
				}
				retP = NewPacket(BP, FormatBP(n.node.blockchain.GetBlocksFromIndex(index)))
			}
		default:
		}
		n.sendChannel <- retP
	}
}

// Address returns the address of the node
func (n *NodeServer) Address() string {
	return n.communicator.Address()
}

func (n *NodeServer) doesPeerExist(peer string) bool {
	for _, addr := range n.peers {
		if peer == addr {
			return true
		}
	}
	return false
}

func (n *NodeServer) savePeers(config *JSONConfig, peers []string) {
	for _, peer := range peers {
		n.savePeer(config, peer)
	}
}

func (n *NodeServer) savePeer(config *JSONConfig, peer string) {
	confPeers := config.Peers
	splat := strings.Split(confPeers, ";")
	for _, confPeer := range splat {
		if confPeer == peer {
			return
		}
	}
	config.Peers += fmt.Sprintf(";%s", peer)
}

func (n *NodeServer) addPeers(peers []string) {
	for _, peer := range peers {
		n.addPeer(peer)
	}
}

func (n *NodeServer) addPeer(peer string) {
	peerSplat := strings.Split(peer, ":")
	if len(peerSplat) == 2 {
		peer = peerSplat[0]
		if !n.doesPeerExist(peer) && peer != n.Address() {
			n.mutex.Lock()
			n.peers = append(n.peers, peer)
			fmt.Printf("%s is a new peer\n", peer)
			n.mutex.Unlock()
		}
	}
}

func (n *NodeServer) requestBlockchain() {
	for _, peer := range n.peers {
		p := NewPacket(BR, []byte{})
		p, err := n.communicator.SR1(peer, p)
		if err != nil {
			continue
		}
		if p.Type() != SCM {
			continue
		}
		index, _, err := UnformatSCM(p.data)
		if err != nil {
			continue
		}
		if index <= n.node.blockchain.GetLatestIndex() {
			continue
		}
		p = NewPacket(FT, FormatFT(index-n.node.blockchain.GetLatestIndex()))
		p, err = n.communicator.SR1(peer, p)
		if err != nil {
			continue
		}
		if p.Type() != BP {
			continue
		}
		blocks, err := UnformatBP(p.data)
		if err != nil {
			continue
		}
		allBlocks := blocks
		for !n.node.blockchain.CompareBlockchains(blocks) {
			p = NewPacket(IS, FormatIS(blocks[0].index))
			p, err = n.communicator.SR1(peer, p)
			if err != nil {
				continue
			}
			if p.Type() != BP {
				continue
			}
			blocks, err = UnformatBP(p.data)
			if err != nil {
				continue
			}
			allBlocks = append(blocks, allBlocks...)
		}
		if !n.node.blockchain.IsUpdating() {
			n.node.blockchain.SetUpdating(true)
			n.node.blockchain.ReplaceBlocks(allBlocks)
			n.node.blockchain.SetUpdating(false)
			fmt.Printf("Updated blockchain from %s\n", peer)
			n.node.PrintBlockchain()
		}
	}
}

func (n *NodeServer) requestPeers() {
	for _, peer := range n.peers {
		p := NewPacket(PR, []byte{})
		p, err := n.communicator.SR1(peer, p)
		if err != nil {
			continue
		}
		if p.Type() != PA {
			continue
		}
		addrs := UnformatPA(p.data)
		n.addPeers(addrs)
	}
}

func (n *NodeServer) requestPool() {
	for _, peer := range n.peers {
		p := NewPacket(TPR, []byte{})
		p, err := n.communicator.SR1(peer, p)
		if err != nil {
			continue
		}
		if p.Type() != STPM {
			continue
		}
		trans, err := UnformatSTPM(p.data)
		if err != nil {
			continue
		}
		for _, t := range trans {
			if !n.node.transactionPool.DoesExists(t) {
				n.node.transactionPool.addTransaction(t)
				fmt.Printf("Added a new transaction from %s\n", peer)
			}
		}
	}
}
