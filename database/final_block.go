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

package database

import (
	//"database/sql"
	"fmt"
	"strconv"

	"github.com/ecoball/eballscan/data"
	//"github.com/muesli/cache2go"
)

var (
	Max_Final_Height      int
	curr_max_final_height int
)

func initFinal_block() (err error) {
	// Create the "blocks" table.
	if _, err = cockroachDb.Exec(
		`create table if not exists final_blocks (height int primary key, timeStamp int,
		hash varchar(70), prevHash varchar(70), CMBlockHash varchar(70), TrxRootHash varchar(70), 
		StateDeltaRootHash varchar(70), MinorBlocksHash varchar(70), StateHashRoot varchar(70), MinorBlockCount int, TrxCount int, ProposalPubKey varchar(512), EpochNo int)`); err != nil {
		log.Fatal(err)
		return err
	}

	sqlStr := "select count(0) from final_blocks"
	if err := cockroachDb.QueryRow(sqlStr).Scan(&curr_max_final_height); nil != err {
		return err
	}

	if curr_max_final_height > 0 {
		sqlStr = "select max(height) from final_blocks"
		if err := cockroachDb.QueryRow(sqlStr).Scan(&Max_Final_Height); nil != err {
			return err
		}
	}

	/*if _, err = cockroachDb.Exec(
		`drop table if exists final_blocks`); err != nil {
		log.Fatal(err)
		return
	}*/

	return
}

func AddFinal_block(height, timeStamp, minorBlockCount, TrxCount, EpochNo int, hash, prevHash, CMBlockHash, TrxRootHash, StateDeltaRootHash, MinorBlocksHash, StateHashRoot, ProposalPubKey string) (err error) {
	var values string
	values = fmt.Sprintf(`(%d, %d, '%s', '%s', '%s', '%s', '%s', '%s', '%s', %d, %d, '%s', %d)`, height, timeStamp, hash, prevHash, CMBlockHash, TrxRootHash,
		StateDeltaRootHash, MinorBlocksHash, StateHashRoot, minorBlockCount, TrxCount, ProposalPubKey, EpochNo)
	values = "insert into final_blocks(height, timeStamp, hash, prevHash, CMBlockHash, TrxRootHash, StateDeltaRootHash, MinorBlocksHash, StateHashRoot, MinorBlockCount, TrxCount, ProposalPubKey, EpochNo) values" + values
	_, err = cockroachDb.Exec(values)
	if nil != err {
		//log.Fatal(err)
		return err
	}

	return
}

func QueryOneFinalBlock(height int) (*data.Final_blockInfo, int, error) {
	var (
		max_height, timeStamp, minorBlockCount, trxCount, epochNo                                                            int
		hash, prevHash, CMBlockHash, TrxRootHash, StateDeltaRootHash, MinorBlocksHash, StateHashRoot, ProposalPubKey, sqlStr string
	)

	queryStr := "select max(height) from final_blocks"
	if err := cockroachDb.QueryRow(queryStr).Scan(&max_height); nil != err {
		return nil, -1, err
	}

	sqlStr = fmt.Sprintf("%d", height)
	sqlStr = "select timeStamp, hash, prevHash, CMBlockHash, TrxRootHash, StateDeltaRootHash, MinorBlocksHash, StateHashRoot, MinorBlockCount, TrxCount, ProposalPubKey, EpochNo from final_blocks where height = " + sqlStr
	if err := cockroachDb.QueryRow(sqlStr).Scan(&timeStamp, &hash, &prevHash, &CMBlockHash, &TrxRootHash, &StateDeltaRootHash, &MinorBlocksHash, &StateHashRoot, &minorBlockCount, &trxCount, &ProposalPubKey, &epochNo); nil != err {
		return nil, -1, err
	}
	return &data.Final_blockInfo{timeStamp / 1e6, hash, prevHash, CMBlockHash, TrxRootHash, StateDeltaRootHash, MinorBlocksHash, StateHashRoot, minorBlockCount, trxCount, ProposalPubKey, epochNo}, max_height, nil
}

func QueryFinalBlock(index, num int) ([]*data.Final_blockInfoH, int, error) {
	//var rows *sql.Rows
	if 1 == index {
		sqlStr := "select max(height) from final_blocks"
		if err := cockroachDb.QueryRow(sqlStr).Scan(&curr_max_final_height); nil != err {
			return nil, -1, err
		}

	}

	var pageNum int
	if curr_max_final_height%num == 0 {
		pageNum = curr_max_final_height / num
	} else {
		pageNum = curr_max_final_height/num + 1
	}

	querysql := "select height, timeStamp, hash, prevHash, CMBlockHash, TrxRootHash, StateDeltaRootHash, MinorBlocksHash, StateHashRoot, MinorBlockCount, TrxCount, ProposalPubKey, EpochNo from final_blocks order by timeStamp desc limit "
	querysql = querysql + strconv.Itoa(num) + " offset " + strconv.Itoa((index-1)*num)
	rows, err := cockroachDb.Query(querysql)
	if err != nil {
		log.Fatal(err)
		return nil, -1, err
	}
	defer rows.Close()

	Final_blockInfoH := []*data.Final_blockInfoH{}
	for rows.Next() {
		var (
			height, timeStamp, trxCount, minorBlockCount, epochNo                                                        int
			hash, prevHash, CMBlockHash, TrxRootHash, StateDeltaRootHash, MinorBlocksHash, StateHashRoot, ProposalPubKey string
		)

		if err = rows.Scan(&height, &timeStamp, &hash, &prevHash, &CMBlockHash, &TrxRootHash, &StateDeltaRootHash, &MinorBlocksHash, &StateHashRoot, &minorBlockCount, &trxCount, &ProposalPubKey, &epochNo); err != nil {
			log.Fatal(err)
			break
		}

		Final_blockInfoH = append(Final_blockInfoH, &data.Final_blockInfoH{data.Final_blockInfo{timeStamp / 1e6, hash, prevHash, CMBlockHash, TrxRootHash, StateDeltaRootHash,
			MinorBlocksHash, StateHashRoot, minorBlockCount, trxCount, ProposalPubKey, epochNo}, height})
	}

	//blockinfo := data.BlockInfo{hash, prevHash, merkleHash, stateHash, countTxs, timestamp, numTransaction}
	return Final_blockInfoH, pageNum, nil
}
