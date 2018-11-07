package main

import (
	"bufio"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"sync"
	"time"

	net "github.com/libp2p/go-libp2p-net"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	pstoremem "github.com/libp2p/go-libp2p-peerstore/pstoremem"
)

var port int
var listen bool

func main() {
	flag.IntVar(&port, "port", 1234, "Port to listen and accept connection to")
	flag.BoolVar(&listen, "listen", false, "set to True if the node just waits for a connection") //not working
	flag.Parse()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	n := Node{
		blockchain:  &Blockchain{blocks: make([]*Block, 0), hashMap: make(map[string]*Block)},
		mutex:       &sync.Mutex{},
		privKey:     key,
		pubKey:      key.PublicKey,
		recvChannel: make(chan *Packet),
		scmChannel:  make(chan *Packet),
		stpmChannel: make(chan *Packet),
	}
	n.server = &NodeServer{recvChannel: n.recvChannel, scmChannel: n.scmChannel, stpmChannel: n.stpmChannel, node: &n, peers: pstore.NewPeerstore(pstoremem.NewKeyBook(), pstoremem.NewAddrBook(), pstoremem.NewPeerMetadata()), mutex: &sync.Mutex{}}
	var ip string
	ipRaw, err := getIPAddress()
	if err != nil {
		fmt.Println("Enter your IPv4:")
		fmt.Scanf("%s\n", &ip)
	} else {
		ip = ipRaw.String()
	}
	n.server.address = ip
	err = n.server.newHost(port, n.privKey)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	var str string
	fmt.Println("1 - listen, 2 - connect")
	fmt.Scanf("%s\n", &str)
	if str == "1" {
		fmt.Println("Waiting for node to connect...")
		n.server.host.SetStreamHandler(P2Pprotocol, func(s net.Stream) {
			fmt.Println("Stream connected.")
			rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))//change to blocker reader/writer
			bytes := make([]byte, 1024)
			allBytes := make([]byte, 0)
			reAd := 0
			err := error(nil)
			length := 0
			for !(reAd == 0 && err == io.EOF) {
				time.Sleep(time.Second)
				reAd, err = rw.Read(bytes)
				if err != nil {
					fmt.Println(err.Error())
				}
				fmt.Println(reAd)
				if reAd != 0 {
					allBytes = append(allBytes, bytes...)
					bytes = make([]byte, 1024)
					length += reAd
				}
			}
			fmt.Println(string(allBytes[:length]))
			s.Close()
			fmt.Println("Stream Closed")
		})
		time.Sleep(time.Minute * 3)
	} else {
		fmt.Println("Enter IPv4 to connect to:")
		var addr string
		fmt.Scanf("%s\n", &addr)
		fmt.Println("Enter peer ID to connect to:")
		var id string
		fmt.Scanf("%s\n", &id)
		fmt.Println("Enter string to send to peer:")
		var str string
		fmt.Scanf("%s\n", &str)
		fmt.Printf("Connecting to /ip4/%s/tcp/%d/ipfs/%s...\n", addr, port, id)
		s, err := n.server.openStream(fmt.Sprintf("/ip4/%s/tcp/%d/ipfs/%s", addr, port, id))
		if err != nil {
			fmt.Print(err.Error())
			return
		}
		fmt.Println("Stream connected.")
		// Create a buffer stream for non blocking read and write.
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
		fmt.Printf("Sending %s to node...\n", str)
		written, _ := rw.WriteString(str)
		fmt.Printf("\"%s\" was sent to node.\n", str[:written])
		time.Sleep(time.Minute)
		s.Close()
		fmt.Println("Stream Closed")
	}
}
