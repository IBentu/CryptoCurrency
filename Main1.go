package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"sync"

	pstore "github.com/libp2p/go-libp2p-peerstore"
	pstoremem "github.com/libp2p/go-libp2p-peerstore/pstoremem"
)

const (
	lstnPort = 3445
)

func main() {
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
	fmt.Println("Enter your IPv4:")
	var ip string
	fmt.Scanf("%s\n", &ip)
	n.server.address = ip
	err = n.server.newHost(lstnPort, n.privKey)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println("Enter IPv4 to connect to:")
	var addr string
	fmt.Scanf("%s\n", &addr)
	fmt.Println("Enter peer ID to connect to:")
	var id string
	fmt.Scanf("%s\n", &id)
	s, err := n.server.openStream(fmt.Sprintf("/ip4/%s/tcp/%d/ipfs/%s", addr, lstnPort, id))
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	n.server.HandleStream(s)
	for {
	}
}
