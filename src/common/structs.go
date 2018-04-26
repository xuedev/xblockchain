package common

import (
	"math/big"
)

/////////////////////define block////////////////////////
type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

/////////////////////define block chain////////////////////////
type Blockchain struct {
	Blocks []*Block
}

/////////////////////define ProofOfWork////////////////////////
type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

const TargetBits int64 = 16
