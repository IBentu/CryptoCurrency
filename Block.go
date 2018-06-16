package main

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
)

// TODO Block doc
type Block struct {
	index        int
	timestamp    int64
	transactions []Transaction
	miner        ecdsa.PublicKey
	hash         string
	prevHash     string
	filler       string
}

// updateHash updates the block hash
func (b *Block) updateHash() {
	hash := sha256.New()
	data := string(b.index) + string(b.timestamp) + transactionSliceToString(b.transactions) + pubKeyToString(b.miner) + b.prevHash + b.filler
	hash.Write([]byte(data))
	hashChecksum := hash.Sum(nil)
	b.hash = hex.EncodeToString(hashChecksum)
}

//transactionSliceToByteSlice returns a byte slice that can be hashed
func transactionSliceToString(transactions []Transaction) string {
	str := ""
	for i := 0; i < len(transactions); i++ {
		str += transactions[i].toString()
	}
	return str
}

// verifyFiller checks if the block hash is approved by the consensus
func (b *Block) verifyFiller() bool {
	return false
}
