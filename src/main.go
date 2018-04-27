// main.go
package main

import (
	"github.com/xuedev/xblockchain/src/common"
)

func main() {

	blockchain := common.NewBlockChain()
	defer blockchain.Db.Close()
	cli := common.CLI{blockchain}
	//cli.AddBlock("xgg")
	cli.Run()

}
