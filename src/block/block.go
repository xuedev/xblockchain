// block.go
package block

import (
	"crypto/sha256"

	"github.com/xuedev/xblockchain/src/algorithm"
	"github.com/xuedev/xblockchain/src/common"
	"github.com/xuedev/xblockchain/src/util"
)

/////////////////////define block func////////////////////////
//func (b *common.Block) SetHash() {
//	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
//	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
//	hash := sha256.Sum256(headers)
//	b.Hash = hash[:]
//}

func NewBlock(data string, preBlockHash []byte) *common.Block {
	bb := &common.Block{util.GetTimestampInMilli(), []byte(data), preBlockHash, []byte(""), 0}
	//block.SetHash()
	pow := algorithm.NewProofOfWork(bb)
	nonce, hash := algorithm.RunPOW(pow)
	bb.Hash = hash[:]
	bb.Nonce = nonce

	return bb
}

////////////////define blockchain func/////////////////////////

func NewGenesisBlock() *common.Block {
	first256 := sha256.Sum256([]byte("xuedev"))
	return NewBlock("Genesis Block", first256[:])
}

func AddBlock(bc *common.Blockchain, data string) {
	if len(bc.Blocks) < 1 {
		bc.Blocks = append(bc.Blocks, NewGenesisBlock())
	}

	preBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(data, preBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}

func NewBlockChain() *common.Blockchain {
	return &common.Blockchain{[]*common.Block{NewGenesisBlock()}}
}
