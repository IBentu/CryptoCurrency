package main

import (

	//"math/rand"

	"fmt"
	"time"

	ec "github.com/IBentu/CryptoCurrency/EClib"
)

func main() {
	runNode()
	//testWallet()
}

func runNode() {
	config, err := readJSON()
	checkError(err)
	var node Node
	if config.Node.FirstInit {
		config.Node.FirstInit = false
		err = writeJSON(config)
		checkError(err)
		node.firstInit(config)
	} else {
		node.init(config)
	}
}

func testWallet() {
	config, err := readJSON()
	checkError(err)
	var node Node
	node.firstInit(config)
	_, pub := ec.ECGenerateKey()
	node.mine()
	for i := 0; i < 2; i++ {
		node.makeTransaction(pub, 5)
	}
	node.mine()
	fmt.Println("done!")
	time.Sleep(10 * time.Minute)
}
