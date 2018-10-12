package peerstructs

import (
	p2pCrypto "github.com/libp2p/go-libp2p-crypto"
	p2pPeer "github.com/libp2p/go-libp2p-peer"
)

//KeyBook is a struct that implements the KeyBook interface in go-libp2p-peerstore
type KeyBook struct {
	keys map[p2pPeer.ID]p2pCrypto.PubKey
}

// NewKeyBook returns a new empty KeyBook
func NewKeyBook() KeyBook {
	return KeyBook{keys: make(map[p2pPeer.ID]p2pCrypto.PubKey)}
}

// PubKey stores the public key of a peer.
func (k *KeyBook) PubKey(id p2pPeer.ID) p2pCrypto.PubKey {
	return k.keys[id]
}

// AddPubKey stores the public key of a peer.
func (k *KeyBook) AddPubKey(id p2pPeer.ID, pubKey p2pCrypto.PubKey) error {
	k.keys[id] = pubKey
	return nil
}

// PrivKey returns the private key of a peer.
func (k *KeyBook) PrivKey(id p2pPeer.ID) p2pCrypto.PrivKey {
	return nil
}

// AddPrivKey stores the private key of a peer.
func (k *KeyBook) AddPrivKey(id p2pPeer.ID, privKey p2pCrypto.PrivKey) error {
	return nil
}

// PeersWithKeys returns all the peer IDs stored in the KeyBook
func (k *KeyBook) PeersWithKeys() p2pPeer.IDSlice {
	IDs := make(p2pPeer.IDSlice, 0, len(k.keys))
	for id := range k.keys {
		IDs = append(IDs, id)
	}
	return IDs
}
