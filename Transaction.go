package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"math/big"
	"strconv"
)

// Transaction is a single transaction and is saved on the blockchain in it
type Transaction struct {
	senderKey    ecdsa.PublicKey
	recipientKey ecdsa.PublicKey
	amount       int
	timestamp    int64
	hash         string
	signR        *big.Int
	signS        *big.Int
}

// toString returns all the Transaction's fields that need to be hashed as a formatted
func (t *Transaction) toString() string {
	return pubKeyToString(t.senderKey) + pubKeyToString(t.recipientKey) + string(t.amount) + string(t.timestamp)
}

//transactionSliceToByteSlice returns a byte slice that can be hashed
func transactionSliceToString(transactions []*Transaction) string {
	str := ""
	for i := 0; i < len(transactions); i++ {
		str += transactions[i].toString()
	}
	return str
}

// returns a string of a PublicKey
func pubKeyToString(k ecdsa.PublicKey) string {
	return string(k.X.Bytes()) + string(k.Y.Bytes())
}

// hashTransaction hashes the transaction
func (t *Transaction) hashTransaction() {
	hash := sha256.New()
	hash.Write([]byte(t.toString()))
	checksum := hash.Sum(nil)
	t.hash = hex.EncodeToString(checksum)
}

// sign signs a Transaction with a PrivateKey
func (t *Transaction) sign(k *ecdsa.PrivateKey) error {
	r, s, err := ecdsa.Sign(rand.Reader, k, []byte(t.hash))
	if err == nil {
		t.signR = r
		t.signS = s
	}
	return err
}

// formatTransaction formats a Transaction to a []byte
func (t *Transaction) formatTransaction() []byte {
	data := append(formatPublicKey((&t.senderKey)), []byte("\\\000")...)
	data = append(append(data, formatPublicKey(&t.recipientKey)...), []byte("\\\000")...)
	data = append(append(data, []byte(strconv.Itoa(t.amount))...), []byte("\\\000")...)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(t.timestamp))
	data = append(append(data, b...), []byte("\\\000")...)
	data = append(append(data, []byte(t.hash)...), []byte("\\\000")...)
	data = append(append(data, t.signR.Bytes()...), []byte("\\\000")...)
	data = append(append(data, t.signS.Bytes()...), []byte("\\\000")...)
	return data
}

// unformatTransaction formats a []byte to a Transaction
func unformatTransaction(data []byte) (*Transaction, error) {
	splat := bytes.Split(data, []byte("\\\000"))
	var t Transaction
	t.senderKey = unformatPublicKey(splat[0])
	t.recipientKey = unformatPublicKey(splat[1])
	amount, err := strconv.Atoi(string(splat[2]))
	if err != nil {
		return nil, err
	}
	t.amount = amount
	t.timestamp = int64(binary.LittleEndian.Uint64(splat[3]))
	t.hash = string(splat[4])
	r := big.NewInt(0)
	r.SetBytes(splat[5])
	t.signR = r
	s := big.NewInt(0)
	s.SetBytes(splat[6])
	t.signS = s
	return &t, nil
}

func formatPublicKey(key *ecdsa.PublicKey) []byte {
	return append(append(key.X.Bytes(), []byte("/\000")...), key.Y.Bytes()...)
}

func unformatPublicKey(data []byte) ecdsa.PublicKey {
	splat := bytes.Split(data, []byte("/\000"))
	var key ecdsa.PublicKey
	key.Curve = elliptic.P256()
	x := big.NewInt(0)
	x.SetBytes(splat[0])
	key.X = x
	y := big.NewInt(0)
	y.SetBytes(splat[1])
	key.Y = y
	return key
}

// Equals compares two Transactions
func (t *Transaction) Equals(t2 *Transaction) bool {
	return t.hash == t2.hash
}
