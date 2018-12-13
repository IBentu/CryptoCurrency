package main

import (
	"bufio"
	"fmt"
	"net"
)

//Communicator is a struct that handles the
type Communicator struct {
	address        string
	recievedPacket chan *Packet
	answerPacket   chan *Packet
}

//NewCommunicator creates a new Communicator and returns it
func NewCommunicator(address string, recievedPacket, answerPacket chan *Packet) *Communicator {
	return &Communicator{address: address, recievedPacket: recievedPacket, answerPacket: answerPacket}
}

// Send sends 1 Packet to address and returns the recieved packet
func (c *Communicator) Send(address string, p Packet) (*Packet, error) {
	conn, _ := net.Dial("tcp", fmt.Sprintf("%s:%d", address, ListenPort))
	// read in input from stdin
	_, err := fmt.Fprintf(conn, string(append(p.bytes(), byte('\n'))))
	if err != nil {
		return nil, err
	}
	// listen for reply
	msg, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return nil, err
	}
	msg = msg[:len(msg)-1]
	return ToPacket([]byte(msg)), nil
}

// Listen listens for oncoming connections, recieves 1 Packet and sends one packet back
func (c *Communicator) Listen() error {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", ListenPort))
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		msg, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			conn.Close()
			continue
		}
		msg = msg[:len(msg)-1]
		p := ToPacket([]byte(msg))
		c.recievedPacket <- p
		p = <-c.answerPacket
		_, err = conn.Write(append(p.bytes(), byte('\n')))
		if err != nil {
			fmt.Println(err)
		}
		conn.Close()
	}
}

// Address returns the address of this communicator
func (c *Communicator) Address() string {
	return c.address
}
