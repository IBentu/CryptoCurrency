package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"sync"
	"testing"
)

func TestCheckBalance(t *testing.T) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Error(err)
	}
	node := Node{
		pubKey:     key.PublicKey,
		blockchain: &Blockchain{},
	}
	node.blockchain.mutex = &sync.Mutex{}
	node.blockchain.blocks = make([]*Block, 2)
	node.blockchain.blocks[0], node.blockchain.blocks[1] = &Block{}, &Block{}

	node.blockchain.blocks[0].transactions = make([]*Transaction, 2)
	node.blockchain.blocks[1].transactions = make([]*Transaction, 1)

	node.blockchain.blocks[0].miner = ecdsa.PublicKey{}
	node.blockchain.blocks[1].miner = node.pubKey
	node.blockchain.blocks[0].transactions[0] = &Transaction{amount: 3, recipientKey: node.pubKey, senderKey: ecdsa.PublicKey{}}
	node.blockchain.blocks[0].transactions[1] = &Transaction{amount: 9, recipientKey: node.pubKey, senderKey: ecdsa.PublicKey{}}
	node.blockchain.blocks[1].transactions[0] = &Transaction{amount: 2, recipientKey: ecdsa.PublicKey{}, senderKey: node.pubKey}
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

func TestMine(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	var node Node
	key1, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Error(err)
	}
	node.privKey = key1
	node.pubKey = key1.PublicKey
	node.mutex = &sync.Mutex{}
	node.blockchain = &Blockchain{mutex: &sync.Mutex{}, hashMap: make(map[string]*Block), blocks: make([]*Block, 1)}
	node.blockchain.blocks[0] = &Block{hash: "0000000020422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3", index: 0}
	node.transactionPool = &TransactionPool{transactions: make([]*Transaction, 3), mutex: &sync.Mutex{}}

	key2, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Error(err)
	}
	pKey1 := key1.PublicKey
	pKey2 := key2.PublicKey
	node.transactionPool.transactions[0] = &Transaction{amount: 0, senderKey: pKey1, recipientKey: pKey2}
	node.transactionPool.transactions[0].hashTransaction()
	err = node.transactionPool.transactions[0].sign(key1)
	if err != nil {
		t.Error(err)
	}
	node.transactionPool.transactions[1] = &Transaction{amount: 0, senderKey: pKey2, recipientKey: pKey1}
	node.transactionPool.transactions[1].hashTransaction()
	err = node.transactionPool.transactions[1].sign(key2)
	if err != nil {
		t.Error(err)
	}
	node.transactionPool.transactions[2] = &Transaction{amount: 0, senderKey: pKey1, recipientKey: pKey2}
	node.transactionPool.transactions[2].hashTransaction()
	err = node.transactionPool.transactions[2].sign(key1)
	if err != nil {
		t.Error(err)
	}

	expected := Block{index: 1, prevHash: node.blockchain.blocks[0].hash}
	if !node.mine() {
		t.Error("No transactions in transactionPool")
	} else {
		if node.blockchain.length() != 2 || node.blockchain.blocks[1].index != expected.index || node.blockchain.blocks[1].prevHash != expected.prevHash {
			t.Error("Invalid block mined")
		}
	}
}

func TestTranactionMarshaler(t *testing.T) {
	tran := Transaction{
		amount: 30,
	}
	b, err := tran.MarshalJSON()
	if err != nil {
		t.Error(err)
	}
	newT := Transaction{}
	err = newT.UnmarshalJSON(b)
	if err != nil {
		t.Error(err)
	}
	if newT.amount != tran.amount {
		t.Error("Transactions not the same")
	}
}

func TestBlockMarshaler(t *testing.T) {
	b := Block{
		index: 30,
	}
	bytes, err := b.MarshalJSON()
	if err != nil {
		t.Error(err)
	}
	newB := Block{}
	err = newB.UnmarshalJSON(bytes)
	if err != nil {
		t.Error(err)
	}
	if newB.index != b.index {
		t.Error("Blocks not the same")
	}
}
