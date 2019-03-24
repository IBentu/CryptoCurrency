package main

import (
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

// Type returns the packet type
func (p *Packet) Type() string {
	return p.requestType
}
