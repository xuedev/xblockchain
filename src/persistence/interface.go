// interface.go
package persistence

import (
	"github.com/xuedev/xblockchain/src/common"
)

type BlockSerialize interface {
	Serialize(b *common.Block) []byte
	DeserializeBlock(d []byte) *common.Block
}
