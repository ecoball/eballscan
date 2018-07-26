package data

import (
	"fmt"
	"sync"

	"github.com/ecoball/go-ecoball/common/elog"
)

var (
	Blocks = Block{BlocksInfo: make(map[int]BlockInfo, 0)}
	log    = elog.NewLogger("data", elog.DebugLog)
)

type BlockInfo struct {
	Hash       string
	PrevHash   string
	MerkleHash string
	StateHash  string
	CountTxs   int
}

type Block struct {
	BlocksInfo map[int]BlockInfo

	sync.RWMutex
}

func (this *Block) Add(hight int, info BlockInfo) {
	this.Lock()
	defer this.Unlock()

	if _, ok := this.BlocksInfo[hight]; ok {
		return
	}
	this.BlocksInfo[hight] = info
}

func PrintBlock() string {
	Blocks.RLock()
	defer Blocks.RUnlock()

	result := "hight\thash\tprevHash\tmerkleHash\tstateHash\tcountTxs\n"
	for k, v := range Blocks.BlocksInfo {
		result += fmt.Sprintf("%d\t%s\t%s\t%s\t%s\t%d\n", k, v.Hash, v.PrevHash, v.MerkleHash, v.StateHash, v.CountTxs)
	}

	return result
}
