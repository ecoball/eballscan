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
	//MaxHeight int
	nodes_counts int
)

func initNode() (err error) {
	// Create the "blocks" table.
	if _, err = cockroachDb.Exec(
		`create table if not exists nodes (publicKey varchar(1024), port varchar(70), adderss varchar(70), committee_blockHeight int)`); err != nil {
		log.Fatal(err)
		return err
	}

	sqlStr := "select count(0) from nodes"
	if err := cockroachDb.QueryRow(sqlStr).Scan(&nodes_counts); nil != err {
		return err
	}

	/*if _, err = cockroachDb.Exec(
		`drop table if exists nodes`); err != nil {
		log.Fatal(err)
		return
	}*/

	return
}

func AddNode(publicKey, port, adderss string, committee_blockHeight int) (err error) {
	var values string
	values = fmt.Sprintf(`('%s', '%s', '%s', %d)`,publicKey, port, adderss, committee_blockHeight)
	values = "insert into nodes(publicKey, port, adderss, committee_blockHeight) values" + values
	_, err = cockroachDb.Exec(values)
	if nil != err {
		//log.Fatal(err)
		return err
	}

	return
}

func QueryNodesByHeight(blockHeight int)([]*data.NodeInfoH, error) {
	sqlStr := fmt.Sprintf("%d", blockHeight)
	sqlStr = "select publicKey, port, adderss, committee_blockHeight from nodes where committee_blockHeight = " + sqlStr

	rows, err := cockroachDb.Query(sqlStr)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	NodeInfoH := []*data.NodeInfoH{}
	for rows.Next() {
		var (
			committee_blockHeight       int
			publicKey, port, adderss    string
		)

		if err = rows.Scan(&publicKey, &port, &adderss, &committee_blockHeight); err != nil {
			log.Fatal(err)
			break
		}

	    NodeInfoH = append(NodeInfoH, &data.NodeInfoH{data.NodeInfo{publicKey, adderss, port}, committee_blockHeight})
	}

	return NodeInfoH, nil
}

func QueryOneNode(publicKey string) (*data.NodeInfoH, error) {
	var (
		committee_blockHeight       int
		port, adderss, sqlStr    string
	)

	sqlStr = "select publicKey, port, adderss, committee_blockHeight from nodes where publicKey = '" + publicKey + "'"
	if err := cockroachDb.QueryRow(sqlStr).Scan(&publicKey, &port, &adderss, &committee_blockHeight); nil != err {
		return nil, err
	}
	return &data.NodeInfoH{data.NodeInfo{publicKey, adderss, port}, committee_blockHeight}, nil
}

func QueryNodes(index, num int) ([]*data.NodeInfoH, int, error) {
	//var rows *sql.Rows
	if 1 == index{
		sqlStr := "select count(0) from nodes"
		if err := cockroachDb.QueryRow(sqlStr).Scan(&nodes_counts); nil != err {
			return nil, -1, err
		}
	
	}

	var pageNum int
	if nodes_counts % num == 0{
		pageNum = nodes_counts/num
	}else{
		pageNum = nodes_counts/num + 1
	}

	querysql := "select publicKey, port, adderss, committee_blockHeight from nodes order by committee_blockHeight desc limit "
	querysql = querysql + strconv.Itoa(num) + " offset " + strconv.Itoa((index-1)*num)
	rows, err := cockroachDb.Query(querysql)
	if err != nil {
		log.Fatal(err)
		return nil, -1, err
	}
	defer rows.Close()

	NodeInfoH := []*data.NodeInfoH{}
	for rows.Next() {
		var (
			committee_blockHeight       int
			publicKey, port, adderss    string
		)

		if err = rows.Scan(&publicKey, &port, &adderss, &committee_blockHeight); err != nil {
			log.Fatal(err)
			break
		}

	    NodeInfoH = append(NodeInfoH, &data.NodeInfoH{data.NodeInfo{publicKey, adderss, port}, committee_blockHeight})
	}

	return NodeInfoH, pageNum, nil
}
