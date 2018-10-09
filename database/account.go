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
	"time"

	"github.com/ecoball/eballscan/data"
	"github.com/muesli/cache2go"
)

func initAccount() (err error) {
	// Create the "transactions" table.
	if _, err = cockroachDb.Exec(
		`create table if not exists accounts(name varchar(70) primary key, 
		createdAt int, updatedAt int, abi text)`); err != nil {
		log.Fatal(err)
		return
	}

	//Load the data of transactions into the cache
	var rows *sql.Rows
	rows, err = cockroachDb.Query("select name, createdAt, updatedAt, abi from accounts")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			updatedAt, createdAt     int
			abi, name string
		)

		if err = rows.Scan(&name, &createdAt, &updatedAt, &abi); err != nil {
			log.Fatal(err)
			break
		}
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

func AddAccount(txType, timeStamp, blockHeight int, hash, permission, txFrom, address string) (err error) {
	var values string
	values = fmt.Sprintf(`('%s', %d, %d, '%s', '%s', '%s', %d)`, hash, txType, timeStamp, permission, txFrom, address, blockHeight)
	values = "insert into transactions(hash, txType, timeStamp, permission, txFrom, address, blockHeight) values" + values
	_, err = cockroachDb.Exec(values)
	if nil != err {
		log.Fatal(err)
	}

	return
}

func QueryOneAccount(hash string) (*data.TransactionInfo, error) {
	var (
		txType, timeStamp, blockHight       int
		permission, txFrom, address, sqlStr string
	)

	sqlStr = "select txType, timeStamp, permission, txFrom, address, blockHight from transactions where hash = '" + hash + "'"
	if err := cockroachDb.QueryRow(sqlStr).Scan(&txType, &timeStamp, &permission, &txFrom, &address, &blockHight); nil != err {
		return nil, err
	}
	return &data.TransactionInfo{txType, time.Unix(int64(timeStamp), 0).Format("2006-01-02 15:04:05"), permission, txFrom, address, blockHight}, nil
}
