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
	"strconv"
	"time"

	"github.com/ecoball/eballscan/data"
	"github.com/ecoball/eballscan/database"
	"github.com/ecoball/go-ecoball/common"
	"github.com/ecoball/go-ecoball/common/elog"
	"github.com/ecoball/go-ecoball/core/shard"
	"github.com/ecoball/go-ecoball/core/types"
	"github.com/ecoball/go-ecoball/spectator/info"
)

var (
	log = elog.NewLogger("notify", elog.DebugLog)
)

func Dispatch(one info.OneNotify) {
	switch one.InfoType {
	case info.InfoBlock:
		switch one.BlockType {
		case 0:
			if err := handleBlock(one.Info); nil != err {
				log.Error("handleBlock error: ", err)
			}
			break
		default:

		}
	/*case info.ShardBlock:
		switch one.BlockType {
		case mpb.Identify_APP_MSG_CM_BLOCK:
			if err := handleCommittee_block(one.Info); nil != err {
				log.Error("handleCommittee_block error: ", err)
			}
			break
		case mpb.Identify_APP_MSG_FINAL_BLOCK:
			if err := handleFinal_block(one.Info); nil != err {
				log.Error("handleFinal_block error: ", err)
			}
			break
		case mpb.Identify_APP_MSG_MINOR_BLOCK:
			if err := handleMinor_block(one.Info); nil != err {
				log.Error("handleMinor_block error: ", err)
			}
			break
		case mpb.Identify_APP_MSG_VC_BLOCK:
			if err := handleViewchangeblock(one.Info); nil != err {
				log.Error("handleViewchangeblock error: ", err)
			}
			break
		default:

		}
	*/
	default:

	}
}

func handleCommittee_block(info []byte) error {
	oneBlock := shard.CMBlock{}
	if err := oneBlock.Deserialize(info); nil != err {
		log.Fatal(err)
	}

	var nodeCounts int = 0
	for _, v := range oneBlock.Shards {
		nodeCounts += len(v.Member)
	}

	//add Committee_blocks
	err := database.AddCommittee_block(int(oneBlock.Height), int(oneBlock.Nonce), int(oneBlock.Timestamp), nodeCounts, oneBlock.Hash().HexString(), oneBlock.PrevHash.HexString(),
		oneBlock.ShardsHash.HexString(), common.ToHex(oneBlock.LeaderPubKey), oneBlock.Candidate.Port, oneBlock.Candidate.Address, common.ToHex(oneBlock.Candidate.PublicKey))
	if err != nil {
		log.Error("AddCommittee_block error: ", err)
		return err
	}

	//add nodes
	for _, v := range oneBlock.Shards {
		for _, vv := range v.Member {
			if err := database.AddNode(common.ToHex(vv.PublicKey), vv.Port, vv.Address, int(oneBlock.Height)); nil != err {
				log.Error("AddNode error: ", err)
				return err
			}
		}
	}

	return nil
}

func handleFinal_block(info []byte) error {
	oneBlock := shard.FinalBlock{}
	if err := oneBlock.Deserialize(info); nil != err {
		log.Fatal(err)
	}

	//add final block
	if err := database.AddFinal_block(int(oneBlock.Height), int(oneBlock.Timestamp), len(oneBlock.MinorBlocks), int(oneBlock.TrxCount), int(oneBlock.EpochNo), oneBlock.Hash().HexString(),
		oneBlock.PrevHash.HexString(), oneBlock.CMBlockHash.HexString(), oneBlock.TrxRootHash.HexString(), oneBlock.StateDeltaRootHash.HexString(),
		oneBlock.MinorBlocksHash.HexString(), oneBlock.StateHashRoot.HexString(), common.ToHex(oneBlock.ProposalPubKey)); nil != err {
		return err
	}

	//add minor block
	for _, oneMinorBlock := range oneBlock.MinorBlocks {
		if err := database.AddMinor_block(int(oneMinorBlock.Height), int(oneMinorBlock.Timestamp), int(oneMinorBlock.ShardId), int(oneBlock.Height), int(oneMinorBlock.CMEpochNo), -1, oneMinorBlock.Hash().HexString(),
			oneMinorBlock.PrevHash.HexString(), oneMinorBlock.TrxHashRoot.HexString(), oneMinorBlock.StateDeltaHash.HexString(), oneMinorBlock.CMBlockHash.HexString(),
			common.ToHex(oneMinorBlock.ProposalPublicKey)); nil != err {
			return err
		}
	}

	return nil
}

