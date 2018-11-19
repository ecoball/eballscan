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
	"os"
	"github.com/ecoball/go-ecoball/common/elog"
	_ "github.com/lib/pq"
)

var (
	cockroachDb *sql.DB
	log         = elog.NewLogger("database", elog.DebugLog)
)

func init() {
	// Connect to the "blockchain" database.
	var err error
	cockroachDb, err = sql.Open("postgres", "postgresql://eballscan@localhost:26257/blockchain?sslmode=disable")
	if err != nil {
		log.Fatal("connecting to the database error: ", err)
		os.Exit(1)
	}

	//init block
	err = initBlock()
	if err != nil {
		log.Fatal("initialize block error: ", err)
		os.Exit(1)
	}

	err = initAccount()
	if err != nil {
		log.Fatal("initialize account error: ", err)
		os.Exit(1)
	}

	err = initCommittee_block()
	if err != nil {
		log.Fatal("initialize Committee_block error: ", err)
		os.Exit(1)
	}

	err = initFinal_block()
	if err != nil {
		log.Fatal("initialize final_block error: ", err)
		os.Exit(1)
	}

	err = initMinor_block()
	if err != nil {
		log.Fatal("initialize Minor_block error: ", err)
		os.Exit(1)
	}

	err = initNode()
	if err != nil {
		log.Fatal("initialize node error: ", err)
		os.Exit(1)
	}

	err = initViewchangeblock()
	if err != nil {
		log.Fatal("initialize Viewchangeblock error: ", err)
		os.Exit(1)
	}

	//init transaction
	err = initTransaction()
	if err != nil {
		log.Fatal("initialize transaction error: ", err)
		os.Exit(1)
	}
}
