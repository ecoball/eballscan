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
	"sync"

	"github.com/ecoball/eballscan/data"
	"github.com/ecoball/eballscan/syn"
	"github.com/ecoball/go-ecoball/common/elog"
	"github.com/ecoball/go-ecoball/core/types"
	_ "github.com/lib/pq"
)

var (
	CockroachDb *sql.DB
	DbMutex     sync.Mutex
	log             = elog.NewLogger("database", elog.DebugLog)
	Tx_index    int = 0
)

func init() {
	// Connect to the "bank" database.
	var err error
	CockroachDb, err = sql.Open("postgres", "postgresql://eballscan@localhost:26257/blockchain?sslmode=disable")
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}

	// Create the "blocks" table.
	if _, err = CockroachDb.Exec(
		`create table if not exists blocks (hight int primary key, 
			hash varchar(70), prevHash varchar(70), merkleHash varchar(70), stateHash varchar(70), countTxs int)`); err != nil {
		log.Fatal(err)
	}
	if _, err = CockroachDb.Exec(
		`create table if not exists transactions (version int primary key, 
			ty int, from int, permission varchar(70), addr int, nonce int, timeStamp int)`); err != nil {
		log.Fatal(err)
	}
	// Print out the balances.
	rows, errQuery := CockroachDb.Query("select hight, hash, prevHash, merkleHash, stateHash, countTxs from blocks")
	if errQuery != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			hight, countTxs                       int
			hash, prevHash, merkleHash, stateHash string
		)

		if err := rows.Scan(&hight, &hash, &prevHash, &merkleHash, &stateHash, &countTxs); err != nil {
			log.Fatal(err)
		}

		data.Blocks.Add(hight, data.BlockInfo{hash, prevHash, merkleHash, stateHash, countTxs})

		if hight > syn.MaxHight {
			syn.MaxHight = hight
		}
	}
	txrows, txerrQuery := CockroachDb.Query("select version, ty, from, permission, addr, nonce, timeStamp from transactions")
	if txerrQuery != nil {
		log.Fatal(err)
	}
	defer txrows.Close()

	for txrows.Next() {
		var (
			version, nonce, timeStamp, ty, from, addr int
			permission                                string
		)

		if err := rows.Scan(&version, &ty, &from, &permission, &addr, &nonce, &timeStamp); err != nil {
			log.Fatal(err)
		}

		data.Txs.Add(Tx_index, data.TxInfo{version, ty, from, permission, addr, nonce, timeStamp})
		Tx_index++
	}
}

func AddBlock(hight, countTxs int, hash, prevHash, merkleHash, stateHash string, ts []*types.Transaction) error {
	DbMutex.Lock()
	defer DbMutex.Unlock()

	var values string
	values = fmt.Sprintf(`(%d, '%s', '%s', '%s', '%s', %d)`, hight, hash, prevHash, merkleHash, stateHash, countTxs)
	values = "insert into blocks(hight, hash, prevHash, merkleHash, stateHash, countTxs) values" + values
	_, err := CockroachDb.Exec(values)
	if nil != err {
		return err
	}
	data.Blocks.Add(hight, data.BlockInfo{hash, prevHash, merkleHash, stateHash, countTxs})

	var transaction string
	for _, v := range ts {

		transaction = fmt.Sprintf(`(%d, '%d', '%d', '%s', '%d', '%d','%d')`, v.Version, v.Type, v.From, v.Permission, v.Addr, v.Nonce, v.TimeStamp)
		transaction = "insert into transactions(version, ty, from, permission, addr, nonce, timeStamp) values" + values
		_, err := CockroachDb.Exec(transaction)
		if nil != err {
			return err
		}
		data.Txs.Add(Tx_index, data.TxInfo{int(v.Version), int(v.Type), int(v.From), v.Permission, int(v.Addr), int(v.Nonce), int(v.TimeStamp)})
		Tx_index++
	}
	return nil
}
