package main

import (
	"encoding/json"
	"math/big"
	"os"
)

func main() {
	file, err := os.Open("./Settings.json")
	if err != nil {
		panic(err)
	}
	data := [1024]byte{}
	i, err := file.Read(data[:])
	checkError(err)
	settings := JSONSettings{}
	err = json.Unmarshal(data[:i], settings)
	checkError(err)
	node := Node{}
	if settings.FirstInit {
		settings.FirstInit = false
		settData, err := json.Marshal(settings)
		checkError(err)
		_, err = file.Write(settData)
		checkError(err)
		err = file.Close()
		checkError(err)
		node.firstInit()
	} else {
		node.init()
	}
}

// JSONSettings is a data type for the json settings file
type JSONSettings struct {
	FirstInit  bool           `json:"FirstInit"`
	PrivateKey JSONPrivateKey `json:"PrivateKey"`
	Address    string         `json:"Address"`
}

// JSONPrivateKey is a data sub-type for the json settings file
type JSONPrivateKey struct {
	PublicKey JSONPublicKey `json:"PublicKey"`
	D         big.Int       `json:"D"`
}

// JSONPublicKey is a data sub-type for the json settings file
type JSONPublicKey struct {
	X big.Int `json:"X"`
	Y big.Int `json:"Y"`
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
