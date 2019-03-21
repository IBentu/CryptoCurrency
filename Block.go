package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

// Block is the database for the transaction, blockchain node
type Block struct {
	index        int
	timestamp    int64
	miner        string
	transactions []*Transaction
	prevHash     string
	filler       *big.Int
	hash         string
}

// updateHash updates the block hash
func (b *Block) updateHash() {
	hash := sha256.New()
	data := fmt.Sprintf("%d%d%s%s%s%d", b.index, b.timestamp, b.miner,
		transactionSliceToString(b.transactions), b.prevHash, b.filler)
	hash.Write([]byte(data))
	hashChecksum := hash.Sum(nil)
	b.hash = hex.EncodeToString(hashChecksum)
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

// ToBytes converts a Block to an array of bytes
func (b *Block) ToBytes() ([]byte, error) {
	return b.MarshalJSON()
}

// ToBlock converts an array of bytes and returns a pointer to a Block
func ToBlock(data []byte) (*Block, error) {
	b := &Block{}
	err := b.UnmarshalJSON(data)
	if err != nil {
		return nil, err
	}
	return b, nil
}
