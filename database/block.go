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
	"database/sql"
	"fmt"

	"github.com/ecoball/eballscan/data"
	"github.com/muesli/cache2go"
)

var (
	MaxHight int
)

func initBlock() (err error) {
	// Create the "blocks" table.
	if _, err = cockroachDb.Exec(
		`create table if not exists blocks (hight int primary key, timeStamp int, numTransaction int,
			hash varchar(70), prevHash varchar(70), merkleHash varchar(70), stateHash varchar(70), countTxs int)`); err != nil {
		log.Fatal(err)
		return
	}

	/*if _, err = cockroachDb.Exec(
		`drop table if exists blocks`); err != nil {
		log.Fatal(err)
		return
	}*/

	//Load the data of blocks into the cache
	var rows *sql.Rows
	rows, err = cockroachDb.Query("select hight, timeStamp, numTransaction, hash, prevHash, merkleHash, stateHash, countTxs from blocks")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			hight, countTxs, numTransaction       int
			hash, prevHash, merkleHash, stateHash string
			timestamp int
		)

		if err = rows.Scan(&hight, &timestamp, &numTransaction, &hash, &prevHash, &merkleHash, &stateHash, &countTxs); err != nil {
			log.Fatal(err)
			break
		}

		data.AddBlock(hight, &data.BlockInfo{hash, prevHash, merkleHash, stateHash, countTxs, timestamp, numTransaction})

		if hight > MaxHight {
			MaxHight = hight
			data.Length = hight
		}
	}

	//set loader
	data.Blocks.SetDataLoader(func(key interface{}, args ...interface{}) *cache2go.CacheItem {
		hight, ok := key.(int)
		if !ok {
			return nil
		}

		val, err := QueryOneBlock(hight)
		if nil != err {
			return nil
		}

		item := cache2go.NewCacheItem(hight, data.BLOCK_SPAN, val)
		return item
	})

	return
}

func AddBlock(hight, countTxs, Timestamp, NumTransaction int, hash, prevHash, merkleHash, stateHash string) (err error) {
	var values string
	values = fmt.Sprintf(`(%d, %d, %d, '%s', '%s', '%s', '%s', %d)`, hight, Timestamp, NumTransaction, hash, prevHash, merkleHash, stateHash, countTxs)
	values = "insert into blocks(hight, timeStamp, numTransaction, hash, prevHash, merkleHash, stateHash, countTxs) values" + values
	_, err = cockroachDb.Exec(values)
	if nil != err {
		log.Fatal(err)
	}

	return
}

func QueryOneBlock(hight int) (*data.BlockInfo, error) {
	var (
		countTxs, numTransaction, timestamp           int
		hash, prevHash, merkleHash, stateHash, sqlStr string
	)

	sqlStr = fmt.Sprintf("%d", hight)
	sqlStr = "select timeStamp, numTransaction, hash, prevHash, merkleHash, stateHash, countTxs from blocks where hight = " + sqlStr
	if err := cockroachDb.QueryRow(sqlStr).Scan(&timestamp, &numTransaction, &hash, &prevHash, &merkleHash, &stateHash, &countTxs); nil != err {
		return nil, err
	}
	return &data.BlockInfo{hash, prevHash, merkleHash, stateHash, countTxs, timestamp, numTransaction}, nil
}
