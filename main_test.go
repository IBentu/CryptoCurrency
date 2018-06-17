package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"
)

func TestCheckBalance(t *testing.T) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Error(err)
	}
	node := Node{
		pubKey:     key.PublicKey,
		blockchain: []Block{},
	}
	node.blockchain = make([]Block, 2)

	node.blockchain[0].transactions = make([]Transaction, 2)
	node.blockchain[1].transactions = make([]Transaction, 1)

	node.blockchain[0].miner = ecdsa.PublicKey{}
	node.blockchain[1].miner = node.pubKey
	node.blockchain[0].transactions[0] = Transaction{amount: 3, recipientKey: node.pubKey, senderKey: ecdsa.PublicKey{}}
	node.blockchain[0].transactions[1] = Transaction{amount: 9, recipientKey: node.pubKey, senderKey: ecdsa.PublicKey{}}
	node.blockchain[1].transactions[0] = Transaction{amount: 2, recipientKey: ecdsa.PublicKey{}, senderKey: node.pubKey}
	sum := node.checkBalance(node.pubKey)
	expected := 60
	if sum != expected {
		t.Errorf("The balance was %d, expected %d", sum, expected)
	}
}

func TestVerifyPOW(t *testing.T) {
	block := Block{hash: "00000000bdce850ab654f102f0f57dbb4fb09516852fe94298f7fac77a77d8ef"}
	if !block.verifyPOW() {
		t.Error("POW verification failed")
	}
}
