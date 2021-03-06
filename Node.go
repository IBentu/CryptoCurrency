package main

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	ec "github.com/IBentu/CryptoCurrency/EClib"
)

// Node is the client for miners
type Node struct {
	privKey         string
	pubKey          string
	blockchain      *Blockchain
	transactionPool *TransactionPool
	server          *NodeServer
	mutex           *sync.Mutex
}

const (

	// UpdateInterval is the update time interval (in seconds)
	// of the blockchain, peers and transactions
	UpdateInterval = 15

	// SaveInterval is the save time interval (in seconds)
	// of the blockchain and peers
	SaveInterval = 20
)

// init initiates the Node by loading a json settings file
func (n *Node) init(config *JSONConfig) {
	n.mutex = &sync.Mutex{}
	n.privKey = config.Node.PrivateKey
	n.pubKey = config.Node.PublicKey
	n.server = &NodeServer{}
	n.server.init(n, config)
	n.blockchain = &Blockchain{}
	n.blockchain.init()
	n.transactionPool = &TransactionPool{}
	n.transactionPool.init()
	n.updateFromPeers()
	go n.periodicSave()
	fmt.Println("The node is up!")
	n.PrintBlockchain()
}

// PrintBlockchain Prints the hashed of the blockchain
func (n *Node) PrintBlockchain() {
	fmt.Printf("The hashes of the blockchain:\n%s", n.blockchain.HashString())
}

// saveConfig saves the node's data in the config file
func (n *Node) saveConfig() error {
	config, err := readJSON()
	if err != nil {
		return err
	}
	n.mutex.Lock()
	peers := n.server.peers
	n.mutex.Unlock()
	n.server.savePeers(config, peers)
	return writeJSON(config)
}

// verifyTransaction checks the blockchain if the transaction is legal (enough credits to send), and verifies the transactionSign, and also double spending
func (n *Node) verifyTransaction(t *Transaction) bool {
	signed := ec.ECVerify(t.hash, t.sign, t.senderKey)
	hash := ec.ECHashString(t.toHashString())
	validHash := hash == t.hash
	validBalance := t.amount <= n.checkBalance(t.senderKey)
	noDoubleSpending := !n.blockchain.DoesTransactionExist(t)
	return signed && validBalance && noDoubleSpending && validHash
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
		block.nuance = big.NewInt(counter)
		counter++
		block.updateHash()
		if block.verifyPOW() {
			if n.blockchain.IsBlockValid(block) { // incase the blockchain was updated while mining
				n.transactionPool.addTransactions(transactionsToMake)
				return false
			}
			n.blockchain.AddBlock(&block)
			n.PrintBlockchain()
			return true
		}
	}
}

// checkBalance goes through the blockchain, checks and returns the balance of a certain PublicKey
func (n *Node) checkBalance(key string) int {
	sum := 0
	for i := 1; i < n.blockchain.Length(); i++ {
		currBlock := n.blockchain.GetBlock(i)
		if currBlock.miner == key {
			sum += 20
		}
		for j := 0; j < len(currBlock.transactions); j++ {
			if currBlock.transactions[j].senderKey == key {
				sum -= currBlock.transactions[j].amount
			} else if currBlock.transactions[j].recipientKey == key {
				sum += currBlock.transactions[j].amount
			}
		}
	}
	return sum
}

// makeTransaction create a transaction adds it to the pool and returns true if transaction is legal,
// otherwise it returns false
func (n *Node) makeTransaction(recipient string, amount int) bool {
	var t Transaction
	cb := n.checkBalance(n.pubKey)
	if amount > cb {
		return false
	}
	t.amount = amount
	t.recipientKey = recipient
	t.senderKey = n.pubKey
	t.timestamp = GetCurrentMillis()
	t.hash = ec.ECHashString(t.toHashString())
	t.sign = ec.ECSign(t.hash, n.privKey, n.pubKey)
	n.transactionPool.addTransaction(&t)
	return true
}

// updateFromPeers updates the blockchain, peers and transactions from
// the peer-nodes.
func (n *Node) updateFromPeers() {
	go func() {
		for {
			n.server.requestBlockchain()
			time.Sleep(time.Second * UpdateInterval)
		}
	}()
	go func() {
		for {
			n.server.requestPeers()
			time.Sleep(time.Second * UpdateInterval)
		}
	}()
	go func() {
		for {
			n.server.requestPool()
			time.Sleep(time.Second * UpdateInterval)
		}
	}()
}

// periodicSave saves the blockchain and the config every 30 seconds
func (n *Node) periodicSave() {
	for {
		err1 := n.saveConfig()
		err2 := n.blockchain.saveBlockchain()
		if err1 != nil {
			fmt.Printf("could not save config, error:\n	%s\n", err1)
		}
		if err2 != nil {
			fmt.Println(err2)
		}
		time.Sleep(time.Second * SaveInterval)
	}
}
