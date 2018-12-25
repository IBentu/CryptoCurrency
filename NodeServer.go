package main

import (
	"crypto/ecdsa"
	"sync"
)

// NodeServer is the server of the node and it is responsible for communication between nodes
type NodeServer struct {
	node         *Node
	peers        []*ecdsa.PublicKey
	mutex        *sync.Mutex
	communicator *Communicator
	recvChannel  chan *Packet
	sendChannel  chan *Packet
}

const (

	// ListenPort is the IP on which the Server is listening to
	ListenPort = 1625

	// SCM is Sync-Chain-Message
	SCM = "SCM"
	// FT is From-Top
	FT = "FT"
	// IS is Index-Specific
	IS = "IS"
	// STPM is Sync-Transaction-Pool-Message
	STPM = "STPM"
	// NT is New-Transaction
	NT = "NT"
	// PA is Peer-Addresses
	PA = "PA"
	// BP is Blocks-Packet
	BP = "BP"
)

func (n *NodeServer) init(node *Node, address string, recvChannel, sendChannel, stmpChannel chan *Packet, privKey *ecdsa.PrivateKey) {
}

func (n *NodeServer) firstInit(conf *JSONConfig, node *Node, privKey *ecdsa.PrivateKey) {
	n.node = node
	n.mutex = &sync.Mutex{}
	n.peers = make([]*ecdsa.PublicKey, 0)
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

func (n *NodeServer) requestBlockchain() {

}

func (n *NodeServer) requestPeers() {

}
