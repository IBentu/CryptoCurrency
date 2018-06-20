package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
	"sync"
)

// NodeServer is the server of the node and it is responsible for communication between nodes
type NodeServer struct {
	peers       []string
	address     string
	mutex       *sync.Mutex
	recvChannel chan []byte
	sendChannel chan []byte
}

const (
	// networkProtocol is the protocol used in the network layer (UDP with IPv4)
	networkProtocol string = "udp"
)

func (n *NodeServer) init(address string, recvChannel, sendChannel chan []byte) {
}

func (n *NodeServer) firstInit(address string, recvChannel, sendChannel chan []byte) {
	n.mutex = &sync.Mutex{}
	n.peers = make([]string, 0) // read from a certain file a few first peers
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
	data := make([]byte, 1024)
	_, address, err := conn.ReadFromUDP(data)
	if err != nil {
		fmt.Print(err)
		return
	}
	n.addPeer(address)
	port := make([]byte, 2)
	binary.LittleEndian.PutUint16(port, uint16(address.Port))
	addressBytes := []byte{address.IP[15], address.IP[14], address.IP[13], address.IP[12]}
	addressBytes = append(addressBytes, port...)
	n.recvChannel <- append(addressBytes, data...)
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

// sendToPeers receives data from channel and sends to given address (from data)
func (n *NodeServer) sendToPeers() {
	for {
		sendData := <-n.sendChannel
		address := sendData[:6]
		port := binary.LittleEndian.Uint16(address[4:])
		addressS := strconv.Itoa(int(address[0])) + "." + strconv.Itoa(int(address[1])) + "." + strconv.Itoa(int(address[2])) + "." + strconv.Itoa(int(address[3])) + ":" + strconv.Itoa(int(port))
		addr, err := net.ResolveUDPAddr(networkProtocol, addressS)
		if err != nil {
			fmt.Println(err)
			continue
		}
		conn, err := net.DialUDP(networkProtocol, nil, addr)
		if err != nil {
			fmt.Println(err)
			continue
		}
		_, err = conn.Write(sendData[6:])
		if err != nil {
			fmt.Println(err)
		}
	}
}
