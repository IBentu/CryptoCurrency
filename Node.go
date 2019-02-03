package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"time"
)

// Node is the client for miners
type Node struct {
	privKey         *ecdsa.PrivateKey
	pubKey          ecdsa.PublicKey
	blockchain      *Blockchain
	transactionPool *TransactionPool
	server          *NodeServer
	recvChannel     chan *Packet
	scmChannel      chan *Packet
	stpmChannel     chan *Packet
	mutex           *sync.Mutex
}

// init initiates the Node by loading a json settings file
func (n *Node) init(config *JSONConfig) {
	n.mutex = &sync.Mutex{}
	n.privKey = &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     big.NewInt(config.Node.PrivateKey.PublicKey.X),
			Y:     big.NewInt(config.Node.PrivateKey.PublicKey.Y),
		},
		D: big.NewInt(config.Node.PrivateKey.D),
	}
	n.pubKey = n.privKey.PublicKey
	n.recvChannel = make(chan *Packet)
	n.scmChannel = make(chan *Packet)
	n.stpmChannel = make(chan *Packet)
	n.server = &NodeServer{}
	n.server.init(n, config.Addr, n.recvChannel, n.scmChannel, n.stpmChannel, n.privKey)
	n.blockchain = &Blockchain{}
	n.blockchain.init()
	n.transactionPool = &TransactionPool{}
	n.transactionPool.init()
	go n.updateChain()
	go n.updatePeers()
	go n.updatePool()
}

//firstInit initiates the Node for the first time, and saves to a json settings file
func (n *Node) firstInit(conf *JSONConfig) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Print(err)
		return
	}
	n.privKey = key
	n.pubKey = key.PublicKey
	n.mutex = &sync.Mutex{}
	n.recvChannel = make(chan *Packet)
	n.scmChannel = make(chan *Packet)
	n.stpmChannel = make(chan *Packet)
	n.blockchain.firstInit()
	n.transactionPool.firstInit()
	n.server = &NodeServer{}
	n.server.firstInit(conf, n, n.privKey)
	n.blockchain.firstInit()
	n.server.requestBlockchain()
	n.server.requestPeers()
	err1 := n.saveConfig()
	err2 := n.blockchain.saveBlockchain()
	err3 := n.server.savePeers()
	if err1 != nil {
		fmt.Printf("could not save config\n%s\n", err1)
	}
	if err2 != nil {
		fmt.Printf("could not save blockchain\n%s\n", err2)
	}
	if err3 != nil {
		fmt.Printf("could not save peers\n%s\n", err3)
	}
	go n.updateChain()
	go n.updatePeers()
	go n.updatePool()
}

// saveConfig saves the node's data in the config file
func (n *Node) saveConfig() error {
	config, err := readJSON()
	if err != nil {
		return err
	}
	n.mutex.Lock()
	config.Node.PrivateKey.D = n.privKey.D.Int64()
	config.Node.PrivateKey.PublicKey.X = n.pubKey.X.Int64()
	config.Node.PrivateKey.PublicKey.Y = n.pubKey.Y.Int64()
	config.Addr = n.server.Address()
	n.mutex.Unlock()
	return writeJSON(config)
}

// verifyTransaction checks the blockchain if the transaction is legal (enough credits to send), and verifies the transactionSign, and also double spending
func (n *Node) verifyTransaction(t *Transaction) bool {
	signed := ecdsa.Verify(&t.senderKey, []byte(t.hash), t.signR, t.signS)
	validBalance := t.amount <= n.checkBalance(t.senderKey)
	noDoubleSpending := !n.blockchain.DoesTransactionExist(t)
	return signed && validBalance && noDoubleSpending
}

// mine creates a block using the TransactionPool, returns true if a block was created and false otherwise
func (n *Node) mine() bool {
	var block Block
	block.miner = n.pubKey
	transactionsToMake := make([]*Transaction, 0)
	for n.transactionPool.length() > 0 && len(transactionsToMake) < 5 {
		t := n.transactionPool.remove()
		if n.verifyTransaction(t) {
			transactionsToMake = append(transactionsToMake, t)
		}
	}
	block.transactions = transactionsToMake
	block.timestamp = GetCurrentMillis()
	block.index = n.blockchain.GetLatestIndex() + 1
	block.prevHash = n.blockchain.GetLatestHash()
	var counter int64
	for {
		block.filler = big.NewInt(counter)
		counter++
		block.updateHash()
		if block.verifyPOW() {
			if n.blockchain.IsBlockValid(block) { // incase the blockchain was updated while mining
				n.transactionPool.addTransactions(transactionsToMake)
				return false
			}
			n.blockchain.AddBlock(&block)
			return true
		}
	}
}

// checkBalance goes through the blockchain, checks and returns the balance of a certain PublicKey
func (n *Node) checkBalance(key ecdsa.PublicKey) int {
	sum := 0
	for i := 0; i < n.blockchain.Length(); i++ {
		if n.blockchain.GetBlock(i).miner == key {
			sum += 50 // decide how much money to reward miners. for now 50
		}
		for j := 0; j < len(n.blockchain.GetBlock(i).transactions); j++ {
			if n.blockchain.GetBlock(i).transactions[j].senderKey == key {
				sum -= n.blockchain.GetBlock(i).transactions[j].amount
			} else if n.blockchain.GetBlock(i).transactions[j].recipientKey == key {
				sum += n.blockchain.GetBlock(i).transactions[j].amount
			}
		}
	}
	return sum
}

// makeTransaction create a trnsaction adds it to the pool and returns true if transaction is legal,
// otherwise it returns false
func (n *Node) makeTransaction(recipient ecdsa.PublicKey, amount int) bool {
	var t Transaction
	if amount < n.checkBalance(n.pubKey) {
		return false
	}
	t.amount = amount
	t.recipientKey = recipient
	t.senderKey = n.pubKey
	t.timestamp = GetCurrentMillis()
	t.hashTransaction()
	err := t.sign(n.privKey)
	if err != nil {
		return false
	}
	n.transactionPool.addTransaction(&t)
	return true
}

// updatePool requests an update for the transactionPool from peers
func (n *Node) updatePool() {
	for {
		n.server.requestPool()
		time.Sleep(time.Minute)
	}
}

// updatePeers requests an update for the transactionPool from peers
func (n *Node) updatePeers() {
	for {
		n.server.requestPeers()
		time.Sleep(time.Minute)
	}
}

// updateChain requests an update for the blockchain from peers
func (n *Node) updateChain() {
	for {
		n.server.requestBlockchain()
		time.Sleep(time.Minute)
	}
}
