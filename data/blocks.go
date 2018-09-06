// Copyright 2018 The eballscan Authors
// This file is part of the eballscan.
//
// The eballscan is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The eballscan is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the eballscan. If not, see <http://www.gnu.org/licenses/>.

package data

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/ecoball/go-ecoball/common/elog"
	"github.com/muesli/cache2go"
)

const (
	BLOCK_SPAN time.Duration = 20 * time.Second
)

var (
	Blocks = cache2go.Cache("Blocks")
	log    = elog.NewLogger("data", elog.DebugLog)
)

type BlockInfo struct {
	Hash       string
	PrevHash   string
	MerkleHash string
	StateHash  string
	CountTxs   int
}

func AddBlock(hight int, info *BlockInfo) {
	Blocks.Add(hight, BLOCK_SPAN, *info)

}

//新加内容
var (
	Blockss         = Block{BlocksInfo: make(map[int]BlockInfo)}
	BlockInfoHArray []BlockInfoh
)

//新加内容：页面展示信息
type BlockInfoh struct {
	BlockInfo
	Height int
}

//新加内容
type Block struct {
	BlocksInfo map[int]BlockInfo

	sync.RWMutex
}

//原ADD函数
func (this *Block) Add(hight int, info *BlockInfo) {
	this.Lock()
	defer this.Unlock()

	if _, ok := this.BlocksInfo[hight]; ok {
		return
	}
	this.BlocksInfo[hight] = *info
}

//新加内容
func PrintBlock() string {
	Blockss.RLock()
	defer Blockss.RUnlock()
	for k, v := range Blockss.BlocksInfo {
		One := BlockInfoh{}
		One.Height = k
		One.Hash = v.Hash
		One.PrevHash = v.PrevHash
		One.MerkleHash = v.MerkleHash
		One.StateHash = v.StateHash
		One.CountTxs = v.CountTxs
		BlockInfoHArray = append(BlockInfoHArray, One)

	}
	buf, _ := json.Marshal(BlockInfoHArray)
	result := string(buf)
	return result

}
