package main

import (
	"encoding/json"
	"math/big"
)

// JSONTransaction is a struct intended for Json encoding and decoding
type JSONTransaction struct {
	SenderKey    string `json:"senderKey"`
	RecipientKey string `json:"recipientKey"`
	Amount       int    `json:"amount"`
	Timestamp    int64  `json:"timestamp"`
	Hash         string `json:"hash"`
	Sign         string `json:"sign"`
}

// MarshalJSON is an Implementation of Marshaler
func (t *Transaction) MarshalJSON() ([]byte, error) {
	jt := JSONTransaction{
		SenderKey:    t.senderKey,
		RecipientKey: t.recipientKey,
		Amount:       t.amount,
		Timestamp:    t.timestamp,
		Hash:         t.hash,
		Sign:         t.sign,
	}
	return json.Marshal(jt)
}

// UnmarshalJSON is an Implementation of Unmarshaler
func (t *Transaction) UnmarshalJSON(data []byte) error {
	var jt JSONTransaction
	if err := json.Unmarshal(data, &jt); err != nil {
		return err
	}
	*t = Transaction{
		senderKey:    jt.SenderKey,
		recipientKey: jt.RecipientKey,
		amount:       jt.Amount,
		timestamp:    jt.Timestamp,
		hash:         jt.Hash,
		sign:         jt.Sign,
	}
	return nil
}

//------------------------------------------------------------------------------------------------------------------

// JSONBlock is a struct intended for Json encoding and decoding
type JSONBlock struct {
	Index        int            `json:"index"`
	Timestamp    int64          `json:"timestamp"`
	Transactions []*Transaction `json:"transactions"`
	Miner        string         `json:"miner"`
	Hash         string         `json:"hash"`
	PrevHash     string         `json:"prevHash"`
	Filler       *big.Int       `json:"filler"`
}

// MarshalJSON is an Implementation of Marshaler
func (b *Block) MarshalJSON() ([]byte, error) {
	jb := JSONBlock{
		Index:        b.index,
		Timestamp:    b.timestamp,
		Transactions: b.transactions,
		Miner:        b.miner,
		Hash:         b.hash,
		PrevHash:     b.prevHash,
		Filler:       b.filler,
	}
	return json.Marshal(jb)
}

// UnmarshalJSON is an Implementation of Unmarshaler
func (b *Block) UnmarshalJSON(data []byte) error {
	var jb JSONBlock
	if err := json.Unmarshal(data, &jb); err != nil {
		return err
	}
	*b = Block{
		index:        jb.Index,
		timestamp:    jb.Timestamp,
		transactions: jb.Transactions,
		miner:        jb.Miner,
		hash:         jb.Hash,
		prevHash:     jb.PrevHash,
		filler:       jb.Filler,
	}
	return nil
}

//------------------------------------------------------------------------------------------------------------------------------

// JSONNode is a data type for the json settings file
type JSONNode struct {
	FirstInit  bool   `json:"FirstInit"`
	PrivateKey string `json:"PrivateKey"`
	PublicKey  string `json:"PublicKey"`
}

// checkError calls panic() with the recieved error in case err != nil
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

//--------------------------------------------------------------------------------------------------------------

//JSONConfig is
type JSONConfig struct {
	Addr  string
	Node  JSONNode
	Peers string
}
