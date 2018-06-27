package main

import (
	"bytes"
)

// Packet is the struct for transferring data between Nodes
type Packet struct {
	srcAddress  string
	dstAddress  string
	requestType string
	data        []byte
}

//  bytes converts the Packet to a slice of bytes
func (p *Packet) bytes() []byte {
	s := p.srcAddress + "~\000" + p.dstAddress + "~\000" + p.requestType + "~\000"
	return append([]byte(s), p.data...)
}

// to Packet converts a slice of bytes into a Packet
func toPacket(b []byte) *Packet {
	bs := bytes.Split(b, []byte("~\000"))
	p := &Packet{
		srcAddress:  string(bs[0]),
		dstAddress:  string(bs[1]),
		requestType: string(bs[2]),
		data:        bs[3],
	}
	return p
}
