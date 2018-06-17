package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"sync"
	"time"
)

// TODO Node doc
type Node struct {
	privKey         *ecdsa.PrivateKey
	pubKey          ecdsa.PublicKey
	blockchain      []Block
	transactionPool []Transaction
	server          NodeServer
	EllipticCurve   elliptic.Curve
	recvChannel     chan []byte
	mutex           *sync.Mutex
}

/*
Init Functions
*/

// init initiates the Node by either loading a settings file or calling firstInit
func (n *Node) init() {
}

//firstInit initiates the Node for the first time saves to settings file
func (n *Node) firstInit() {
}

//verifyBlock verifies the Block is valid
func (n *Node) verifyBlock(b Block) bool {
	switch {
	case b.prevHash != n.blockchain[len(n.blockchain)-1].hash:
		return false
	case b.verifyPOW():
		return false
	default:
		return true
	}
}

// verifyTransaction checks the blockchain if the transaction is legal (enough credits to send), and verifies the transactionSign
func (n *Node) verifyTransaction(t Transaction) bool {
	return ecdsa.Verify(&t.senderKey, []byte(t.hash), t.signR, t.signS)
}

// mine creates a block using the TransactionPool
func (n *Node) mine() Block {
	return Block{}
}

// checkBalance goes through the blockchain, checks and returns the balance of a certain PublicKey
func (n *Node) checkBalance(key ecdsa.PublicKey) int {
	sum := 0
	for i := 0; i < len(n.blockchain); i++ {
		if n.blockchain[i].miner == key {
			sum += 0 // decide how much money to reward miners
		}
		for j := 0; j < len(n.blockchain[i].transactions); j++ {
			if n.blockchain[i].transactions[j].senderKey == key {
				sum -= n.blockchain[i].transactions[j].amount
			} else if n.blockchain[i].transactions[j].recipientKey == key {
				sum += n.blockchain[i].transactions[j].amount
			}
		}
	}
	return sum
}

// makeTransaction create a trnsaction adds it to the pool and returns true if transaction is legal,
// otherwise it returns false
func (n *Node) makeTransaction(recipient ecdsa.PublicKey, amount int) bool {
	t := Transaction{}
	if amount < n.checkBalance(n.pubKey) {
		return false
	}
	t.amount = amount
	t.recipientKey = recipient
	t.senderKey = n.pubKey
	t.timestamp = getCurrentMillis()
	t.hashTransaction()
	err := t.sign(n.privKey)
	if err != nil {
		return false
	}
	n.mutex.Lock()
	n.transactionPool = append(n.transactionPool, t)
	n.mutex.Unlock()
	return true
}

// handleSCM handles every SCM
func (n *Node) handleSCM(hash string, index int) {
}

// compare SCM compares by Blockchain Sync Protocol (refer to protocol doc):
// 		returns 0 if scenario 1
// 		returns 0 if scenario 2.i
// 		returns 0 if scenario 2.ii.a
// 		returns -1 if scenario 2.ii.b
// 		returns 1 if scenario 3.i
// 		returns -2 if scenario 3.ii
func (n *Node) compareSCM(hash string, index int) int {
	return 0
}

// getServerData checks to see if a request is received from the NodeServer
func (n *Node) getServerData() {
	for {
		data := <-n.recvChannel
		go n.handleRequest(data)
	}
}

// handleRequest handles the requests from the dataQueue in the NodeServer
func (n *Node) handleRequest(request []byte) {
}

// getCurrentMillis returns the current time in millisecs
func getCurrentMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
