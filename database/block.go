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
	MaxHight int
	curr_max_hight int
)

func initBlock() (err error) {
	// Create the "blocks" table.
	if _, err = cockroachDb.Exec(
		`create table if not exists blocks (hight int primary key, timeStamp int,
			hash varchar(70), prevHash varchar(70), merkleHash varchar(70), stateHash varchar(70), countTxs int)`); err != nil {
		log.Fatal(err)
		return
	}

	sqlStr := "select count(0) from blocks"
	if err := cockroachDb.QueryRow(sqlStr).Scan(&curr_max_hight); nil != err {
		return err
	}

	/*if _, err = cockroachDb.Exec(
		`drop table if exists blocks`); err != nil {
		log.Fatal(err)
		return
	}*/

	/*if _, err = cockroachDb.Exec(
		`create sequence if not exists blocks_id_seq   
		minvalue 1  
		maxvalue 9223372036854775807  
		start 1  
		increment 1  
		cache 1;
		`); err != nil {
		log.Fatal(err)
		return
	}*/

	//Load the data of blocks into the cache
	var rows *sql.Rows
	rows, err = cockroachDb.Query("select hight, timeStamp, hash, prevHash, merkleHash, stateHash, countTxs from blocks")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			hight, countTxs, timestamp       int
			hash, prevHash, merkleHash, stateHash string
		)

		if err = rows.Scan(&hight, &timestamp, &hash, &prevHash, &merkleHash, &stateHash, &countTxs); err != nil {
			log.Fatal(err)
			break
		}

		data.AddBlock(hight, &data.BlockInfo{hash, prevHash, merkleHash, stateHash, countTxs, timestamp})

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

		val, _, err := QueryOneBlock(hight)
		if nil != err {
			return nil
		}

		item := cache2go.NewCacheItem(hight, data.BLOCK_SPAN, val)
		return item
	})

	return
}

func AddBlock(hight, countTxs, timestamp int, hash, prevHash, merkleHash, stateHash string) (err error) {
	var values string
	values = fmt.Sprintf(`(%d, %d, '%s', '%s', '%s', '%s', %d)`, hight, timestamp, hash, prevHash, merkleHash, stateHash, countTxs)
	values = "insert into blocks(hight, timeStamp, hash, prevHash, merkleHash, stateHash, countTxs) values" + values
	_, err = cockroachDb.Exec(values)
	if nil != err {
		//log.Fatal(err)
		return err
	}

	data.AddBlock(hight, &data.BlockInfo{hash, prevHash, merkleHash, stateHash, countTxs, timestamp})

	if hight > MaxHight {
		MaxHight = hight
		data.Length = hight
	}

	return
}

func QueryOneBlock(hight int) (*data.BlockInfo, int, error) {
	var (
		countTxs, timestamp,max_hight          int
		hash, prevHash, merkleHash, stateHash, sqlStr string
	)

	queryStr := "select count(0) from blocks"
	if err := cockroachDb.QueryRow(queryStr).Scan(&max_hight); nil != err {
		return nil, -1, err
	}

	sqlStr = fmt.Sprintf("%d", hight)
	sqlStr = "select timeStamp, hash, prevHash, merkleHash, stateHash, countTxs from blocks where hight = " + sqlStr
	if err := cockroachDb.QueryRow(sqlStr).Scan(&timestamp, &hash, &prevHash, &merkleHash, &stateHash, &countTxs); nil != err {
		return nil, -1, err
	}
	return &data.BlockInfo{hash, prevHash, merkleHash, stateHash, countTxs, timestamp/1e6}, max_hight, nil
}

func QueryBlock(index, num int) ([]*data.BlockInfoh, int, error) {
	//var rows *sql.Rows
	if 1 == index{
		sqlStr := "select count(0) from blocks"
		if err := cockroachDb.QueryRow(sqlStr).Scan(&curr_max_hight); nil != err {
			return nil, -1, err
		}
	
	}

	var pageNum int
	if curr_max_hight % num == 0{
		pageNum = curr_max_hight/num
	}else{
		pageNum = curr_max_hight/num + 1
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
			hight, countTxs, timestamp   int
			hash, prevHash, merkleHash, stateHash string
		)

		if err = rows.Scan(&hight, &timestamp, &hash, &prevHash, &merkleHash, &stateHash, &countTxs); err != nil {
			log.Fatal(err)
			break
		}

	    BlockInfoh = append(BlockInfoh, &data.BlockInfoh{data.BlockInfo{hash, prevHash, merkleHash, stateHash, countTxs, timestamp/1e6}, hight})
		//return &data.BlockInfoh{data.BlockInfo{hash, prevHash, merkleHash, stateHash, countTxs, timestamp, numTransaction}, hight}, nil
	}

	//blockinfo := data.BlockInfo{hash, prevHash, merkleHash, stateHash, countTxs, timestamp, numTransaction}
	//return &data.BlockInfoh{data.BlockInfo{hash, prevHash, merkleHash, stateHash, countTxs, timestamp, numTransaction}, hight}, nil
	return BlockInfoh, pageNum, nil
}
