package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"sync"
)

// TODO Node doc
type Node struct {
	privKey         *ecdsa.PrivateKey
	pubKey          ecdsa.PublicKey
	blockchain      []Block
	transactionPool []Transaction
	server          NodeServer
	EllipticCurve   elliptic.Curve
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

/*
Blockchain Functions
*/

//verifyBlock verifies the Block is valid
func (n *Node) verifyBlock(b Block) bool {
	return false
}

// verifyTransaction checks the blockchain if the transaction is legal (enough credits to send), and verifies the senderSign
func (n *Node) verifyTransaction(t Transaction) bool {
	return false
}

// mine creates a block using the TransactionPool
func (n *Node) mine() Block {
	return Block{}
}

// verifyPOW verifies if the Proof-of-Work is valid in a certain block
func (n *Node) verifyPOW(b Block) bool {
	return false
}

// checkBalance goes through the blockchain, checks and returns the balance of a certain PublicKey
func (n *Node) checkBalance(key ecdsa.PublicKey) int {
	return 0
}

//makeTransaction create a trnsaction adds it to the pool and returns true if transaction is legal,
//otherwise it returns false
func (n *Node) makeTransaction(recipient ecdsa.PublicKey, amount int) bool {
	return false
}

/*
NodeServer Functions
*/

// checks to see if a request is received from the NodeServer
func (n *Node) sampleDataQ() {
	for {
		n.mutex.Lock()
		if len(n.server.dataQueue) > 0 {
			var request []byte
			if len(n.server.dataQueue) > 1 {
				request = n.server.dataQueue[0]
				n.server.dataQueue = make([][]byte, 0)
			} else if len(n.server.dataQueue) == 1 {
				request = n.server.dataQueue[0]
				n.server.dataQueue = n.server.dataQueue[1:]
			}
			go n.handleRequest(request)
		}
		n.mutex.Unlock()
	}
}

// handleRequest handles the requests from the dataQueue in the NodeServer
func (n *Node) handleRequest(request []byte) {
}
