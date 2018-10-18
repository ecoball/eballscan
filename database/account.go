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
	"strconv"

	"github.com/ecoball/eballscan/data"
	"github.com/muesli/cache2go"
)

func initAccount() (err error) {
	// Create the "accounts" table.
	if _, err = cockroachDb.Exec(
		`create table if not exists accounts(name varchar(70) primary key, 
		timeStamp int, balance int, token varchar(32))`); err != nil {
		log.Fatal(err)
		return
	}

	/*if _, err = cockroachDb.Exec(
		`drop table if exists accounts`); err != nil {
		log.Fatal(err)
		return
	}*/

	var rows *sql.Rows
	rows, err = cockroachDb.Query("select name, timeStamp, balance, token from accounts")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer rows.Close()

	root_is_exist := false
	for rows.Next() {
		var (
			timestamp, balance       int
			name,token            string
		)

		if err = rows.Scan(&name, &timestamp, &balance, &token); err != nil {
			log.Fatal(err)
			break
		}

		if "root" == name {
			root_is_exist = true
		}

		data.AddAccount(name, &data.AccountInfo{timestamp, balance, token})
	}

	if !root_is_exist {//自动创建root账户
		timeStamp := time.Now().UnixNano()
		if err = AddAccount("root", "ABA", int(timeStamp), 70000); err != nil {
			log.Fatal(err)
		}
	}

	data.Accounts.SetDataLoader(func(key interface{}, args ...interface{}) *cache2go.CacheItem {
		name, ok := key.(string)
		if !ok {
			return nil
		}

		val, err := QueryOneAccount(name)
		if nil != err {
			return nil
		}

		item := cache2go.NewCacheItem(name, data.ACCOUNT_SPAN, *val)
		return item
	})

	return
}

func AddAccount(name, token string, timeStamp, balance int)(err error) {
	var values string
	values = fmt.Sprintf(`('%s', %d, %d, '%s')`, name, timeStamp, balance, token)
	values = "insert into accounts(name, timeStamp, balance, token) values" + values
	_, err = cockroachDb.Exec(values)
	if nil != err {
		return err
	}

	data.AddAccount(name, &data.AccountInfo{timeStamp, balance, token})

	return nil
}

func QueryOneAccount(name string) (*data.AccountInfo, error) {
	var (
		timeStamp, balance       int
		token                    string
	)

	sqlStr := "select timeStamp, balance, token from accounts where name = '" + name + "'"
	if err := cockroachDb.QueryRow(sqlStr).Scan(&timeStamp, &balance, &token); nil != err {
		return nil, err
	}
	return &data.AccountInfo{timeStamp/1e6, balance, token}, nil
}

func QueryAccountBalance(name string) (int, error) {
	var balance       int

	sqlStr := "select balance from accounts where name = '" + name + "'"
	if err := cockroachDb.QueryRow(sqlStr).Scan(&balance); nil != err {
		return -1, err
	}
	return balance, nil
}

func UpdateAccountBalance(name string, balance int) error {
	sqlStr := "update accounts set balance = " + strconv.Itoa(balance) + " where name = '" + name + "'"
	_, err := cockroachDb.Exec(sqlStr)
	if nil != err {
		return err
	}
	return nil
}



func QueryAccounts(num, index int) ([]*data.AccountInfoh, int, error) {
	var pageNum, counts int
	sqlStr := "select count(0) from accounts"
	if err := cockroachDb.QueryRow(sqlStr).Scan(&counts); nil != err {
		return nil, -1, err
	}
	
	if counts % num == 0{
		pageNum = counts/num
	}else{
		pageNum = counts/num + 1
	}


	querySql := "select * from accounts order by timeStamp desc limit " + strconv.Itoa(num) + " offset " + strconv.Itoa((index-1)*num)

	rows, err := cockroachDb.Query(querySql)
	if err != nil {
		log.Fatal(err)
		return nil, -1, err
	}
	defer rows.Close()

	accounts := []*data.AccountInfoh{}
	for rows.Next() {
		var (
			timeStamp, balance       int
			name, token 			string
		)

		if err = rows.Scan(&name, &timeStamp, &balance, &token); err != nil {
			log.Fatal(err)
			break
		}

	    accounts = append(accounts, &data.AccountInfoh{data.AccountInfo{timeStamp/1e6, balance, token}, name})
	}

	return accounts, pageNum, nil
}
