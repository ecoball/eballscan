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
	Max_ViewChange_Height int
	curr_max_view_change_height int
)

func initViewchangeblock() (err error) {
	// Create the "blocks" table.
	if _, err = cockroachDb.Exec(
		`create table if not exists viewchangeblocks (height int primary key, timeStamp int,
		hash varchar(70), prevHash varchar(70), port varchar(70), adderss varchar(70), publicKey varchar(70), Round int, CMEpochNo int, FinalBlockHeight int)`); err != nil {
		log.Fatal(err)
		return err
	}

	sqlStr := "select count(0) from viewchangeblocks"
	if err := cockroachDb.QueryRow(sqlStr).Scan(&curr_max_view_change_height); nil != err {
		return err
	}

	if curr_max_view_change_height > 0{
		sqlStr = "select max(height) from viewchangeblocks"
		if err := cockroachDb.QueryRow(sqlStr).Scan(&Max_ViewChange_Height); nil != err {
			return err
		}
	}

	/*if _, err = cockroachDb.Exec(
		`drop table if exists viewchangeblocks`); err != nil {
		log.Fatal(err)
		return
	}*/

	return
}

func AddViewchangeblock(height, timeStamp, Round, CMEpochNo, FinalBlockHeight int, hash, prevHash, port, adderss, publicKey string) (err error) {
	var values string
	values = fmt.Sprintf(`(%d, %d, '%s', '%s', '%s', '%s', '%s', %d, %d, %d)`, height, timeStamp, hash, prevHash, port, adderss, publicKey, Round, CMEpochNo, FinalBlockHeight)
	values = "insert into viewchangeblocks(height, timeStamp, hash, prevHash, port, adderss, publicKey, Round, CMEpochNo, FinalBlockHeight) values" + values
	_, err = cockroachDb.Exec(values)
	if nil != err {
		fmt.Println(err)
		return err
	}

	return
}

func QueryOneViewChangeBlock(height int) (*data.ViewChange_blockInfo, int, error) {
	var (
		max_height, timeStamp, Round, CMEpochNo, FinalBlockHeight   int
		hash, prevHash, port, address, publicKey, sqlStr string
	)

	queryStr := "select max(height) from viewchangeblocks"
	if err := cockroachDb.QueryRow(queryStr).Scan(&max_height); nil != err {
		return nil, -1, err
	}

	sqlStr = fmt.Sprintf("%d", height)
	sqlStr = "select timeStamp, hash, prevHash, port, adderss, publicKey, Round, CMEpochNo, FinalBlockHeight from viewchangeblocks where height = " + sqlStr
	if err := cockroachDb.QueryRow(sqlStr).Scan(&timeStamp, &hash, &prevHash, &port, &address, &publicKey, &Round, &CMEpochNo, &FinalBlockHeight); nil != err {
		return nil, -1, err
	}
	return &data.ViewChange_blockInfo{timeStamp/1e6, hash, prevHash, data.NodeInfo{publicKey, address, port}, Round, CMEpochNo, FinalBlockHeight}, max_height, nil
}

func QueryViewChangeBlock(index, num int) ([]*data.ViewChange_blockInfoH, int, error) {
	//var rows *sql.Rows
	if 1 == index{
		sqlStr := "select max(height) from viewchangeblocks"
		if err := cockroachDb.QueryRow(sqlStr).Scan(&curr_max_view_change_height); nil != err {
			return nil, -1, err
		}
	
	}

	var pageNum int
	if curr_max_view_change_height % num == 0{
		pageNum = curr_max_view_change_height/num
	}else{
		pageNum = curr_max_view_change_height/num + 1
	}

	querysql := "select height, timeStamp, hash, prevHash, port, adderss, publicKey, Round, CMEpochNo, FinalBlockHeight from viewchangeblocks order by timeStamp desc limit "
	querysql = querysql + strconv.Itoa(num) + " offset " + strconv.Itoa((index-1)*num)
	rows, err := cockroachDb.Query(querysql)
	if err != nil {
		log.Fatal(err)
		return nil, -1, err
	}
	defer rows.Close()

	ViewChange_blockInfoH := []*data.ViewChange_blockInfoH{}
	for rows.Next() {
		var (
			height, timeStamp, Round, CMEpochNo, FinalBlockHeight   int
			hash, prevHash, port, address, publicKey string
		)

		if err = rows.Scan(&height, &timeStamp, &hash, &prevHash, &port, &address, &publicKey, &Round, &CMEpochNo, &FinalBlockHeight); err != nil {
			log.Fatal(err)
			break
		}

		ViewChange_blockInfoH = append(ViewChange_blockInfoH, &data.ViewChange_blockInfoH{data.ViewChange_blockInfo{timeStamp/1e6, hash, prevHash, data.NodeInfo{publicKey, address, port},
			Round, CMEpochNo, FinalBlockHeight},height})
	}

	//blockinfo := data.BlockInfo{hash, prevHash, merkleHash, stateHash, countTxs, timestamp, numTransaction}
	return ViewChange_blockInfoH, pageNum, nil
}

func QueryViewChangeBlockByFinalBlockHeight(FinalBlockHeight int)([]*data.ViewChange_blockInfoH, error) {
	sqlStr := fmt.Sprintf("%d", FinalBlockHeight)
	sqlStr = "select * from viewchangeblocks where FinalBlockHeight = " + sqlStr

	rows, err := cockroachDb.Query(sqlStr)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	ViewChange_blockInfoH := []*data.ViewChange_blockInfoH{}
	for rows.Next() {
		var (
			height, timeStamp, Round, CMEpochNo, FinalBlockHeight   int
			hash, prevHash, port, address, publicKey string
		)

		if err = rows.Scan(&height, &timeStamp, &hash, &prevHash, &port, &address, &publicKey, &Round, &CMEpochNo, &FinalBlockHeight); err != nil {
			log.Fatal(err)
			break
		}

		ViewChange_blockInfoH = append(ViewChange_blockInfoH, &data.ViewChange_blockInfoH{data.ViewChange_blockInfo{timeStamp/1e6, hash, prevHash, data.NodeInfo{publicKey, address, port},
			Round, CMEpochNo, FinalBlockHeight},height})
	}

	//blockinfo := data.BlockInfo{hash, prevHash, merkleHash, stateHash, countTxs, timestamp, numTransaction}
	return ViewChange_blockInfoH, nil
}