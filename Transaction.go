package main

import "fmt"

// Transaction is a single transaction and is saved on the blockchain in it
type Transaction struct {
	senderKey    string
	recipientKey string
	amount       int
	timestamp    int64
	hash         string
	sign         string
}

// toString returns all the Transaction's fields that need to be hashed as a formatted
func (t *Transaction) toHashString() string {
	return fmt.Sprintf("%s%s%d%d", t.senderKey, t.recipientKey, t.amount, t.timestamp)
}

//transactionSliceToByteSlice returns a string that can be hashed
func transactionSliceToHashString(transactions []*Transaction) string {
	str := ""
	for i := 0; i < len(transactions); i++ {
		str += transactions[i].toHashString()
	}
	return str
}

// Format formats a Transaction to a []byte
func (t *Transaction) Format() ([]byte, error) {
	return t.MarshalJSON()
}

// UnformatTransaction formats a []byte to a Transaction
func UnformatTransaction(data []byte) (*Transaction, error) {
	t := &Transaction{}
	err := t.UnmarshalJSON(data)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// Equals compares two Transactions
func (t *Transaction) Equals(t2 *Transaction) bool {
	return t.hash == t2.hash
}
