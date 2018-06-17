package main

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
)

// Block is the database for the transaction, blockchain node
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

// verifyPOW verifies if the Proof-of-Work is valid in the block
func (b *Block) verifyPOW() bool {
	hashBytes := []byte(b.hash)
	leadingZeros := 5 // leadingZeros is the number of leading zeros required for the POW
	for i := 0; i < leadingZeros; i++ {
		if hashBytes[i] != 48 { // 48 is the value of the char '0'
			return false
		}
	}
	return true
}
