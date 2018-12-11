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

package data

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/muesli/cache2go"
)

const (
	TRANSACTION_SPAN time.Duration = 10 * time.Second
)

var (
	Transactions = cache2go.Cache("Transactions")

	THashArray []string
)

type TransactionInfo struct {
	TxType      int
	TimeStamp   string
	Permission  string
	TxFrom      string
	Address     string
	BlockHeight int
	ShardId     int
}
type TransactionInfoH struct {
	TransactionInfo
	Hash string
}

func AddTransaction(hash string, info *TransactionInfo) {
	Transactions.Add(hash, TRANSACTION_SPAN, *info)
}
func PrintTransaction() string {
	Transactions.RLock()
	defer Transactions.RUnlock()
	var TransactionInfoHArray []TransactionInfoH
	for _, hash := range THashArray {

		res, err := Transactions.Value(hash)

		if err == nil {
			One := TransactionInfoH{}
			One.Hash = hash
			One.TxType = res.Data().(*TransactionInfo).TxType
			One.TimeStamp = res.Data().(*TransactionInfo).TimeStamp
			One.Permission = res.Data().(*TransactionInfo).Permission
			One.TxFrom = res.Data().(*TransactionInfo).TxFrom
			One.Address = res.Data().(*TransactionInfo).Address
			One.BlockHeight = res.Data().(*TransactionInfo).BlockHeight
			TransactionInfoHArray = append(TransactionInfoHArray, One)
		} else {
			fmt.Println("Error retrieving value from cache:", err)
		}
	}
	buf, _ := json.Marshal(TransactionInfoHArray)
	result := string(buf)
	return result

}
