package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
)

// Transaction is a single transaction and is saved on the blockchain in it
type Transaction struct {
	senderKey    ecdsa.PublicKey
	recipientKey ecdsa.PublicKey
	amount       int
	timestamp    int64
	hash         string
	signR        *big.Int
	signS        *big.Int
}

func (t *Transaction) toString() string {
	return pubKeyToString(t.senderKey) + pubKeyToString(t.recipientKey) + string(t.amount) + string(t.timestamp)
}

// returns a string of a PublicKey
func pubKeyToString(k ecdsa.PublicKey) string {
	return string(k.X.Bytes()) + string(k.Y.Bytes())
}

// hashTransaction hashes the transaction
func (t *Transaction) hashTransaction() {
	hash := sha256.New()
	hash.Write([]byte(t.toString()))
	checksum := hash.Sum(nil)
	t.hash = hex.EncodeToString(checksum)
}

// sign signs a Transaction with a PrivateKey
func (t *Transaction) sign(k *ecdsa.PrivateKey) error {
	r, s, err := ecdsa.Sign(rand.Reader, k, []byte(t.hash))
	if err == nil {
		t.signR = r
		t.signS = s
	}
	return err
}
