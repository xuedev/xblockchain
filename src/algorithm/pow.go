package algorithm

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"

	"github.com/xuedev/xblockchain/src/common"
)

func Int2Hex(val int64) []byte {
	return []byte(fmt.Sprintf("%x", val))
}

func NewProofOfWork(b *common.Block) *common.ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-common.TargetBits))
	pow := &common.ProofOfWork{b, target}
	return pow
}

func preparePOWData(block *common.Block, nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			block.PrevBlockHash,
			block.Data,
			Int2Hex(block.Timestamp),
			Int2Hex(int64(common.TargetBits)),
			Int2Hex(int64(nonce)),
		},
		[]byte{},
	)
	return data
}

func RunPOW(pow *common.ProofOfWork) (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0
	fmt.Printf("Mining the block containing \"%s\"\n", pow.Block.Data)
	for nonce < math.MaxInt64 {
		data := preparePOWData(pow.Block, nonce)
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

func ValidatePOW(block *common.Block) bool {
	var hashInt big.Int

	data := preparePOWData(block, block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	pow := NewProofOfWork(block)
	isValid := hashInt.Cmp(pow.Target) == -1

	return isValid
}
