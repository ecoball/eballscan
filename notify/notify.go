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

package notify

import (
	"time"

	"github.com/ecoball/eballscan/data"
	"github.com/ecoball/eballscan/database"
	"github.com/ecoball/go-ecoball/common"
	"github.com/ecoball/go-ecoball/core/types"
	"github.com/ecoball/go-ecoball/spectator/info"
	"github.com/ontio/ontology/common/log"
)

func Dispatch(one info.OneNotify) {
	switch one.InfoType {
	case info.InfoBlock:
		if err := handleBlock(one.Info); nil != err {
			log.Error("handleBlock error: ", err)
		}
	default:

	}
}

func handleBlock(info []byte) error {
	oneBlock := types.Block{}
	if err := oneBlock.Deserialize(info); nil != err {
		log.Fatal(err)
		return err
	}

	//add block
	if err := database.AddBlock(int(oneBlock.Height), int(oneBlock.CountTxs), common.ToHex(oneBlock.Hash.Bytes()), common.ToHex(oneBlock.PrevHash.Bytes()),
		common.ToHex(oneBlock.MerkleHash.Bytes()), common.ToHex(oneBlock.StateHash.Bytes())); nil != err {
		log.Fatal(err)
		return err
	}

	data.AddBlock(int(oneBlock.Height), &data.BlockInfo{common.ToHex(oneBlock.Hash.Bytes()), common.ToHex(oneBlock.PrevHash.Bytes()),
		common.ToHex(oneBlock.MerkleHash.Bytes()), common.ToHex(oneBlock.StateHash.Bytes()), int(oneBlock.CountTxs)})

	//add transactions
	for _, v := range oneBlock.Transactions {
		if err := database.AddTransaction(int(v.Type), int(v.TimeStamp), int(oneBlock.Height), common.ToHex(v.Hash.Bytes()),
			v.Permission, v.From.String(), v.Addr.String()); nil != err {
			log.Fatal(err)
			return err
		}
		data.AddTransaction(common.ToHex(v.Hash.Bytes()), &data.TransactionInfo{int(v.Type), time.Unix(v.TimeStamp, 0).Format("2006-01-02 15:04:05"),
			v.Permission, v.From.String(), v.Addr.String(), int(oneBlock.Height)})
	}

	return nil
}
