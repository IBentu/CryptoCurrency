package peerstructs

import peer "github.com/libp2p/go-libp2p-peer"

// PeerMetadata is a struct
type PeerMetadata struct {
}

// NewPeerMetadata returns a new PeerMetadata struct
func NewPeerMetadata() PeerMetadata {
	return PeerMetadata{}
}

// Get ...
func (m *PeerMetadata) Get(p peer.ID, key string) (interface{}, error) {
	return nil, nil
}

// Put ...
func (m *PeerMetadata) Put(p peer.ID, key string, val interface{}) error {
	return nil
}
