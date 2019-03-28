package main

import (
	"fmt"

	ec "github.com/IBentu/CryptoCurrency/EClib"
)

func main() {
	runNode()
}

// runNode runs the initiates the node and runs it
func runNode() {
	config, err := readJSON()
	checkError(err)
	var node Node
	if config.Node.FirstInit {
		priv, pub := ec.ECGenerateKey()
		fmt.Printf("Generated Keys:\n    Private: %s\n    Public: %s\n", priv, pub)
		config.Node.FirstInit = false
		config.Node.PrivateKey = priv
		config.Node.PublicKey = pub
		err = writeJSON(config)
		checkError(err)
	}
	node.init(config)
	select {}
}
