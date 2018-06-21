package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"sync"
)

// NodeServer is the server of the node and it is responsible for communication between nodes
type NodeServer struct {
	peers       []*net.UDPAddr
	address     string
	mutex       *sync.Mutex
	recvChannel chan *Packet
	sendChannel chan *Packet
}

// Packet is the struct for transferring data between Nodes
type Packet struct {
	srcAddress  string
	dstAddress  string
	requestType string
	data        []byte
}

const (
	// networkProtocol is the protocol used in the network layer (UDP with IPv4)
	networkProtocol string = "udp"
)

func (n *NodeServer) init(address string, recvChannel, sendChannel chan *Packet) {
}

func (n *NodeServer) firstInit(address string, recvChannel, sendChannel chan *Packet) {
	n.mutex = &sync.Mutex{}
	n.peers = make([]*net.UDPAddr, 0) // read from a certain file a few first peers
	n.address = address
	n.recvChannel = recvChannel
	n.sendChannel = sendChannel
	//go n.listenForPeers()
	//go n.SyncBlockchain()
	//go n.SyncTransactionPool()
	//go n.sendToPeers()
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
	defer conn.Close()
	recvBytes := make([]byte, 4096)
	length, address, err := conn.ReadFromUDP(recvBytes)
	if err != nil {
		fmt.Print(err)
		return
	}
	var p Packet
	buffer := bytes.NewBuffer(recvBytes[:length])
	decoder := gob.NewDecoder(buffer)
	decoder.Decode(&p)
	n.addPeer(address)
	n.recvChannel <- &p
}

// addPeer calls doesPeerExist and adds the address to Peers if the address can not be found in there
func (n *NodeServer) addPeer(address *net.UDPAddr) {
	if !n.doesPeerExist(address) {
		n.mutex.Lock()
		n.peers = append(n.peers, address)
		n.mutex.Unlock()
	}
}

// doesPeerExist checks if the connected peer is already listed in the
// Peers slice
func (n *NodeServer) doesPeerExist(address *net.UDPAddr) bool {
	n.mutex.Lock()
	peersLen := len(n.peers)
	n.mutex.Unlock()
	for i := 0; i < peersLen; i++ {
		n.mutex.Lock()
		if address.String() == n.peers[i].String() {
			n.mutex.Unlock()
			return true
		}
		n.mutex.Unlock()
	}
	return false
}

// requestPeers requests peers from known nodes according to protocol
func (n *NodeServer) requestPeers() {
}

// sendToPeer receives data from channel and sends to given address (from data)
func (n *NodeServer) sendToPeer() {
	var buffer bytes.Buffer
	for {
		p := <-n.sendChannel
		dstAddr, err := net.ResolveUDPAddr(networkProtocol, p.dstAddress)
		if err != nil {
			fmt.Println(err)
			continue
		}
		srcAddr, err := net.ResolveUDPAddr(networkProtocol, p.srcAddress)
		if err != nil {
			fmt.Println(err)
			continue
		}
		conn, err := net.DialUDP(networkProtocol, srcAddr, dstAddr)
		if err != nil {
			fmt.Println(err)
			continue
		}
		encoder := gob.NewEncoder(&buffer)
		encoder.Encode(p)
		_, err = conn.Write(buffer.Bytes())
		if err != nil {
			fmt.Println(err)
		}
		conn.Close()
		buffer.Reset()
	}
}
