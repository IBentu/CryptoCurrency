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
func (n *Node) init() {
	settings, err := readJSON()
	checkError(err)
	n.mutex = &sync.Mutex{}
	n.privKey = &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     big.NewInt(settings.PrivateKey.PublicKey.X),
			Y:     big.NewInt(settings.PrivateKey.PublicKey.Y),
		},
		D: big.NewInt(settings.PrivateKey.D),
	}
	n.pubKey = n.privKey.PublicKey
	n.recvChannel = make(chan *Packet)
	n.scmChannel = make(chan *Packet)
	n.stpmChannel = make(chan *Packet)
	n.server.init(n, settings.Address, n.recvChannel, n.scmChannel, n.stpmChannel, n.privKey)
	n.blockchain.init()
	n.transactionPool.init()

	//go n.makeSCM()
	//go n.makeSTPM()
}

//firstInit initiates the Node for the first time, and saves to a json settings file
func (n *Node) firstInit() {
	// check file
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Print(err)
		return
	}
	n.privKey = key
	n.pubKey = key.PublicKey
	n.mutex = &sync.Mutex{}
	n.blockchain.firstInit()
	n.transactionPool.firstInit()
	n.recvChannel = make(chan *Packet)
	n.scmChannel = make(chan *Packet)
	n.stpmChannel = make(chan *Packet)
	IP, err := getIPAddress()
	if err != nil {
		fmt.Print(err)
		return
	}
	n.server.firstInit(n, IP, n.recvChannel, n.scmChannel, n.stpmChannel, n.privKey)
	n.blockchain.firstInit()
	// update blockchain + transactionPool
	n.saveData()
}

// saveData saves the node's data in the settings file and the sqlite database
func (n *Node) saveData() {
	settings, err := readJSON()
	if err != nil {
		fmt.Print(err)
		return
	}
	n.mutex.Lock()
	settings.PrivateKey.D = n.privKey.D.Int64()
	settings.PrivateKey.PublicKey.X = n.pubKey.X.Int64()
	settings.PrivateKey.PublicKey.Y = n.pubKey.Y.Int64()
	settings.Address = n.server.address
	n.mutex.Unlock()
	err = writeJSON(settings)
	if err != nil {
		fmt.Print(err)
		return
	}
}

// verifyTransaction checks the blockchain if the transaction is legal (enough credits to send), and verifies the transactionSign
func (n *Node) verifyTransaction(t *Transaction) bool {
	signed := ecdsa.Verify(&t.senderKey, []byte(t.hash), t.signR, t.signS)
	validBalance := t.amount <= n.checkBalance(t.senderKey)
	return signed && validBalance
}

// mine creates a block using the TransactionPool, returns true if a block was created and false otherwise
func (n *Node) mine() bool {
	n.mutex.Lock()
	if n.transactionPool.length() == 0 {
		n.mutex.Unlock()
		return false
	}
	n.mutex.Unlock()
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
	block.index = n.blockchain.getLatestIndex() + 1
	block.prevHash = n.blockchain.getLatestHash()
	var counter int64
	for {
		block.filler = big.NewInt(counter)
		counter++
		block.updateHash()
		if block.verifyPOW() {
			if n.blockchain.isBlockValid(block) { // incase the blockchain was updated while mining
				block.index = n.blockchain.getLatestIndex() + 1
				block.prevHash = n.blockchain.getLatestHash()
				continue
			}
			n.blockchain.addBlock(&block)
			return true
		}
	}
}

// checkBalance goes through the blockchain, checks and returns the balance of a certain PublicKey
func (n *Node) checkBalance(key ecdsa.PublicKey) int {
	sum := 0
	for i := 0; i < n.blockchain.length(); i++ {
		if n.blockchain.getBlock(i).miner == key {
			sum += 50 // decide how much money to reward miners. for now 50
		}
		for j := 0; j < len(n.blockchain.getBlock(i).transactions); j++ {
			if n.blockchain.getBlock(i).transactions[j].senderKey == key {
				sum -= n.blockchain.getBlock(i).transactions[j].amount
			} else if n.blockchain.getBlock(i).transactions[j].recipientKey == key {
				sum += n.blockchain.getBlock(i).transactions[j].amount
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

func (n *Node) makeSCM() {
	for {
		time.Sleep(5 * time.Second)
		p := NewPacket(SCM, FormatSCM(n.blockchain.getLatestIndex(), n.blockchain.getLatestHash()))
		n.scmChannel <- p
	}
}

/*
CompareSCM compares by Blockchain Sync Protocol (refer to protocol doc):
returns 0 if scenario 1
returns 0 if scenario 2
returns the difference between indexes otherwise
*/
func (n *Node) CompareSCM(index int) int {
	currIndex := n.blockchain.getLatestIndex()
	if index <= currIndex {
		return 0
	}
	return index - currIndex
}

func (n *Node) makeSTMP() {
	for {
		time.Sleep(5 * time.Second)
		p := NewPacket(STPM, n.transactionPool.formatSTPM())
		n.stpmChannel <- p
	}
}
