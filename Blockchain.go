package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"sync"
)

// Blockchain is the database for all the blocks
type Blockchain struct {
	blocks   []*Block
	mutex    *sync.Mutex
	updating bool
}

// SetUpdating changes the update status of the blockchain
func (bc *Blockchain) SetUpdating(status bool) {
	bc.mutex.Lock()
	bc.updating = status
	bc.mutex.Unlock()
}

// IsUpdating checks if the blockchain is updating
func (bc *Blockchain) IsUpdating() bool {
	bc.mutex.Lock()
	status := bc.updating
	bc.mutex.Unlock()
	return status
}

// init initiates the blockchain at node startup
func (bc *Blockchain) init() {
	bc.blocks = []*Block{}
	bc.mutex = &sync.Mutex{}
	bc.updating = false
	fmt.Println(bc.readBlockchain())
}

func (bc *Blockchain) saveBlockchain() error {
	currDir, err := os.Getwd()
	if err != nil {
		return err
	}
	errList := ""
	if bc.IsUpdating() {
		return errors.New("cannot save blockchain while it's in use")
	}
	bc.mutex.Lock()
	for i := 1; i < len(bc.blocks); i++ {
		dir := path.Join(currDir, fmt.Sprintf("Config/Blockchain/%d.block", i))
		data, err := bc.blocks[i].MarshalJSON()
		if err != nil {
			errList += strconv.Itoa(i) + " "
			continue
		}
		err = ioutil.WriteFile(dir, data, 0644)
		if err != nil {
			errList += strconv.Itoa(i) + " "
			continue
		}
	}
	bc.mutex.Unlock()
	if len(errList) > 2 {
		return fmt.Errorf("failed to save blocks at indexes: %s", errList)
	}
	return nil
}

func (bc *Blockchain) readBlockchain() error {
	currDir, err := os.Getwd()
	if err != nil {
		return err
	}
	blocks := make([]*Block, 0)
	bc.mutex.Lock()
	dir := path.Join(currDir, "Config/Blockchain/0.block")
	data, err := ioutil.ReadFile(dir)
	if err != nil {
		fmt.Println(errors.New("Error reading the origin block"))
		os.Exit(1)
	}
	b := &Block{}
	err = b.UnmarshalJSON(data)
	if err != nil {
		fmt.Println(errors.New("Error reading the origin block"))
		os.Exit(1)
	}
	blocks = append(blocks, b)
	i := 1
	for ; ; i++ {
		dir = path.Join(currDir, fmt.Sprintf("Config/Blockchain/%d.block", i))
		data, err := ioutil.ReadFile(dir)
		if err != nil {
			break
		}
		b := &Block{}
		err = b.UnmarshalJSON(data)
		if err != nil {
			break
		}
		blocks = append(blocks, b)
	}
	bc.blocks = blocks
	bc.mutex.Unlock()
	if i == 1 {
		return errors.New("loaded the origin of the blockchain")
	}
	return fmt.Errorf("loaded blockchain from the origin to index %d", i-1)
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

// DoesTransactionExist checks if a given transaction already happened in the blockchain
func (bc *Blockchain) DoesTransactionExist(t *Transaction) bool {
	blocks := bc.getCopy()
	for _, block := range blocks {
		for _, transaction := range block.transactions {
			if transaction.hash == t.hash {
				return true
			}
		}
	}
	return false
}

// getCopy returns a copy of the blockchain
func (bc *Blockchain) getCopy() []Block {
	copy := []Block{}
	bc.mutex.Lock()
	for _, block := range bc.blocks {
		copy = append(copy, *block)
	}
	bc.mutex.Unlock()
	return copy
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

// GetBlocksFromTop returns the number of blocks from the top of the blockchain from the received number
func (bc *Blockchain) GetBlocksFromTop(num int) []*Block {
	index := bc.GetLatestIndex() - num
	bc.mutex.Lock()
	blocks := bc.blocks[index:]
	bc.mutex.Unlock()
	return blocks
}

// GetBlocksFromIndex returns blocks from the specified index until ten block before it (or the genesis block)
func (bc *Blockchain) GetBlocksFromIndex(index int) []*Block {
	firstIndex := index - 10
	if firstIndex < 0 {
		firstIndex = 0
	}
	bc.mutex.Lock()
	blocks := bc.blocks[firstIndex:index]
	bc.mutex.Unlock()
	return blocks
}

// CompareBlockchains compares the current blockchain top block's hash the recieved blocks's bottom block's hash
// and returns true if they are the same
func (bc *Blockchain) CompareBlockchains(blocks []*Block) bool {
	hash, err := bc.GetHash(blocks[0].index)
	if err != nil {
		return false
	}
	if hash == blocks[0].hash {
		return true
	}
	return false
}

// ReplaceBlocks replaces a part of the blockchain with the recieved blocks
func (bc *Blockchain) ReplaceBlocks(blocks []*Block) {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()
	index := blocks[0].index
	bc.blocks = bc.blocks[:index]
	bc.blocks = append(bc.blocks, blocks...)
}

// HashString returns a string of the hashes of the blockchain
func (bc *Blockchain) HashString() string {
	str := ""
	bc.mutex.Lock()
	defer bc.mutex.Unlock()
	for _, b := range bc.blocks {
		str += b.hash + "\n"
	}
	return str
}
