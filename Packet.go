package main

import (
	"bytes"
	"errors"
)

// Packet is the struct for transferring data between Nodes
type Packet struct {
	requestType string
	data        []byte
}

var (
	// ErrPacketType is an error for a packet with the wrong message type
	ErrPacketType = errors.New("Wrong Packet Type")
)

// NewPacket returns a new packet
func NewPacket(request string, data []byte) *Packet {
	return &Packet{
		requestType: request,
		data:        data,
	}
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
