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
	Max_Committee_Height      int
	curr_max_committee_height int
)

func initCommittee_block() (err error) {
	// Create the "blocks" table.
	if _, err = cockroachDb.Exec(
		`create table if not exists committee_blocks (height int primary key, timeStamp int,
			hash varchar(70), prevHash varchar(70), shardsHash varchar(70), leaderPubKey varchar(512), port varchar(70), adderss varchar(70), publicKey varchar(1024), nonce int, nodeCounts int)`); err != nil {
		log.Fatal(err)
		return err
	}

	sqlStr := "select count(0) from committee_blocks"
	if err := cockroachDb.QueryRow(sqlStr).Scan(&curr_max_committee_height); nil != err {
		return err
	}

	if curr_max_committee_height > 0 {
		sqlStr = "select max(height) from committee_blocks"
		if err := cockroachDb.QueryRow(sqlStr).Scan(&Max_Committee_Height); nil != err {
			return err
		}
	}

	/*if _, err = cockroachDb.Exec(
		`drop table if exists committee_blocks`); err != nil {
		log.Fatal(err)
		return
	}*/

	return
}

func AddCommittee_block(height, nonce, timestamp, nodeCounts int, hash, prevHash, shardsHash, leaderPubKey, port, adderss, publicKey string) (err error) {
	var values string
	values = fmt.Sprintf(`(%d, %d, '%s', '%s', '%s', '%s', '%s', '%s', '%s', %d, %d)`, height, timestamp, hash, prevHash, shardsHash, leaderPubKey, port, adderss, publicKey, nonce, nodeCounts)
	values = "insert into committee_blocks(height, timeStamp, hash, prevHash, shardsHash, leaderPubKey, port, adderss, publicKey, nonce, nodeCounts) values" + values
	_, err = cockroachDb.Exec(values)
	if nil != err {
		//log.Fatal(err)
		fmt.Println(err)
		return err
	}

	return
}

func QueryOneCommitteeBlock(height int) (*data.Committee_blockInfo, int, error) {
	var (
		max_height, nonce, timestamp, nodeCounts                                   int
		hash, prevHash, shardsHash, leaderPubKey, publicKey, address, port, sqlStr string
	)

	queryStr := "select max(height) from committee_blocks"
	if err := cockroachDb.QueryRow(queryStr).Scan(&max_height); nil != err {
		return nil, -1, err
	}

	sqlStr = fmt.Sprintf("%d", height)
	sqlStr = "select timeStamp, hash, prevHash, shardsHash, leaderPubKey, port, adderss, publicKey, nonce, nodeCounts from committee_blocks where height = " + sqlStr
	if err := cockroachDb.QueryRow(sqlStr).Scan(&timestamp, &hash, &prevHash, &shardsHash, &leaderPubKey, &port, &address, &publicKey, &nonce, &nodeCounts); nil != err {
		return nil, -1, err
	}
	return &data.Committee_blockInfo{timestamp / 1e6, hash, prevHash, shardsHash, leaderPubKey, data.NodeInfo{publicKey, address, port}, nonce, nodeCounts}, max_height, nil
}

func QueryOneCommitteeBlockByHash(hash string) (*data.Committee_blockInfoH, int, error) {
	var (
		max_height, nonce, timestamp, nodeCounts, blockHight                 int
		prevHash, shardsHash, leaderPubKey, publicKey, address, port, sqlStr string
	)

	queryStr := "select max(height) from committee_blocks"
	if err := cockroachDb.QueryRow(queryStr).Scan(&max_height); nil != err {
		return nil, -1, err
	}

	sqlStr = "select height, timeStamp, prevHash, shardsHash, leaderPubKey, port, adderss, publicKey, nonce, nodeCounts from committee_blocks where hash = '" + hash + "'"
	if err := cockroachDb.QueryRow(sqlStr).Scan(&blockHight, &timestamp, &prevHash, &shardsHash, &leaderPubKey, &port, &address, &publicKey, &nonce, &nodeCounts); nil != err {
		return nil, -1, err
	}
	return &data.Committee_blockInfoH{data.Committee_blockInfo{timestamp / 1e6, hash, prevHash, shardsHash, leaderPubKey, data.NodeInfo{publicKey, address, port}, nonce, nodeCounts}, blockHight}, max_height, nil
}

func QueryCommitteeBlock(index, num int) ([]*data.Committee_blockInfoH, int, error) {
	//var rows *sql.Rows
	if 1 == index {
		sqlStr := "select max(height) from committee_blocks"
		if err := cockroachDb.QueryRow(sqlStr).Scan(&curr_max_committee_height); nil != err {
			return nil, -1, err
		}

	}

	var pageNum int
	if curr_max_committee_height%num == 0 {
		pageNum = curr_max_committee_height / num
	} else {
		pageNum = curr_max_committee_height/num + 1
	}

	querysql := "select * from committee_blocks order by timeStamp desc limit "
	querysql = querysql + strconv.Itoa(num) + " offset " + strconv.Itoa((index-1)*num)
	rows, err := cockroachDb.Query(querysql)
	if err != nil {
		log.Fatal(err)
		return nil, -1, err
	}
	defer rows.Close()

	Committee_blockInfos := []*data.Committee_blockInfoH{}
	for rows.Next() {
		var (
			height, nonce, timestamp, nodeCounts                               int
			hash, prevHash, shardsHash, leaderPubKey, publicKey, address, port string
		)

		if err = rows.Scan(&height, &timestamp, &hash, &prevHash, &shardsHash, &leaderPubKey, &port, &address, &publicKey, &nonce, &nodeCounts); err != nil {
			log.Fatal(err)
			break
		}

		Committee_blockInfos = append(Committee_blockInfos, &data.Committee_blockInfoH{data.Committee_blockInfo{timestamp / 1e6, hash, prevHash,
			shardsHash, leaderPubKey, data.NodeInfo{publicKey, address, port}, nonce, nodeCounts}, height})
	}

	//blockinfo := data.BlockInfo{hash, prevHash, merkleHash, stateHash, countTxs, timestamp, numTransaction}
	return Committee_blockInfos, pageNum, nil
}
