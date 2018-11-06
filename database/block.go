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
	"strconv"

	"github.com/ecoball/eballscan/data"
	"github.com/muesli/cache2go"
)

var (
	MaxHeight int
	curr_max_height int
)

func initBlock() (err error) {
	// Create the "blocks" table.
	if _, err = cockroachDb.Exec(
		`create table if not exists blocks (height int primary key, timeStamp int,
			hash varchar(70), prevHash varchar(70), merkleHash varchar(70), stateHash varchar(70), countTxs int)`); err != nil {
		log.Fatal(err)
		return
	}

	sqlStr := "select count(0) from blocks"
	if err := cockroachDb.QueryRow(sqlStr).Scan(&curr_max_height); nil != err {
		return err
	}

	/*if _, err = cockroachDb.Exec(
		`drop table if exists blocks`); err != nil {
		log.Fatal(err)
		return
	}*/

	//Load the data of blocks into the cache
	var rows *sql.Rows
	rows, err = cockroachDb.Query("select height, timeStamp, hash, prevHash, merkleHash, stateHash, countTxs from blocks")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			height, countTxs, timestamp   int
			hash, prevHash, merkleHash, stateHash string
		)

		if err = rows.Scan(&height, &timestamp, &hash, &prevHash, &merkleHash, &stateHash, &countTxs); err != nil {
			log.Fatal(err)
			break
		}

		data.AddBlock(height, &data.BlockInfo{hash, prevHash, merkleHash, stateHash, countTxs, timestamp})

		if height > MaxHeight {
			MaxHeight = height
			data.Length = height
		}
	}

	//set loader
	data.Blocks.SetDataLoader(func(key interface{}, args ...interface{}) *cache2go.CacheItem {
		height, ok := key.(int)
		if !ok {
			return nil
		}

		val, _, err := QueryOneBlock(height)
		if nil != err {
			return nil
		}

		item := cache2go.NewCacheItem(height, data.BLOCK_SPAN, val)
		return item
	})

	return
}

func AddBlock(height, countTxs, timestamp int, hash, prevHash, merkleHash, stateHash string) (err error) {
	var values string
	values = fmt.Sprintf(`(%d, %d, '%s', '%s', '%s', '%s', %d)`, height, timestamp, hash, prevHash, merkleHash, stateHash, countTxs)
	values = "insert into blocks(height, timeStamp, hash, prevHash, merkleHash, stateHash, countTxs) values" + values
	_, err = cockroachDb.Exec(values)
	if nil != err {
		//log.Fatal(err)
		return err
	}

	return
}

func QueryOneBlock(height int) (*data.BlockInfo, int, error) {
	var (
		countTxs, timestamp, max_height          int
		hash, prevHash, merkleHash, stateHash, sqlStr string
	)

	queryStr := "select max(height) from blocks"
	if err := cockroachDb.QueryRow(queryStr).Scan(&max_height); nil != err {
		return nil, -1, err
	}

	sqlStr = fmt.Sprintf("%d", height)
	sqlStr = "select timeStamp, hash, prevHash, merkleHash, stateHash, countTxs from blocks where height = " + sqlStr
	if err := cockroachDb.QueryRow(sqlStr).Scan(&timestamp, &hash, &prevHash, &merkleHash, &stateHash, &countTxs); nil != err {
		return nil, -1, err
	}
	return &data.BlockInfo{hash, prevHash, merkleHash, stateHash, countTxs, timestamp/1e6}, max_height, nil
}

func QueryBlock(index, num int) ([]*data.BlockInfoh, int, error) {
	//var rows *sql.Rows
	if 1 == index{
		sqlStr := "select max(height) from blocks"
		if err := cockroachDb.QueryRow(sqlStr).Scan(&curr_max_height); nil != err {
			return nil, -1, err
		}
	
	}

	var pageNum int
	if curr_max_height % num == 0{
		pageNum = curr_max_height/num
	}else{
		pageNum = curr_max_height/num + 1
	}

	querysql := "select * from blocks order by timeStamp desc limit "
	querysql = querysql + strconv.Itoa(num) + " offset " + strconv.Itoa((index-1)*num)
	rows, err := cockroachDb.Query(querysql)
	if err != nil {
		log.Fatal(err)
		return nil, -1, err
	}
	defer rows.Close()

	BlockInfoh := []*data.BlockInfoh{}
	for rows.Next() {
		var (
			height, countTxs, timestamp   int
			hash, prevHash, merkleHash, stateHash string
		)

		if err = rows.Scan(&height, &timestamp, &hash, &prevHash, &merkleHash, &stateHash, &countTxs); err != nil {
			log.Fatal(err)
			break
		}

	    BlockInfoh = append(BlockInfoh, &data.BlockInfoh{data.BlockInfo{hash, prevHash, merkleHash, stateHash, countTxs, timestamp/1e6}, height})
	}

	//blockinfo := data.BlockInfo{hash, prevHash, merkleHash, stateHash, countTxs, timestamp, numTransaction}
	return BlockInfoh, pageNum, nil
}
