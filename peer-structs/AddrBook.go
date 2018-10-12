package peerstructs

import (
	"context"
	"time"

	p2pPeer "github.com/libp2p/go-libp2p-peer"
	ma "github.com/multiformats/go-multiaddr"
)

// AddrBook stores peer addresses
type AddrBook struct {
	addrs map[p2pPeer.ID][]ma.Multiaddr
}

// NewAddrBook creates a new AddrBook struct
func NewAddrBook() AddrBook {
	return AddrBook{addrs: make(map[p2pPeer.ID][]ma.Multiaddr)}
}

// AddAddr calls AddAddrs(p, []ma.Multiaddr{addr}, ttl)
func (a *AddrBook) AddAddr(p p2pPeer.ID, addr ma.Multiaddr, ttl time.Duration) {
	a.AddAddrs(p, []ma.Multiaddr{addr}, ttl)
}

// AddAddrs gives this AddrBook addresses to use, with a given ttl
// (time-to-live), after which the address is no longer valid.
// If the manager has a longer TTL, the operation is a no-op for that address
func (a *AddrBook) AddAddrs(p p2pPeer.ID, addrs []ma.Multiaddr, ttl time.Duration) {
	a.addrs[p] = append(a.addrs[p], addrs...)
}

// SetAddr calls mgr.SetAddrs(p, addr, ttl)
func (a *AddrBook) SetAddr(p p2pPeer.ID, addr ma.Multiaddr, ttl time.Duration) {
	addrs := make([]ma.Multiaddr, 1)
	addrs[0] = addr
	a.SetAddrs(p, addrs, ttl)
}

// SetAddrs sets the ttl on addresses. This clears any TTL there previously.
// This is used when we receive the best estimate of the validity of an address.
func (a *AddrBook) SetAddrs(p p2pPeer.ID, addrs []ma.Multiaddr, ttl time.Duration) {
}

// UpdateAddrs updates the addresses associated with the given peer that have
// the given oldTTL to have the given newTTL.
func (a *AddrBook) UpdateAddrs(p p2pPeer.ID, oldTTL time.Duration, newTTL time.Duration) {
}

// Addrs returns all known (and valid) addresses for a given peer
func (a *AddrBook) Addrs(p p2pPeer.ID) []ma.Multiaddr {
	return a.addrs[p]
}

// AddrStream returns a channel that gets all addresses for a given
// peer sent on it. If new addresses are added after the call is made
// they will be sent along through the channel as well.
func (a *AddrBook) AddrStream(context.Context, p2pPeer.ID) <-chan ma.Multiaddr { // TODO: RESEARCH AND IMPLEMENT
	return nil
}

// ClearAddrs removes all previously stored addresses
func (a *AddrBook) ClearAddrs(p p2pPeer.ID) {
	a.addrs[p] = make([]ma.Multiaddr, 0)
}

// PeersWithAddrs returns all of the peer IDs stored in the AddrBook
func (a *AddrBook) PeersWithAddrs() p2pPeer.IDSlice {
	IDs := make(p2pPeer.IDSlice, 0, len(a.addrs))
	for id := range a.addrs {
		IDs = append(IDs, id)
	}
	return IDs
}
