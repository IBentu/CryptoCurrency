package main

import (
	"bytes"
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

func (tp *TransactionPool) formatSTPM() []byte {
	var data []byte
	for _, v := range tp.transactions {
		data = append(append(data, v.formatTransaction()...), []byte("|\000")...)
	}
	return data
}

//UnformatSTPM unfomrmats []byte to a slice of Transactions
func UnformatSTPM(data []byte) ([]*Transaction, error) {
	splat := bytes.Split(data, []byte("|\000"))
	trans := make([]*Transaction, 0)
	for _, v := range splat {
		t, err := unformatTransaction(v)
		if err != nil {
			return nil, err
		}
		trans = append(trans, t)
	}
	return trans, nil
}

// DoesExists return true if t exists in the TransactionPool and false otherwise
func (tp *TransactionPool) DoesExists(t *Transaction) bool {
	tp.mutex.Lock()
	trans := tp.transactions
	tp.mutex.Unlock()
	for _, v := range trans {
		if t.Equals(v) {
			return true
		}
	}
	return false
}
