package main

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"fmt"
	"sync"

	libp2p "github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-crypto"
	host "github.com/libp2p/go-libp2p-host"
	net "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	pstoremem "github.com/libp2p/go-libp2p-peerstore/pstoremem"
	ma "github.com/multiformats/go-multiaddr"
)

// NodeServer is the server of the node and it is responsible for communication between nodes
type NodeServer struct {
	node        *Node
	peers       pstore.Peerstore
	address     string
	mutex       *sync.Mutex
	recvChannel chan *Packet
	scmChannel  chan *Packet
	host        host.Host
}

const (

	// ListenPort is the IP on which the Server is listening to
	ListenPort = 1625

	// P2Pprotocol is the peer-to-peer protocol the nodes use
	P2Pprotocol = "/p2p/1.0.0"

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
	// BP is Blocks-Packet
	BP = "BP"
)

func (n *NodeServer) init(node *Node, address string, recvChannel, sendChannel chan *Packet, privKey *ecdsa.PrivateKey) {
}

func (n *NodeServer) firstInit(node *Node, address string, recvChannel, scmChannel chan *Packet, privKey *ecdsa.PrivateKey) {
	n.node = node
	n.mutex = &sync.Mutex{}
	n.peers = pstore.NewPeerstore(pstoremem.NewKeyBook(), pstoremem.NewAddrBook(), pstoremem.NewPeerMetadata())
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

	priv, _, err := crypto.ECDSAKeyPairFromKey(privKey)
	if err != nil {
		return err
	}
	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", listenPort)),
		libp2p.Identity(priv),
	}
	hst, err := libp2p.New(context.Background(), opts...)
	if err != nil {
		return err
	}

	// Build host multiaddress
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", hst.ID().Pretty()))
	addr := hst.Addrs()[0]
	fullAddr := addr.Encapsulate(hostAddr)
	fmt.Printf("I am %s\n", fullAddr)

	n.mutex.Lock()
	n.host = hst
	n.host.SetStreamHandler(P2Pprotocol, n.HandleStream)
	n.mutex.Unlock()

	return nil
}

func (n *NodeServer) openStream(target string) (net.Stream, error) {
	// The following code extracts target's peer ID from the
	// given multiaddress
	ipfsaddr, err := ma.NewMultiaddr(target)
	if err != nil {
		return nil, err
	}

	pid, err := ipfsaddr.ValueForProtocol(ma.P_IPFS)
	if err != nil {
		return nil, err
	}

	peerid, err := peer.IDB58Decode(pid)
	if err != nil {
		return nil, err
	}

	// Decapsulate the /ipfs/<peerID> part from the target
	// /ip4/<a.b.c.d>/ipfs/<peer> becomes /ip4/<a.b.c.d>
	targetPeerAddr, _ := ma.NewMultiaddr(
		fmt.Sprintf("/ipfs/%s", peer.IDB58Encode(peerid)))
	targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)

	// We have a peer ID and a targetAddr so we add it to the pstore
	// so LibP2P knows how to contact it
	n.peers.AddAddr(peerid, targetAddr, pstore.PermanentAddrTTL)
	n.host.Peerstore().AddAddr(peerid, targetAddr, pstore.PermanentAddrTTL)

	fmt.Print("opening stream")
	// make a new stream from host B to host A
	// it should be handled on host A by the handler we set above because
	// we use the same /p2p/1.0.0 protocol
	s, err := n.host.NewStream(context.Background(), peerid, "/p2p/1.0.0")
	if err != nil {
		return nil, err
	}
	return s, nil
}

// HandleStream handles an incoming peer stream
func (n *NodeServer) HandleStream(s net.Stream) {
	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go n.writeStream(rw)
	go n.readStream(rw)

}

func (n *NodeServer) writeStream(rw *bufio.ReadWriter) {
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
}

func (n *NodeServer) readStream(rw *bufio.ReadWriter) {
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
			index, _, err := UnformatSCM(p.data)
			if err != nil {
				fmt.Printf(err.Error())
				continue
			}
			res := n.node.CompareSCM(index)
			if res == 0 {
				continue
			} else { // scenario 3
				newP := NewPacket(FT, FormatFT(res))
				_, err := rw.Write(newP.bytes())
				if err == nil {
					fmt.Print(err.Error())
					continue
				}
				var bytes []byte
				_, err = rw.Read(bytes)
				if err == nil {
					fmt.Print(err.Error())
					continue
				}
				recvP := ToPacket(bytes)
				if recvP.requestType != BP {
					fmt.Print(ErrPacketType.Error())
					continue
				}
				recvBlockchain := UnformatBP(recvP.data)
				if recvBlockchain[0].hash == n.node.blockchain.getLatestHash() {
					n.node.blockchain.addBlocks(recvBlockchain[1:])
				} else {
					newP := NewPacket(IS, FormatIS(recvBlockchain[0].index-1))
					_, err := rw.Write(newP.bytes())
					if err == nil {
						fmt.Print(err.Error())
						continue
					}
					var prevBlocks []*Block
					_, err = rw.Read(bytes)
					if err == nil {
						fmt.Print(err.Error())
						continue
					}
					recvP := ToPacket(bytes)
					if recvP.requestType != BP {
						fmt.Print(ErrPacketType.Error())
						continue
					}
					recvBlockchain := UnformatBP(recvP.data)
					prevBlocks = recvBlockchain
					var sameIndex int
					sameIndex, err = n.node.blockchain.compareBlockchains(recvBlockchain)
					for err != nil {
						newP := NewPacket(IS, FormatIS(recvBlockchain[0].index-1))
						_, err := rw.Write(newP.bytes())
						if err == nil {
							fmt.Print(err.Error())
							continue
						}
						_, err = rw.Read(bytes)
						if err == nil {
							fmt.Print(err.Error())
							continue
						}
						recvP := ToPacket(bytes)
						if recvP.requestType != BP {
							fmt.Print(ErrPacketType.Error())
							continue
						}
						recvBlockchain := UnformatBP(recvP.data)
						prevBlocks = append(prevBlocks, recvBlockchain...)
						sameIndex, err = n.node.blockchain.compareBlockchains(recvBlockchain)
					}
					prevBlocks = append(prevBlocks, recvBlockchain[sameIndex])
					n.node.blockchain.addBlocks(prevBlocks)
				}

			}
		case STPM:
		case NT:
		case PA:
		default:
			fmt.Print(ErrPacketType.Error())
		}
	}
}
