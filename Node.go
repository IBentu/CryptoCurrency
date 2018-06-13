package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
)

// TODO Node doc
type Node struct {
	privKey         *ecdsa.PrivateKey
	pubKey          ecdsa.PublicKey
	blockchain      []Block
	transactionPool []Transaction
	server          NodeServer
	EllipticCurve   elliptic.Curve
	dataChannel     chan []byte
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

// makeTransaction create a trnsaction adds it to the pool and returns true if transaction is legal,
// otherwise it returns false
func (n *Node) makeTransaction(recipient ecdsa.PublicKey, amount int) bool {
	return false
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

// checks to see if a request is received from the NodeServer
func (n *Node) sampleDataQ() {
	for {
		data := <-n.dataChannel
		go n.handleRequest(data)
	}
}

// handleRequest handles the requests from the dataQueue in the NodeServer
func (n *Node) handleRequest(request []byte) {
}
