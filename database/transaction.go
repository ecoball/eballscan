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
	"time"

	"github.com/ecoball/eballscan/data"
	"github.com/muesli/cache2go"
)

var current_transactions_num int

func initTransaction() (err error) {
	// Create the "transactions" table.
	if _, err = cockroachDb.Exec(
		`create table if not exists transactions(hash varchar(70) primary key, 
		txType int, timeStamp int, permission varchar(32), txFrom varchar(32), address varchar(32), blockHeight int, ShardId int,
		foreign key(blockHeight, ShardId) references minor_blocks(height, ShardId))`); err != nil {
		log.Fatal(err)
		return
	}

	/*if _, err = cockroachDb.Exec(
		`drop table if exists transactions`); err != nil {
		log.Fatal(err)
		return
	}*/

	sqlStr := "select count(0) from transactions"
	if err := cockroachDb.QueryRow(sqlStr).Scan(&current_transactions_num); nil != err {
		return err
	}

	//Load the data of transactions into the cache
	var rows *sql.Rows
	rows, err = cockroachDb.Query("select hash, txType, timeStamp, permission, txFrom, address, blockHeight, ShardId from transactions")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			txType, timeStamp, blockHeight, shardId int
			hash, permission, txFrom, address       string
		)

		if err = rows.Scan(&hash, &txType, &timeStamp, &permission, &txFrom, &address, &blockHeight, &shardId); err != nil {
			log.Fatal(err)
			break
		}
		data.THashArray = append(data.THashArray, hash)
		data.AddTransaction(hash, &data.TransactionInfo{txType, time.Unix(int64(timeStamp/1e9), 0).Format("2006-01-02 15:04:05"), permission, txFrom, address, blockHeight, shardId})
	}

	//set loader
	data.Transactions.SetDataLoader(func(key interface{}, args ...interface{}) *cache2go.CacheItem {
		hash, ok := key.(string)
		if !ok {
			return nil
		}

		val, err := QueryOneTransaction(hash)
		if nil != err {
			return nil
		}

		item := cache2go.NewCacheItem(hash, data.TRANSACTION_SPAN, *val)
		return item
	})

	return
}

func AddTransaction(txType, timeStamp, blockHeight, ShardId int, hash, permission, txFrom, address string) (err error) {
	var values string
	values = fmt.Sprintf(`('%s', %d, %d, '%s', '%s', '%s', %d, %d)`, hash, txType, timeStamp, permission, txFrom, address, blockHeight, ShardId)
	values = "insert into transactions(hash, txType, timeStamp, permission, txFrom, address, blockHeight, ShardId) values" + values
	_, err = cockroachDb.Exec(values)
	if nil != err {
		return err
	}

	return
}

func QueryOneTransaction(hash string) (*data.TransactionInfo, error) {
	var (
		txType, timeStamp, blockHeight, shardId int
		permission, txFrom, address, sqlStr     string
	)

	sqlStr = "select txType, timeStamp, permission, txFrom, address, blockHeight, ShardId from transactions where hash = '" + hash + "'"
	if err := cockroachDb.QueryRow(sqlStr).Scan(&txType, &timeStamp, &permission, &txFrom, &address, &blockHeight, &shardId); nil != err {
		return nil, err
	}
	return &data.TransactionInfo{txType, strconv.Itoa(timeStamp / 1e6), permission, txFrom, address, blockHeight, shardId}, nil
}

func QueryTransactionsByAccountName(num, index int, name string) ([]*data.TransactionInfoH, int, int, error) {
	var pageNum, counts int
	sqlStr := "select count(0) from transactions where txFrom = '"
	sqlStr = sqlStr + name + "' or address = '" + name + "'"
	if err := cockroachDb.QueryRow(sqlStr).Scan(&counts); nil != err {
		return nil, -1, -1, err
	}

	if counts%num == 0 {
		pageNum = counts / num
	} else {
		pageNum = counts/num + 1
	}

	querySql := "select hash, txType, timeStamp, permission, txFrom, address, blockHeight, ShardId from transactions where txFrom = '"
	querySql = querySql + name + "' or address = '" + name + "' order by timeStamp desc limit " + strconv.Itoa(num) + " offset " + strconv.Itoa((index-1)*num)

	rows, err := cockroachDb.Query(querySql)
	if err != nil {
		log.Fatal(err)
		return nil, -1, -1, err
	}
	defer rows.Close()

	transactionInfoH := []*data.TransactionInfoH{}
	for rows.Next() {
		var (
			txType, blockHeight, timeStamp, shardId int
			permission, txFrom, address, hash       string
		)

		if err = rows.Scan(&hash, &txType, &timeStamp, &permission, &txFrom, &address, &blockHeight, &shardId); err != nil {
			log.Fatal(err)
			break
		}

		transactionInfoH = append(transactionInfoH, &data.TransactionInfoH{data.TransactionInfo{txType, strconv.Itoa(timeStamp / 1e6), permission, txFrom, address, blockHeight, shardId}, hash})
	}

	return transactionInfoH, pageNum, counts, nil
}

func QueryTransactionsByHeightAndShardId(blockHeight, shardId int) ([]*data.TransactionInfoH, error) {
	sqlStr := fmt.Sprintf("%d", blockHeight)
	sqlStr = "select hash, txType, timeStamp, permission, txFrom, address, blockHeight from transactions where blockHeight = " + sqlStr + " and ShardId=" + fmt.Sprintf("%d", shardId)

	rows, err := cockroachDb.Query(sqlStr)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	transactionInfoH := []*data.TransactionInfoH{}
	for rows.Next() {
		var (
			txType, blockHeight, timeStamp    int
			permission, txFrom, address, hash string
		)

		if err = rows.Scan(&hash, &txType, &timeStamp, &permission, &txFrom, &address, &blockHeight); err != nil {
			log.Fatal(err)
			break
		}

		transactionInfoH = append(transactionInfoH, &data.TransactionInfoH{data.TransactionInfo{txType, strconv.Itoa(timeStamp / 1e6), permission, txFrom, address, blockHeight, shardId}, hash})
	}

	return transactionInfoH, nil
}

func QueryTransaction(index, num int) ([]*data.TransactionInfoH, int, error) {
	if 1 == index {
		sqlStr := "select count(0) from transactions"
		if err := cockroachDb.QueryRow(sqlStr).Scan(&current_transactions_num); nil != err {
			return nil, -1, err
		}

	}

	var pageNum int
	if current_transactions_num%num == 0 {
		pageNum = current_transactions_num / num
	} else {
		pageNum = current_transactions_num/num + 1
	}

	sqlStr := "select hash, txType, timeStamp, permission, txFrom, address, blockHeight, ShardId from transactions order by timeStamp desc limit "
	sqlStr = sqlStr + strconv.Itoa(num) + " offset " + strconv.Itoa((index-1)*num)

	rows, err := cockroachDb.Query(sqlStr)
	if err != nil {
		log.Fatal(err)
		return nil, -1, err
	}
	defer rows.Close()

	transactionInfoH := []*data.TransactionInfoH{}
	for rows.Next() {
		var (
			txType, blockHeight, timeStamp, shardId int
			permission, txFrom, address, hash       string
		)

		if err = rows.Scan(&hash, &txType, &timeStamp, &permission, &txFrom, &address, &blockHeight, &shardId); err != nil {
			log.Fatal(err)
			break
		}

		transactionInfoH = append(transactionInfoH, &data.TransactionInfoH{data.TransactionInfo{txType, strconv.Itoa(timeStamp / 1e6), permission, txFrom, address, blockHeight, shardId}, hash})
	}

	return transactionInfoH, pageNum, nil
}
