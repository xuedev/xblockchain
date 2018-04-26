// main.go
package main

import (
	"fmt"
	"math/big"

	"github.com/xuedev/xblockchain/src/algorithm"
	"github.com/xuedev/xblockchain/src/block"
)

func main() {

	test := big.NewInt(1)
	fmt.Printf("0x%x\n", test)
	test.Lsh(test, uint(256-24))
	fmt.Printf("0x%x\n", test)

	bc := block.NewBlockChain()
	block.AddBlock(bc, "second block")

	for _, bb := range bc.Blocks {
		fmt.Print("\n")
		fmt.Printf("Pre.Hash:%x\n", bb.PrevBlockHash)
		fmt.Printf("Data:%s\n", bb.Data)
		fmt.Printf("Hash:%x\n", bb.Hash)
		fmt.Println("validate:", algorithm.ValidatePOW(bb))

	}
}
