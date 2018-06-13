package main

import "crypto/ecdsa"

// TODO Block doc
type Block struct {
	miner       ecdsa.PublicKey
	timestamp   int64
	hash        string
	prevHash    string
	Transaction []Transaction
	filler      string
}

// updateHash updates the block hash
func (b *Block) updateHash() {
}

// verifyFiller checks if the block hash is approved by the consensus
func (b *Block) verifyFiller() bool {
	return false
}
