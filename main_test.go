package main

import (
	"crypto/ecdsa"
	"math/big"
	"testing"
)

func TestCheckBalance(t *testing.T) {
	key := ecdsa.PublicKey{}
	node := Node{
		pubKey:     ecdsa.PublicKey{},
		blockchain: []Block{},
	}
	node.blockchain = make([]Block, 2)

	node.blockchain[0].transactions = make([]Transaction, 2)
	node.blockchain[1].transactions = make([]Transaction, 1)

	node.blockchain[0].miner = ecdsa.PublicKey{X: big.NewInt(0)}
	node.blockchain[1].miner = ecdsa.PublicKey{}
	node.blockchain[0].transactions[0] = Transaction{amount: 3, recipientKey: ecdsa.PublicKey{}, senderKey: ecdsa.PublicKey{X: big.NewInt(0)}}
	node.blockchain[0].transactions[1] = Transaction{amount: 9, recipientKey: ecdsa.PublicKey{}, senderKey: ecdsa.PublicKey{X: big.NewInt(0)}}
	node.blockchain[1].transactions[0] = Transaction{amount: 2, recipientKey: ecdsa.PublicKey{X: big.NewInt(3)}, senderKey: ecdsa.PublicKey{}}
	sum := node.checkBalance(key)
	expected := 60
	if sum != expected {
		t.Errorf("The balance was %d, expected %d", sum, expected)
	}
}
