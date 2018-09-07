// Copyright 2018 The go-ecoball Authors
// This file is part of the go-ecoball library.
//
// The go-ecoball library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ecoball library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ecoball library. If not, see <http://www.gnu.org/licenses/>.

package types

import (
	"encoding/json"
	errIn "errors"
	"fmt"
	"github.com/ecoball/go-ecoball/account"
	"github.com/ecoball/go-ecoball/common"
	"github.com/ecoball/go-ecoball/common/errors"
	"github.com/ecoball/go-ecoball/core/bloom"
	"github.com/ecoball/go-ecoball/core/pb"
	"github.com/ecoball/go-ecoball/core/trie"
)

type Block struct {
	*Header
	CountTxs     uint32
	Transactions []*Transaction
	Shards       []Shard
}

func NewBlock(chainID common.Hash, prevHeader *Header, stateHash common.Hash, headerPayload Payload, consensusData ConsensusData, txs []*Transaction, cpu, net float64, timeStamp int64) (*Block, error) {
	if nil == prevHeader {
		return nil, errors.New(log, "invalid parameter preHeader")
	}
	var Bloom bloom.Bloom
	var hashes []common.Hash
	for _, t := range txs {
		hashes = append(hashes, t.Hash)
		Bloom.Add(t.Hash.Bytes())
		Bloom.Add(common.IndexToBytes(t.From))
		Bloom.Add(common.IndexToBytes(t.Addr))
	}
	merkleHash, err := trie.GetMerkleRoot(hashes)
	if err != nil {
		return nil, err
	}

	var cpuLimit, netLimit float64
	if cpu < (BlockCpuLimit / 10) {
		cpuLimit = prevHeader.Receipt.BlockCpu * 1.01
		if cpuLimit > VirtualBlockCpuLimit {
			cpuLimit = VirtualBlockCpuLimit
		}
	} else {
		cpuLimit = prevHeader.Receipt.BlockCpu * 0.99
		if cpuLimit < BlockCpuLimit {
			cpuLimit = BlockCpuLimit
		}
	}
	if net < (BlockNetLimit / 10) {
		netLimit = prevHeader.Receipt.BlockNet * 1.01
		if netLimit > VirtualBlockNetLimit {
			netLimit = VirtualBlockNetLimit
		}
	} else {
		netLimit = prevHeader.Receipt.BlockNet * 0.99
		if netLimit < BlockNetLimit {
			netLimit = BlockNetLimit
		}
	}

	header, err := NewHeader(headerPayload, VersionHeader, chainID, prevHeader.Height+1, prevHeader.Hash, merkleHash, stateHash, consensusData, Bloom, cpuLimit, netLimit, timeStamp)
	if err != nil {
		return nil, err
	}
	block := Block{
		Header:       header,
		CountTxs:     uint32(len(txs)),
		Transactions: txs,
		Shards:       nil,
	}
	return &block, nil
}

func (b *Block) CmBlockSetData(Shards []Shard) {
	b.Shards = Shards
}

func (b *Block) SetSignature(account *account.Account) error {
	return b.Header.SetSignature(account)
}

func (b *Block) GetTransaction(hash common.Hash) (*Transaction, error) {
	for _, tx := range b.Transactions {
		if hash.Equals(&tx.Hash) {
			return tx, nil
		}
	}
	return nil, errIn.New("can't find this transaction")
}

func (b *Block) IsExistedTransaction(hash common.Hash) bool {
	for _, tx := range b.Transactions {
		if hash.Equals(&tx.Hash) {
			return true
		}
	}
	return false
}

func GenesesBlockInitConsensusData(timestamp int64) *ConsensusData {
	conData, err := InitConsensusData(timestamp)
	if err != nil {
		log.Debug(err)
		return nil
	}
	return conData
}

func (b *Block) protoBuf() (*pb.BlockTx, error) {
	var block pb.BlockTx
	var err error
	block.Header, err = b.Header.protoBuf()
	if err != nil {
		return nil, err
	}
	var pbTxs []*pb.Transaction
	for _, tx := range b.Transactions {
		pbTx, err := tx.protoBuf()
		if err != nil {
			return nil, err
		}
		pbTxs = append(pbTxs, pbTx)
	}
	block.Transactions = append(block.Transactions, pbTxs...)
	return &block, nil
}

/**
 *  @brief converts a structure into a sequence of characters
 *  @return []byte - a sequence of characters
 */
func (b *Block) Serialize() (data []byte, err error) {
	p, err := b.protoBuf()
	if err != nil {
		return nil, err
	}
	data, err = p.Marshal()
	if err != nil {
		return nil, err
	}
	return data, nil
}

/**
 *  @brief converts a sequence of characters into a structure
 *  @param data - a sequence of characters
 */
func (b *Block) Deserialize(data []byte) error {
	if len(data) == 0 {
		return errors.New(log, "input data's length is zero")
	}
	var pbBlock pb.BlockTx
	if err := pbBlock.Unmarshal(data); err != nil {
		return err
	}
	dataHeader, err := pbBlock.Header.Marshal()
	if err != nil {
		return err
	}

	b.Header = new(Header)
	err = b.Header.Deserialize(dataHeader)
	if err != nil {
		return err
	}

	var txs []*Transaction
	for _, tx := range pbBlock.Transactions {
		b, err := tx.Marshal()
		if err != nil {
			return err
		}
		t := new(Transaction)
		if err := t.Deserialize(b); err != nil {
			return err
		}
		txs = append(txs, t)
	}

	b.CountTxs = uint32(len(txs))
	b.Transactions = txs

	return nil
}

func (b *Block) Show(format bool) {
	fmt.Println(b.JsonString(format))
}

func (b *Block) JsonString(format bool) string {
	if !format {
		data, _ := json.Marshal(b)
		return string(data)
	} else {
		data := b.Header.JsonString()
		data += fmt.Sprintf("{CountTxs:%d}", b.CountTxs)
		for _, v := range b.Transactions {
			data += v.JsonString()
		}
		return string(data)
	}
}

func (b *Block) Blk2BlkTx() (*pb.BlockTx, error) {
	block, err := b.protoBuf()
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (b *Block) BlkTx2Blk(blktx pb.BlockTx) error {
	dataHeader, err := blktx.Header.Marshal()
	if err != nil {
		return err
	}
	b.Header = new(Header)
	err = b.Header.Deserialize(dataHeader)
	if err != nil {
		return err
	}
	var txs []*Transaction
	for _, tx := range blktx.Transactions {
		b, err := tx.Marshal()
		if err != nil {
			return err
		}
		t := new(Transaction)
		if err := t.Deserialize(b); err != nil {
			return err
		}
		txs = append(txs, t)
	}
	b.CountTxs = uint32(len(txs))
	b.Transactions = txs
	return nil
}
