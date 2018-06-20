package main

import (
	"sync"
)

// TransactionPool is the data structure that holds pending transactions
type TransactionPool struct {
	transactions []*Transaction
	mutex        *sync.Mutex
}

// firstInit is the function that initiates the TP for the first time
func (tp *TransactionPool) firstInit() {
}

// init initiates the TP after each startup
func (tp *TransactionPool) init() {
}

// length returns the length of the transactions slice
func (tp *TransactionPool) length() int {
	tp.mutex.Lock()
	length := len(tp.transactions)
	tp.mutex.Unlock()
	return length
}

// remove returns the first transaction from the transaction slice and removes it from the slice
func (tp *TransactionPool) remove() *Transaction {
	var t *Transaction
	tp.mutex.Lock()
	if len(tp.transactions) > 0 {
		t, tp.transactions = tp.transactions[0], tp.transactions[1:]
	}
	tp.mutex.Unlock()
	return t
}

// addTransaction add a transactin to the pending transaction slice
func (tp *TransactionPool) addTransaction(t *Transaction) {
	tp.mutex.Lock()
	tp.transactions = append(tp.transactions, t)
	tp.mutex.Unlock()
}
