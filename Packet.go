package main

import (
	"bytes"
)

// Packet is the struct for transferring data between Nodes
type Packet struct {
	requestType string
	data        []byte
}

//  bytes converts the Packet to a slice of bytes
func (p *Packet) bytes() []byte {
	s := p.requestType + "~\000"
	return append([]byte(s), p.data...)
}

// ToPacket converts a slice of bytes into a Packet
func ToPacket(b []byte) *Packet {
	bs := bytes.Split(b, []byte("~\000"))
	p := &Packet{
		requestType: string(bs[0]),
		data:        bs[1],
	}
	return p
}
