package main

import (
	"fmt"
	"net"
	"sync"
)

// NodeServer is the server of the node and it is responsible for communication between nodes
type NodeServer struct {
	peers       []*net.Addr
	address     string
	mutex       *sync.Mutex
	recvChannel chan *Packet
	sendChannel chan *Packet
}

const (
	// networkProtocol is the protocol used in the network layer (UDP with IPv4)
	networkProtocol string = "tcp"
)

func (n *NodeServer) init(address string, recvChannel, sendChannel chan *Packet) {
}

func (n *NodeServer) firstInit(address string, recvChannel, sendChannel chan *Packet) {
	n.mutex = &sync.Mutex{}
	n.peers = make([]*net.Addr, 0) // read from a certain file a few first peers
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
	listenAddr, err := net.ResolveTCPAddr(networkProtocol, n.address)
	if err != nil {
		panic(err)
	}
	for {
		listen, err := net.ListenTCP(networkProtocol, listenAddr)
		conn, err := listen.Accept()
		if err == nil {
			go n.handlePeer(conn)
		}
	}
}

// handlePeer handles a connection from another node
func (n *NodeServer) handlePeer(conn net.Conn) {
	defer conn.Close()
	recvBytes := make([]byte, 4096)
	length, err := conn.Read(recvBytes)
	if err != nil {
		fmt.Print(err)
		return
	}
	p := toPacket(recvBytes[:length])
	n.addPeer(conn.RemoteAddr())
	n.recvChannel <- p
}

// addPeer calls doesPeerExist and adds the address to Peers if the address can not be found in there
func (n *NodeServer) addPeer(address net.Addr) {
	addr := &address
	if !n.doesPeerExist(addr) {
		n.mutex.Lock()
		n.peers = append(n.peers, addr)
		n.mutex.Unlock()
	}
}

// doesPeerExist checks if the connected peer is already listed in the
// Peers slice
func (n *NodeServer) doesPeerExist(address *net.Addr) bool {
	n.mutex.Lock()
	peersLen := len(n.peers)
	n.mutex.Unlock()
	for i := 0; i < peersLen; i++ {
		n.mutex.Lock()
		if (*address).String() == (*n.peers[i]).String() {
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
	for {
		p := <-n.sendChannel
		dstAddr, err := net.ResolveTCPAddr(networkProtocol, p.dstAddress)
		if err != nil {
			fmt.Println(err)
			continue
		}
		srcAddr, err := net.ResolveTCPAddr(networkProtocol, n.address)
		if err != nil {
			fmt.Println(err)
			continue
		}
		conn, err := net.DialTCP(networkProtocol, srcAddr, dstAddr)
		if err != nil {
			fmt.Println(err)
			continue
		}
		_, err = conn.Write(p.bytes())
		if err != nil {
			fmt.Println(err)
		}
		conn.Close()
	}
}

// SR1 sends and recieves 1 Packet
func (n *NodeServer) SR1(p *Packet) (*Packet, error) {
	dstAddr, err := net.ResolveUDPAddr(networkProtocol, p.dstAddress)
	if err != nil {
		return nil, err
	}
	srcAddr, err := net.ResolveUDPAddr(networkProtocol, n.address)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP(networkProtocol, srcAddr, dstAddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	_, err = conn.Write(p.bytes())
	if err != nil {
		fmt.Println(err)
	}
	recvBytes := make([]byte, 4096)
	length, _, err := conn.ReadFromUDP(recvBytes)
	if err != nil {
		return nil, err
	}
	return toPacket(recvBytes[:length]), nil
}

func (n *NodeServer) getAddress() string {
	n.mutex.Lock()
	address := n.address
	n.mutex.Unlock()
	return address
}
