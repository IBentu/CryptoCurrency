package main

import (
	"crypto/ecdsa"
)

// TODO Transaction doc
type Transaction struct {
	senderKey    ecdsa.PublicKey
	recipientKey ecdsa.PublicKey
	amount       int
	senderSign   string
	timestamp    int64
}

func (t *Transaction) toString() string {
	return pubKeyToString(t.senderKey) + pubKeyToString(t.recipientKey) + string(t.amount) + t.senderSign + string(t.timestamp)
}

func pubKeyToString(k ecdsa.PublicKey) string {
	return string(k.X.Bytes()) + string(k.Y.Bytes())
}