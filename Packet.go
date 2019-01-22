package main

import (
	"bytes"
	"errors"
)

const (
	// BR is Blockchain Request
	BR = "Blockchain Request"
	// TPR is Transaction-Pool Request
	TPR = "Transaction-Pool Request"
	// PR is Peers Request
	PR = "Peers Request"
	// SCM is Sync-Chain-Message
	SCM = "Sync-Chain-Message"
	// FT is From-Top
	FT = "From-Top"
	// IS is Index-Specific
	IS = "Index-Specific"
	// STPM is Sync-Transaction-Pool-Message
	STPM = "Sync-Transaction-Pool-Message"
	// NT is New-Transaction
	NT = "New-Transaction"
	// PA is Peer-Addresses
	PA = "Peer-Addresses"
	// BP is Blocks-Packet
	BP = "Blocks-Packet"
)

// Packet is the struct for transferring data between Nodes
type Packet struct {
	requestType string
	data        []byte
}

var (
	// ErrPacketType is an error for a packet with the wrong message type
	ErrPacketType = errors.New("Invalid Packet Type")
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

// Type returns the packet type
func (p *Packet) Type() string {
	return p.requestType
}
