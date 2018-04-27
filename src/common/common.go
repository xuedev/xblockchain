package common

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"math"
	"math/big"
	"os"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/xuedev/xblockchain/src/util"
)

/////////////////////define block////////////////////////
type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		return []byte("")
	}
	return result.Bytes()
}

/////////////////////define block chain////////////////////////
/*** version 1****
type Blockchain struct {
	Blocks []*Block
}
*****/
type Blockchain struct {
	Tip []byte
	Db  *bolt.DB
}

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.Tip, bc.Db}

	return bci
}

func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		return nil
	}
	return &block
}

func (i *BlockchainIterator) Next() *Block {
	var block *Block

	i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})

	fmt.Printf("-------------%x\n", block.PrevBlockHash)
	i.currentHash = block.PrevBlockHash

	return block
}

/////////////////////define ProofOfWork////////////////////////
type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func Int2Hex(val int64) []byte {
	return []byte(fmt.Sprintf("%x", val))
}

func (pow *ProofOfWork) prepareData(block *Block, nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			block.PrevBlockHash,
			block.Data,
			Int2Hex(block.Timestamp),
			Int2Hex(int64(TargetBits)),
			Int2Hex(int64(nonce)),
		},
		[]byte{},
	)
	return data
}

func RunPOW(pow *ProofOfWork) (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0
	fmt.Printf("Mining the block containing \"%s\"\n", pow.Block.Data)
	for nonce < math.MaxInt64 {
		data := pow.prepareData(pow.Block, nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])
		fmt.Println("nonce:", nonce)
		fmt.Printf("Result:%x\n", hash)
		fmt.Printf("Target:%x\n", pow.Target)
		fmt.Println("-----------------------")
		//fmt.Println("nonce:", nonce, "result:", hashInt, ":", pow.Target)
		if hashInt.Cmp(pow.Target) == -1 {
			fmt.Printf("\r%x", hash)
			break

		} else {
			nonce++
		}
	}

	fmt.Print("\n\n")
	return nonce, hash[:]

}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-TargetBits))
	pow := &ProofOfWork{b, target}
	return pow
}
func (pow *ProofOfWork) Validate(block *Block) bool {
	var hashInt big.Int

	data := pow.prepareData(block, block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	isValid := hashInt.Cmp(pow.Target) == -1

	return isValid
}

///////////////////////define CLI//////////////////////////////
type CLI struct {
	Bc *Blockchain
}

func (cli *CLI) addBlock(data string) {
	cli.Bc.AddBlock(data)
	fmt.Println("Success!")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  print - Print all the blocks of the blockchain")
	fmt.Println("  add - add block")
}

func (cli *CLI) AddBlock(data string) {
	cli.Bc.AddBlock(data)
}

func (cli *CLI) PrintChain() {
	bci := cli.Bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate(block)))
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

func (cli *CLI) Run() {
	cli.validateArgs()

	switch os.Args[1] {
	case "add":
		cli.AddBlock(os.Args[2])
	case "print":
		cli.PrintChain()
	default:
		cli.printUsage()
		os.Exit(1)
	}

}

////////////////////////////////////////////////////////////////
func NewBlock(data string, preBlockHash []byte) *Block {
	bb := &Block{util.GetTimestampInMilli(), []byte(data), preBlockHash, []byte(""), 0}
	//block.SetHash()

	target := big.NewInt(1)
	target.Lsh(target, uint(256-TargetBits))
	pow := &ProofOfWork{bb, target}

	nonce, hash := RunPOW(pow)
	bb.Hash = hash[:]
	bb.Nonce = nonce

	return bb
}

func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte

	bc.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})

	newBlock := NewBlock(data, lastHash)

	err := bc.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucket))
		b.Put(newBlock.Hash, newBlock.Serialize())
		b.Put([]byte("l"), newBlock.Hash)
		bc.Tip = newBlock.Hash
		fmt.Printf("save newblock:%x\n", newBlock.Hash)
		return nil
	})

	fmt.Println(err)

}

////////////////define blockchain func/////////////////////////

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte(""))
}

func NewBlockChain() *Blockchain {
	var tip []byte
	db, _ := bolt.Open(DbFile, 0600, nil)

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucket))
		if b == nil {
			genesis := NewGenesisBlock()
			b, _ := tx.CreateBucket([]byte(BlocksBucket))
			b.Put(genesis.Hash, genesis.Serialize())
			b.Put([]byte("l"), genesis.Hash)
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})

	bc := Blockchain{tip, db}

	return &bc
}

///////////////////////define argvs////////////////////////////

const TargetBits int64 = 8
const DbFile string = "data"
const BlocksBucket string = "blocks"
