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
