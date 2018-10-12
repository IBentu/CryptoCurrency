package main

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"fmt"
	"sync"

	peerStructs "github.com/CryptoCurrency/peer-structs"
	libp2p "github.com/libp2p/go-libp2p"
	p2pCrypto "github.com/libp2p/go-libp2p-crypto"
	p2pHost "github.com/libp2p/go-libp2p-host"
	p2pNet "github.com/libp2p/go-libp2p-net"
	p2pPeerstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
)

// NodeServer is the server of the node and it is responsible for communication between nodes
type NodeServer struct {
	node        *Node
	peers       *p2pPeerstore.Peerstore
	address     string
	mutex       *sync.Mutex
	recvChannel chan *Packet
	scmChannel  chan *Packet
	host        p2pHost.Host
}

const (

	// ListenPort is the IP on which the Server is listening to
	ListenPort = 1625

	// SCM is Sync-Chain-Message
	SCM = "SCM"
	// FT is From-Top
	FT = "FT"
	// IS is Index-Specific
	IS = "IS"
	// STPM is Sync-Transaction-Pool-Message
	STPM = "STPM"
	// NT is New-Transaction
	NT = "NT"
	// PA is Peer-Addresses
	PA = "PA"
)

func (n *NodeServer) init(node *Node, address string, recvChannel, sendChannel chan *Packet, privKey *ecdsa.PrivateKey) {
}

func (n *NodeServer) firstInit(node *Node, address string, recvChannel, scmChannel chan *Packet, privKey *ecdsa.PrivateKey) {
	n.node = node
	n.mutex = &sync.Mutex{}
	n.peers = p2pPeerstore.NewPeerstore(peerStructs.NewKeyBook(), peerStructs.NewAddrBook(), peerStructs.NewPeerMetadata())
	n.address = address
	n.recvChannel = recvChannel
	n.scmChannel = scmChannel
	err := n.newHost(ListenPort, privKey)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	//go n.listenForPeers()
	//go n.SyncBlockchain()
	//go n.SyncTransactionPool()
	//go n.sendToPeers()
}

func (n *NodeServer) newHost(listenPort int, privKey *ecdsa.PrivateKey) error {

	priv, _, err := p2pCrypto.ECDSAKeyPairFromKey(privKey)
	if err != nil {
		return err
	}
	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", listenPort)),
		libp2p.Identity(priv),
	}
	host, err := libp2p.New(context.Background(), opts...)
	if err != nil {
		return err
	}

	// Build host multiaddress
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", host.ID().Pretty()))
	addr := host.Addrs()[0]
	fullAddr := addr.Encapsulate(hostAddr)
	fmt.Printf("I am %s\n", fullAddr)

	n.mutex.Lock()
	n.host = host
	n.mutex.Unlock()

	return nil
}

// HandleStream handles an incoming peer stream
func (n *NodeServer) HandleStream(s p2pNet.Stream) {
	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	go func() {
		// sends SCM to peer every few seconds
		for {
			scm := <-n.scmChannel
			n.mutex.Lock()
			_, err := rw.Write(scm.bytes())
			if err != nil {
				fmt.Print(err)
			}
			n.mutex.Unlock()
		}
	}()

	for {
		data := make([]byte, 0)
		l, err := rw.Read(data)
		if err != nil {
			fmt.Print(err.Error())
			continue
		}

		if l == 0 {
			continue
		}
		p := ToPacket(data)
		switch p.requestType {
		case SCM:
			index, hash, err := UnformatSCM(p.data)
			if err != nil {
				fmt.Printf(err.Error())
				continue
			}
			res := n.node.CompareSCM(index)
			if res == 0 {
				continue
			} else {
				newP := Packet{
					requestType: FT,
					data:        FormatFT(res),
				}
				_, err := rw.Write(newP.bytes())
				if err == nil {
					fmt.Print(err.Error())
					continue
				}

			}
		case FT:
		case IS:
		case STPM:
		case NT:
		case PA:
		}
	}

}
