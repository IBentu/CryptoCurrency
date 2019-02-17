package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	crand "crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"path"
	"strings"
	"time"
)

func main() {
	//realMain()
	tempMain()
}
func realMain() {
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

func tempMain() {
	config, err := readJSON()
	checkError(err)
	var node Node
	node.init(config)
	key, err := ecdsa.GenerateKey(elliptic.P256(), crand.Reader)
	if err != nil {
		fmt.Print(err)
		return
	}
	for i := 0; i < 8; i++ {
		node.makeTransaction(key.PublicKey, rand.Intn(100))
	}
	node.mine()
	node.mine()
	peerStr := config.Peers
	splat := strings.Split(peerStr, ";")
	node.server.peers = append(node.server.peers, splat...)
	time.Sleep(3 * time.Minute)
}

func readJSON() (*JSONConfig, error) {
	dir, err := os.Getwd()
	checkError(err)
	dir = path.Join(dir, "/config.json")
	data, err := ioutil.ReadFile(dir)
	if err != nil {
		return nil, err
	}
	var config JSONConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, err
}

func writeJSON(config *JSONConfig) error {
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	dir, err := os.Getwd()
	checkError(err)
	dir = path.Join(dir, "/config.json")
	err = ioutil.WriteFile(dir, data, 0644)
	return err
}

// getIPAddress returns the local ip address
func getIPAddress() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, i := range ifaces {
		if strings.Index(i.Name, "Wi-Fi") != 0 || strings.Index(i.Name, "eth0") != 0 {
			continue
		}
		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {

			switch v := addr.(type) {
			case *net.IPNet:
				if v.IP[0] == 0 {
					return v.IP, nil
				}
			case *net.IPAddr:
				return v.IP, nil
			}
		}
	}
	return nil, errors.New("IP not found")
}
