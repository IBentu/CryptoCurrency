package main

import (
	"fmt"
	"net"
	"sync"
)

// TODO NodeServer doc
type NodeServer struct {
	peers       []string
	address     string
	mutex       *sync.Mutex
	dataChannel chan []byte
}

const (
	// networkProtocol is the protocol used in the network layer (UDP with IPv4)
	networkProtocol string = "udp"
)

func (n *NodeServer) init() {
	// init from file settings
}

func (n *NodeServer) firstInit(address string, dataChannel chan []byte) {
	n.mutex = &sync.Mutex{}
	n.peers = make([]string, 0) // read from a certain file a few first peers
	n.address = address
	n.dataChannel = dataChannel
	go n.listenForPeers()
	go n.SyncBlockchain()
	go n.SyncTransactionPool()
}

// SyncBlockchain sends a SCM (refer to protocol doc) to all peers repeatedly.
func (n *NodeServer) SyncBlockchain() {
}

// SyncTransactionPool syncs the transactionpool according to protocol
func (n *NodeServer) SyncTransactionPool() {
}

// listenForPeers listens to other nodes for UDP connections
func (n *NodeServer) listenForPeers() {
	listenAddr, err := net.ResolveUDPAddr(networkProtocol, n.address)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := net.ListenUDP(networkProtocol, listenAddr)
		if err == nil {
			go n.handlePeer(conn)
		}
	}
}

// handlePeer handles a connection from another node
func (n *NodeServer) handlePeer(conn *net.UDPConn) {
	data := make([]byte, 1024)
	_, address, err := conn.ReadFromUDP(data)
	if err != nil {
		fmt.Print(err)
		return
	}
	n.addPeer(address)
	n.dataChannel <- data
}

// addPeer calls doesPeerExist and adds the address to Peers if the address can not be found in there
func (n *NodeServer) addPeer(address net.Addr) {
	if !n.doesPeerExist(address) {
		n.mutex.Lock()
		n.peers = append(n.peers, address.String())
		n.mutex.Unlock()
	}
}

// doesPeerExist checks if the connected peer is already listed in the
// Peers slice
func (n *NodeServer) doesPeerExist(address net.Addr) bool {
	n.mutex.Lock()
	peersLen := len(n.peers)
	n.mutex.Unlock()
	for i := 0; i < peersLen; i++ {
		n.mutex.Lock()
		if address.String() == n.peers[i] {
			return true
		}
		n.mutex.Unlock()
	}
	return false
}

// requestPeers requests peers from known nodes according to protocol
func (n *NodeServer) requestPeers() {

}
