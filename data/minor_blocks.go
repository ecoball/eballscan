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
	"time"
	//"encoding/json"
	//"fmt"
	"github.com/muesli/cache2go"

)

const (
	MINOR_BLOCK_SPAN time.Duration = 10 * time.Second
)

var (
	Minor_blocks = cache2go.Cache("Minor_blocks")

	//THashArray         []string
)


type Minor_blockInfo struct {
	TimeStamp  int
	Hash string
	PrevHash     string
	TrxHashRoot    string
	StateDeltaHash  string
	CMBlockHash    string
	ShardId       int
	ProposalPublicKey  string
	CMEpochNo       int
	CountTxs        int
}
type Minor_blockInfoH struct {
	Minor_blockInfo
	Height int
}