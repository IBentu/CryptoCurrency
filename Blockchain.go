package main

import (
	"errors"
	"sync"
)

// Blockchain is the database for all the blocks
type Blockchain struct {
	blocks  []*Block
	hashMap map[string]*Block
	mutex   *sync.Mutex
}

// init initiates the blockchain at node startup
func (bc *Blockchain) init() {

}

// firstInit initiates the blockchain at the first startup
func (bc *Blockchain) firstInit() {

}

//verifyBlock verifies the Block is valid
func (bc *Blockchain) verifyBlock(b Block) bool {
	switch {
	case b.prevHash != bc.blocks[len(bc.blocks)-1].hash:
		return false
	case !b.verifyPOW():
		return false
	default:
		return true
	}
}

// getLatestIndex returns the indexes of the latests block
func (bc *Blockchain) getLatestIndex() int {
	bc.mutex.Lock()
	length := len(bc.blocks) - 1
	bc.mutex.Unlock()
	return length
}

// getLatestHash returns the indexes of the latests block
func (bc *Blockchain) getLatestHash() string {
	bc.mutex.Lock()
	hash := bc.blocks[len(bc.blocks)-1].hash
	bc.mutex.Unlock()
	return hash
}

func (bc *Blockchain) addBlock(b *Block) {
	bc.mutex.Lock()
	bc.blocks = append(bc.blocks, b)
	bc.hashMap[b.hash] = b
	bc.mutex.Unlock()
}

func (bc *Blockchain) addBlocks(blocks []*Block) {
	for _, b := range blocks {
		bc.addBlock(b)
	}
}

// isBlockValid validates a block is valid hash-wise and index-wise
func (bc *Blockchain) isBlockValid(b Block) bool {
	bc.mutex.Lock()
	valid := bc.blocks[len(bc.blocks)-1].index == b.index && bc.blocks[len(bc.blocks)-1].hash == b.prevHash
	bc.mutex.Unlock()
	return valid
}

// length returns the current length of the blockchain
func (bc *Blockchain) length() int {
	bc.mutex.Lock()
	length := len(bc.blocks)
	bc.mutex.Unlock()
	return length
}

func (bc *Blockchain) getBlock(index int) Block {
	bc.mutex.Lock()
	b := bc.blocks[index]
	bc.mutex.Unlock()
	return *b
}

func (bc *Blockchain) compareBlockchains(blocks []*Block) (int, error) {
	for i := blocks[len(blocks)-1].index; i >= blocks[0].index; i-- {
		if blocks[i].hash == bc.getLatestHash() {
			return i, nil
		}
	}
	return 0, errors.New("No Matched Blocks")
}
