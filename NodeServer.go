package main

import (
	"crypto/ecdsa"
	"fmt"
	"sync"
)

// NodeServer is the server of the node and it is responsible for communication between nodes
type NodeServer struct {
	node         *Node
	peers        []string
	mutex        *sync.Mutex
	communicator *Communicator
	recvChannel  chan *Packet
	sendChannel  chan *Packet
}

const (

	// ListenPort is the IP on which the Server is listening to
	ListenPort = 1625
)

func (n *NodeServer) init(node *Node, address string, recvChannel, sendChannel, stmpChannel chan *Packet, privKey *ecdsa.PrivateKey) {
}

func (n *NodeServer) firstInit(conf *JSONConfig, node *Node, privKey *ecdsa.PrivateKey) {
	n.node = node
	n.mutex = &sync.Mutex{}
	n.peers = make([]string, 0)
	n.recvChannel = make(chan *Packet)
	n.sendChannel = make(chan *Packet)
	n.communicator = NewCommunicator(conf.Addr, n.recvChannel, n.sendChannel)
	go n.communicator.Listen()
	//go n.sendToPeers()
}

func (n *NodeServer) handlePackets() { // NOT FISNISHED
	for {
		p := <-n.recvChannel
		retP := &Packet{requestType: ""}
		switch p.Type() {
		case BR:
			retP = NewPacket(SCM, FormatSCM(n.node.blockchain.GetLatestIndex(), n.node.blockchain.GetLatestHash()))
		case TPR:
			//------------------------------
		case PR:
			retP = NewPacket(PA, FormatPA(n.peers))
		case FT:
			num, err := UnformatFT(p.data)
			if err != nil {
				fmt.Println(err)
				break
			}
			retP = NewPacket(BP, FormatBP(n.node.blockchain.GetBlocksFromTop(num)))
		case IS:
			index, err := UnformatIS(p.data)
			if err != nil {
				fmt.Println(err)
				break
			}
			retP = NewPacket(BP, FormatBP(n.node.blockchain.GetBlocksFromIndex(index)))
		case NT:
			//---
		default:
			fmt.Println(ErrPacketType)
		}
		n.sendChannel <- retP
	}
}

// Address returns the address of the node
func (n *NodeServer) Address() string {
	return n.communicator.Address()
}

func (n *NodeServer) savePeers() error {
	// save to database
	return nil
}

func (n *NodeServer) doesPeerExist(peer string) bool {
	for _, addr := range n.peers {
		if peer == addr {
			return true
		}
	}
	return false
}

func (n *NodeServer) requestBlockchain() { /// TEST!!!
	for _, peer := range n.peers {
		p := NewPacket(BR, []byte{})
		p, err := n.communicator.SR1(peer, p)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if p.Type() != SCM {
			fmt.Println(err)
			continue
		}
		index, _, err := UnformatSCM(p.data)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if index <= n.node.blockchain.GetLatestIndex() {
			continue
		}
		p = NewPacket(FT, FormatFT(index-n.node.blockchain.GetLatestIndex()))
		p, err = n.communicator.SR1(peer, p)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if p.Type() != BP {
			continue
		}
		blocks, err := UnformatBP(p.data)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if n.node.blockchain.CompareBlockchains(blocks) {
			n.node.blockchain.AddBlocks(blocks)
			continue
		}
		for !n.node.blockchain.CompareBlockchains(blocks) {
			p = NewPacket(IS, FormatIS(blocks[0].index))
			p, err = n.communicator.SR1(peer, p)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if p.Type() != BP {
				continue
			}
			blocks, err = UnformatBP(p.data)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
}

func (n *NodeServer) requestPeers() {
	for _, peer := range n.peers {
		p := NewPacket(PR, []byte{})
		p, err := n.communicator.SR1(peer, p)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if p.Type() != PA {
			fmt.Println(ErrPacketType)
			continue
		}
		addrs := UnformatPA(p.data)
		addrsToAdd := []string{}
		for _, addr := range addrs {
			if !n.doesPeerExist(addr) {
				addrsToAdd = append(addrsToAdd, addr)
			}
		}
		n.mutex.Lock()
		n.peers = append(n.peers, addrsToAdd...)
		n.mutex.Unlock()
	}
}

func (n *NodeServer) requestPool() {
	for _, peer := range n.peers {
		p := NewPacket(BR, []byte{})
		p, err := n.communicator.SR1(peer, p)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if p.Type() != STPM {
			fmt.Println(ErrPacketType)
			continue
		}
		trans, err := UnformatSTPM(p.data)
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, t := range trans {
			if !n.node.transactionPool.DoesExists(t) {
				n.node.transactionPool.addTransaction(t)
			}
		}
	}
}
