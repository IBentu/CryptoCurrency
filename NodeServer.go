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

// Address returns the address of the node
func (n *NodeServer) Address() string {
	return n.communicator.Address()
}

func (n *NodeServer) savePeers() error {
	// save to database
	return nil
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

}

func (n *NodeServer) requestPool() {
	
}
