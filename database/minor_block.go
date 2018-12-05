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
	curr_max_minor_height int
)

func initMinor_block() (err error) {
	// Create the "blocks" table.
	if _, err = cockroachDb.Exec(
		`create table if not exists minor_blocks (height int, timeStamp int,
		hash varchar(70), prevHash varchar(70), TrxHashRoot varchar(70), StateDeltaHash varchar(70), 
		CMBlockHash varchar(70), ShardId int, ProposalPublicKey varchar(200), CMEpochNo int, CountTxs int, primary key(height, ShardId))`); err != nil {
		log.Fatal(err)
		return err
	}

	sqlStr := "select count(2) from minor_blocks"
	if err := cockroachDb.QueryRow(sqlStr).Scan(&curr_max_minor_height); nil != err {
		return err
	}

	/*if _, err = cockroachDb.Exec(
		`drop table if exists minor_blocks`); err != nil {
		log.Fatal(err)
		return
	}*/

	return
}

func AddMinor_block(height, timeStamp, ShardId, CMEpochNo, CountTxs int, hash, prevHash, TrxHashRoot, StateDeltaHash, CMBlockHash, ProposalPublicKey string) (err error) {
	var values string
	values = fmt.Sprintf(`(%d, %d, '%s', '%s', '%s', '%s', '%s', %d, '%s', %d, %d)`, height, timeStamp, hash, prevHash, TrxHashRoot, StateDeltaHash,
		CMBlockHash, ShardId, ProposalPublicKey, CMEpochNo, CountTxs)
	values = "insert into minor_blocks(height, timeStamp, hash, prevHash, TrxHashRoot, StateDeltaHash, CMBlockHash, ShardId, ProposalPublicKey, CMEpochNo, CountTxs) values" + values
	_, err = cockroachDb.Exec(values)
	if nil != err {
		//log.Fatal(err)
		return err
	}

	return
}

func QueryOneMinorBlock(height, shardId int) (*data.Minor_blockInfo, int, error) {
	var (
		max_height, timeStamp, CMEpochNo, CountTxs                                          int
		hash, prevHash, TrxHashRoot, StateDeltaHash, CMBlockHash, ProposalPublicKey, sqlStr string
	)

	queryStr := "select max(height) from minor_blocks"
	if err := cockroachDb.QueryRow(queryStr).Scan(&max_height); nil != err {
		return nil, -1, err
	}

	sqlStr = fmt.Sprintf("%d", height)
	shardStr := fmt.Sprintln("%d", shardId)
	sqlStr = "select timeStamp, hash, prevHash, TrxHashRoot, StateDeltaHash, CMBlockHash, ProposalPublicKey, CMEpochNo, CountTxs from minor_blocks where height=" + sqlStr
	sqlStr += " and ShardId="
	sqlStr += shardStr
	if err := cockroachDb.QueryRow(sqlStr).Scan(&timeStamp, &hash, &prevHash, &TrxHashRoot, &StateDeltaHash, &CMBlockHash, &ProposalPublicKey, &CMEpochNo, &CountTxs); nil != err {
		return nil, -1, err
	}
	return &data.Minor_blockInfo{timeStamp / 1e6, hash, prevHash, TrxHashRoot, StateDeltaHash, CMBlockHash, shardId, ProposalPublicKey, CMEpochNo, CountTxs}, max_height, nil
}

func QueryMinorBlock(index, num int) ([]*data.Minor_blockInfoH, int, error) {
	//var rows *sql.Rows
	if 1 == index {
		sqlStr := "select count(2) from minor_blocks"
		if err := cockroachDb.QueryRow(sqlStr).Scan(&curr_max_minor_height); nil != err {
			return nil, -1, err
		}

	}

	var pageNum int
	if curr_max_minor_height%num == 0 {
		pageNum = curr_max_minor_height / num
	} else {
		pageNum = curr_max_minor_height/num + 1
	}

	querysql := "select height, timeStamp, hash, prevHash, TrxHashRoot, StateDeltaHash, CMBlockHash, ShardId, ProposalPublicKey, CMEpochNo, CountTxs from minor_blocks order by timeStamp desc limit "
	querysql = querysql + strconv.Itoa(num) + " offset " + strconv.Itoa((index-1)*num)
	rows, err := cockroachDb.Query(querysql)
	if err != nil {
		log.Fatal(err)
		return nil, -1, err
	}
	defer rows.Close()

	Minor_blockInfoH := []*data.Minor_blockInfoH{}
	for rows.Next() {
		var (
			height, timeStamp, ShardId, CMEpochNo, CountTxs                             int
			hash, prevHash, TrxHashRoot, StateDeltaHash, CMBlockHash, ProposalPublicKey string
		)

		if err = rows.Scan(&height, &timeStamp, &hash, &prevHash, &TrxHashRoot, &StateDeltaHash, &CMBlockHash, &ShardId, &ProposalPublicKey, &CMEpochNo, &CountTxs); err != nil {
			log.Fatal(err)
			break
		}

		Minor_blockInfoH = append(Minor_blockInfoH, &data.Minor_blockInfoH{data.Minor_blockInfo{timeStamp / 1e6, hash, prevHash, TrxHashRoot, StateDeltaHash, CMBlockHash,
			ShardId, ProposalPublicKey, CMEpochNo, CountTxs}, height})
	}

	//blockinfo := data.BlockInfo{hash, prevHash, merkleHash, stateHash, countTxs, timestamp, numTransaction}
	return Minor_blockInfoH, pageNum, nil
}
