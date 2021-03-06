package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

//Communicator is a struct that handles the
type Communicator struct {
	server         *NodeServer
	address        string
	port           int
	recievedPacket chan *Packet
	answerPacket   chan *Packet
	mutex          *sync.Mutex
}

//NewCommunicator creates a new Communicator and returns it
func NewCommunicator(server *NodeServer, address string, recievedPacket, answerPacket chan *Packet, port int) *Communicator {
	return &Communicator{server: server, address: address, recievedPacket: recievedPacket, answerPacket: answerPacket, port: port, mutex: &sync.Mutex{}}
}

// SR1 sends 1 Packet to address and returns the recieved packet
func (c *Communicator) SR1(address string, p *Packet) (*Packet, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	//fmt.Printf("Connecting to %s:%d...\n", address, c.port)
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", address, c.port))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	bytes, err := p.MarshalJSON()
	if err != nil {
		return nil, err
	}
	_, err = fmt.Fprintf(conn, string(append(bytes, '\n')))
	if err != nil {
		return nil, err
	}
	// listen for reply
	msg, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return nil, err
	}
	msg = msg[:len(msg)-1]
	newP := &Packet{}
	return newP, newP.UnmarshalJSON([]byte(msg))
}

// Listen listens for oncoming connections, recieves 1 Packet and sends one packet back
func (c *Communicator) Listen() error {
	fmt.Println("Listening for nodes...")
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", c.port))
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		peerAddr := conn.RemoteAddr().String()
		peerSplat := strings.Split(peerAddr, ":")
		if len(peerSplat) == 2 {
			c.server.addPeer(peerSplat[0])
		}
		//fmt.Printf("Connected to %s\n", peerAddr)
		msg, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			conn.Close()
			//fmt.Printf("Connection with %s closed due to error:\n	%s\n", peerAddr, err)
			continue
		}
		msg = msg[:len(msg)-1]
		p := &Packet{}
		err = p.UnmarshalJSON([]byte(msg))
		if err != nil {
			conn.Close()
			//fmt.Printf("Connection with %s closed due to error:\n	%s\n", peerAddr, err)
			continue
		}
		c.recievedPacket <- p
		p = <-c.answerPacket
		bytes, err := p.MarshalJSON()
		if err != nil {
			conn.Close()
			//fmt.Printf("Connection with %s closed due to error:\n	%s\n", peerAddr, err)
			continue
		}
		_, err = fmt.Fprintf(conn, string(append(bytes, '\n')))
		if err != nil {
			conn.Close()
			//fmt.Printf("Connection with %s closed due to error:\n	%s\n", peerAddr, err)
			continue
		}
		conn.Close()
		//fmt.Printf("Connection with %s closed\n", peerAddr)
	}
}

// Address returns the address of this communicator
func (c *Communicator) Address() string {
	return c.address
}
