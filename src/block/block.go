// block.go
package block

import (
	"bytes"
	"crypto/sha256"
	"strconv"

	"github.com/xuedev/xblockchain/src/util"
)

/////////////////////define block////////////////////////
type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
}

func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}

func NewBlock(data string, preBlockHash []byte) *Block {
	block := &Block{util.GetTimestampInMilli(), []byte(data), preBlockHash, []byte("")}
	block.SetHash()
	return block
}

////////////////define blockchain/////////////////////////
type Blockchain struct {
	Blocks []*Block
}

func NewGenesisBlock() *Block {
	first256 := sha256.Sum256([]byte("xuedev"))
	return NewBlock("Genesis Block", first256[:])
}

func (bc *Blockchain) AddBlock(data string) {
	if len(bc.Blocks) < 1 {
		bc.Blocks = append(bc.Blocks, NewGenesisBlock())
	}

	preBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(data, preBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}

func NewBlockChain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}
