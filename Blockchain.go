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
	//TODO: Load the blockchain from database
}

// firstInit initiates the blockchain at the first startup
func (bc *Blockchain) firstInit() {
	bc.blocks = make([]*Block, 1)
	bc.hashMap = make(map[string]*Block, 1)
	bc.mutex = &sync.Mutex{}
	//put genesis Block values in first block
}

func (bc *Blockchain) saveBlockchain() error {
	//TODO: save to database
	return nil
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

// GetLatestIndex returns the indexes of the latests block
func (bc *Blockchain) GetLatestIndex() int {
	bc.mutex.Lock()
	length := len(bc.blocks) - 1
	bc.mutex.Unlock()
	return length
}

// GetLatestHash returns the indexes of the latests block
func (bc *Blockchain) GetLatestHash() string {
	bc.mutex.Lock()
	hash := bc.blocks[len(bc.blocks)-1].hash
	bc.mutex.Unlock()
	return hash
}

// GetHash returns the hash of the block in the specified index
func (bc *Blockchain) GetHash(index int) (string, error) {
	if index > bc.GetLatestIndex() || index < 0 {
		return "", errors.New("Index Out of Bounds")
	}
	bc.mutex.Lock()
	hash := bc.blocks[index].hash
	bc.mutex.Unlock()
	return hash, nil
}

//AddBlock adds a block to the blockchain
func (bc *Blockchain) AddBlock(b *Block) {
	bc.mutex.Lock()
	bc.blocks = append(bc.blocks, b)
	bc.hashMap[b.hash] = b
	bc.mutex.Unlock()
}

// AddBlocks adds blocks to the blockchain
func (bc *Blockchain) AddBlocks(blocks []*Block) {
	for _, b := range blocks {
		bc.AddBlock(b)
	}
}

// IsBlockValid validates a block is valid hash-wise and index-wise
func (bc *Blockchain) IsBlockValid(b Block) bool {
	bc.mutex.Lock()
	valid := bc.blocks[len(bc.blocks)-1].index == b.index && bc.blocks[len(bc.blocks)-1].hash == b.prevHash
	bc.mutex.Unlock()
	return valid
}

// Length returns the current length of the blockchain
func (bc *Blockchain) Length() int {
	bc.mutex.Lock()
	length := len(bc.blocks)
	bc.mutex.Unlock()
	return length
}

// GetBlock returns a block in the specified index
func (bc *Blockchain) GetBlock(index int) Block {
	bc.mutex.Lock()
	b := bc.blocks[index]
	bc.mutex.Unlock()
	return *b
}

// CompareBlockchains compares the recieved blockchain's bottom block with the current one's
// top and returns the index of the first blocks who match hash-wise
func (bc *Blockchain) CompareBlockchains(blocks []*Block) bool {
	hash, err := bc.GetHash(len(blocks) - 1)
	if err != nil {
		return false
	}
	return hash == blocks[0].hash
}
