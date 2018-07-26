package database

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/ecoball/go-ecoball/common/elog"
	"github.com/ecoball/go-ecoball/explorer/data"
	_ "github.com/lib/pq"
)

var (
	CockroachDb *sql.DB
	DbMutex     sync.Mutex
	log         = elog.NewLogger("database", elog.DebugLog)
	MaxHight    int
)

func init() {
	// Connect to the "bank" database.
	var err error
	CockroachDb, err = sql.Open("postgres", "postgresql://ecoball@localhost:26257/blockchain?sslmode=disable")
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}

	// Create the "blocks" table.
	if _, err = CockroachDb.Exec(
		`create table if not exists blocks (hight int primary key, 
			hash varchar(70), prevHash varchar(70), merkleHash varchar(70), stateHash varchar(70), countTxs int)`); err != nil {
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

		if hight > MaxHight {
			MaxHight = hight
		}
	}
}

func AddBlock(hight, countTxs int, hash, prevHash, merkleHash, stateHash string) error {
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
	return nil
}
