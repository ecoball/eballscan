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
	"sync"
	"time"

	"github.com/muesli/cache2go"
)

const (
	TRANSACTION_SPAN time.Duration = 20 * time.Second
)

var (
	Transactions = cache2go.Cache("Transactions")
)

type TransactionInfo struct {
	TxType     int
	TimeStamp  string
	Permission string
	TxFrom     string
	Address    string
	BlockHight int
}

func AddTransaction(hash string, info *TransactionInfo) {
	Transactions.Add(hash, TRANSACTION_SPAN, *info)
}

//新加内容
var (
	Transactionss         = Transaction{TxsInfo: make(map[string]TransactionInfo)}
	TransactionInfoHArray []TransactionInfoH
)

//新加内容
type TransactionInfoH struct {
	TransactionInfo
	Hash string
}

//新加内容
type Transaction struct {
	TxsInfo map[string]TransactionInfo

	sync.RWMutex
}

//新加内容
func PrintTransaction() string {
	Transactionss.RLock()
	defer Transactions.RUnlock()
	for k, v := range Transactionss.TxsInfo {
		One := TransactionInfoH{}
		One.Hash = k
		One.TxType = v.TxType
		One.TimeStamp = v.TimeStamp
		One.Permission = v.Permission
		One.TxFrom = v.TxFrom
		One.Address = v.Address
		One.BlockHight = v.BlockHight
		TransactionInfoHArray = append(TransactionInfoHArray, One)
	}
	buf, _ := json.Marshal(TransactionInfoHArray)
	result := string(buf)
	return result

}

//原
func (this *Transaction) Add(hash string, info *TransactionInfo) {
	this.Lock()
	defer this.Unlock()

	if _, ok := this.TxsInfo[hash]; ok {
		return
	}
	this.TxsInfo[hash] = *info
}
