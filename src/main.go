// main.go
package main

import (
	"fmt"

	"github.com/xuedev/xblockchain/src/block"
)

func main() {
	bc := block.NewBlockChain()
	bc.AddBlock("second block")

	for _, bb := range bc.Blocks {
		fmt.Print("\n")
		fmt.Printf("Pre.Hash:%x\n", bb.PrevBlockHash)
		fmt.Printf("Data:%s\n", bb.Data)
		fmt.Printf("Hash:%x\n", bb.Hash)

	}
}