func handleMinor_block(info []byte) error {
	oneBlock := shard.MinorBlock{}
	if err := oneBlock.Deserialize(info); nil != err {
		log.Fatal(err)
	}

	//add minor block
	if err := database.AddMinor_block(int(oneBlock.Height), int(oneBlock.Timestamp), int(oneBlock.ShardId), -1, int(oneBlock.CMEpochNo), len(oneBlock.Transactions), oneBlock.Hash().HexString(),
		oneBlock.PrevHash.HexString(), oneBlock.TrxHashRoot.HexString(), oneBlock.StateDeltaHash.HexString(), oneBlock.CMBlockHash.HexString(),
		common.ToHex(oneBlock.ProposalPublicKey)); nil != err {
		return err
	}

	//add transaction
	if err := handleTransaction(oneBlock.Transactions, oneBlock.MinorBlockHeader.GetHeight(), oneBlock.MinorBlockHeader.ShardId); nil != err {
		return err
	}

	return nil
}

func handleViewchangeblock(info []byte) error {
	oneBlock := shard.ViewChangeBlock{}
	if err := oneBlock.Deserialize(info); nil != err {
		log.Fatal(err)
	}

	if err := database.AddViewchangeblock(int(oneBlock.Height), int(oneBlock.Timestamp), int(oneBlock.Round), int(oneBlock.CMEpochNo), int(oneBlock.FinalBlockHeight),
		oneBlock.Hash().HexString(), oneBlock.PrevHash.HexString(), oneBlock.Candidate.Port,
		oneBlock.Candidate.Address, common.ToHex(oneBlock.Candidate.PublicKey)); nil != err {
		return err
	}

	return nil
}

func handleBlock(info []byte) error {
	oneBlock := types.Block{}
	if err := oneBlock.Deserialize(info); nil != err {
		log.Fatal(err)
		return err
	}

	//add block
	if err := database.AddBlock(int(oneBlock.Height), int(oneBlock.CountTxs), int(oneBlock.TimeStamp), common.ToHex(oneBlock.Hash.Bytes()), common.ToHex(oneBlock.PrevHash.Bytes()),
		common.ToHex(oneBlock.MerkleHash.Bytes()), common.ToHex(oneBlock.StateHash.Bytes())); nil != err {
		log.Fatal(err)
		return err
	}

	data.AddBlock(int(oneBlock.Height), &data.BlockInfo{common.ToHex(oneBlock.Hash.Bytes()), common.ToHex(oneBlock.PrevHash.Bytes()),
		common.ToHex(oneBlock.MerkleHash.Bytes()), common.ToHex(oneBlock.StateHash.Bytes()), int(oneBlock.CountTxs), int(oneBlock.TimeStamp)})

	//add transactions
	if err := handleTransaction(oneBlock.Transactions, oneBlock.Height, 0); nil != err {
		return err
	}

	return nil
}

func handleTransaction(trxs []*types.Transaction, height uint64, ShardId uint32) error {
	//add transactions
	for _, v := range trxs {
		if err := database.AddTransaction(int(v.Type), int(v.TimeStamp), int(height), int(ShardId), common.ToHex(v.Hash.Bytes()),
			v.Permission, v.From.String(), v.Addr.String()); nil != err {
			log.Fatal(err)
			return err
		}
		data.AddTransaction(common.ToHex(v.Hash.Bytes()), &data.TransactionInfo{int(v.Type), time.Unix(v.TimeStamp/1000000000, 0).Format("2006-01-02 15:04:05"),
			v.Permission, v.From.String(), v.Addr.String(), int(height), int(ShardId)})

		if v.Type == 0x02 { //新增账号交易处理
			info := new(types.InvokeInfo)
			data, err := v.Payload.Serialize()
			if err != nil {
				log.Info(err)
				return err
			}

			err = info.Deserialize(data)
			if err != nil {
				log.Info(err)
				return err
			}

			if string(info.Method) == "new_account" {
				if err := database.AddAccount(info.Param[0], "ABA", int(v.TimeStamp), 0); nil != err {
					log.Fatal(err)
					return err
				}

			}
		}

		if v.Type == 0x03 { //转账交易处理
			info := new(types.TransferInfo)
			data, err := v.Payload.Serialize()
			if err != nil {
				log.Info(err)
				return err
			}

			err = info.Deserialize(data)
			if err != nil {
				log.Info(err)
				return err
			}

			amount, err := strconv.Atoi(info.Value.String())
			if err != nil {
				log.Info(err)
				return err
			}

			//from账户余额处理
			from := v.From.String()
			from_balance, err := database.QueryAccountBalance(from)
			if err != nil {
				log.Fatal(err)
				return err
			}
			balance := from_balance - amount
			err = database.UpdateAccountBalance(from, balance)
			if err != nil {
				log.Fatal(err)
				return err
			}

			//to账户余额处理
			to := v.Addr.String()
			to_balance, err := database.QueryAccountBalance(to)
			if err != nil {
				log.Fatal(err)
				return err
			}
			balance = to_balance + amount //to账户余额+
			err = database.UpdateAccountBalance(to, balance)
			if err != nil {
				log.Fatal(err)
				return err
			}
		}
	}

	return nil
}
