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
	"fmt"

	"github.com/ecoball/go-ecoball/account"
	"github.com/ecoball/go-ecoball/common"
	"github.com/ecoball/go-ecoball/common/errors"
	"github.com/ecoball/go-ecoball/core/pb"
	"github.com/ecoball/go-ecoball/crypto/secp256k1"
)

const VersionTx = 1

type TxType uint32

const (
	TxDeploy   TxType = 0x01
	TxInvoke   TxType = 0x02
	TxTransfer TxType = 0x03
)

type VmType uint32

const (
	VmWasm   VmType = 0x01
	VmNative VmType = 0x02
)

type Payload interface {
	Serialize() ([]byte, error)
	Deserialize(data []byte) error
	GetObject() interface{}
	Type() uint32
	JsonString() string
	Show()
}

type Transaction struct {
	Version    uint32             `json:"version"`
	ChainID    common.Hash        `json:"chain_id"`
	Type       TxType             `json:"type"`
	From       common.AccountName `json:"from"`
	Permission string             `json:"permission"`
	Addr       common.AccountName `json:"addr"`
	Nonce      uint64             `json:"nonce"`
	TimeStamp  int64              `json:"timeStamp"`
	Payload    Payload            `json:"payload"`
	Signatures []common.Signature `json:"signatures"`
	Hash       common.Hash        `json:"hash"`
	Receipt    TransactionReceipt
}

func NewTransaction(t TxType, from, addr common.AccountName, chainID common.Hash, perm string, payload Payload, nonce uint64, time int64) (*Transaction, error) {
	if payload == nil {
		return nil, errors.New(log, "the transaction's payload is nil")
	}
	tx := Transaction{
		Version:    VersionTx,
		ChainID:    chainID,
		Type:       t,
		From:       from,
		Permission: perm,
		Addr:       addr,
		Nonce:      nonce,
		TimeStamp:  time,
		Payload:    payload,
	}
	if tx.Permission == "" {
		tx.Permission = "active"
	}
	data, err := tx.unSignatureData()
	if err != nil {
		return nil, err
	}
	tx.Hash, err = common.DoubleHash(data)
	if err != nil {
		return nil, err
	}
	tx.Receipt.Hash = tx.Hash
	return &tx, nil
}

func (t *Transaction) SetSignature(account *account.Account) error {
	sigData, err := account.Sign(t.Hash.Bytes())
	if err != nil {
		return err
	}
	sig := common.Signature{}
	sig.SigData = common.CopyBytes(sigData)
	sig.PubKey = common.CopyBytes(account.PublicKey)
	t.Signatures = append(t.Signatures, sig)
	return nil
}

func (t *Transaction) VerifySignature() (bool, error) {
	for _, v := range t.Signatures {
		result, err := secp256k1.Verify(t.Hash.Bytes(), v.SigData, v.PubKey)
		if err != nil || result != true {
			return false, err
		}
	}
	return true, nil
}

func (t *Transaction) unSignatureData() ([]byte, error) {
	payload, err := t.Payload.Serialize()
	if err != nil {
		return nil, err
	}
	p := &pb.TxPayload{
		Version:    t.Version,
		ChainID:    t.ChainID.Bytes(),
		From:       uint64(t.From),
		Permission: []byte(t.Permission),
		Addr:       uint64(t.Addr),
		Payload:    payload,
		Nonce:      t.Nonce,
		Timestamp:  t.TimeStamp,
	}
	b, err := p.Marshal()
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (t *Transaction) protoBuf() (*pb.Transaction, error) {
	payload, err := t.Payload.Serialize()
	if err != nil {
		return nil, err
	}
	var sig []*pb.Signature
	for i := 0; i < len(t.Signatures); i++ {
		s := &pb.Signature{PubKey: t.Signatures[i].PubKey, SigData: t.Signatures[i].SigData}
		sig = append(sig, s)
	}
	p := &pb.Transaction{
		Payload: &pb.TxPayload{
			Version:    t.Version,
			ChainID:    t.ChainID.Bytes(),
			Type:       uint32(t.Type),
			From:       uint64(t.From),
			Permission: []byte(t.Permission),
			Addr:       uint64(t.Addr),
			Payload:    payload,
			Nonce:      t.Nonce,
			Timestamp:  t.TimeStamp,
		},
		Sign:    sig,
		Hash:    t.Hash.Bytes(),
		Receipt: &pb.TransactionReceipt{Hash: t.Hash.Bytes(), Cpu: t.Receipt.Cpu, Net: t.Receipt.Net, Result: common.CopyBytes(t.Receipt.Result)},
	}
	return p, nil
}

/**
 *  @brief converts a structure into a sequence of characters
 *  @return []byte - a sequence of characters
 */
func (t *Transaction) Serialize() ([]byte, error) {
	p, err := t.protoBuf()
	if err != nil {
		return nil, err
	}
	b, err := p.Marshal()
	if err != nil {
		return nil, errors.New(log, fmt.Sprintf("Marshal error:%s", err.Error()))
	}
	return b, nil
}

/**
 *  @brief converts a sequence of characters into a structure
 *  @param data - a sequence of characters
 */
func (t *Transaction) Deserialize(data []byte) error {
	if len(data) == 0 {
		return errors.New(log, "input data's length is zero")
	}

	var txPb pb.Transaction
	if err := txPb.Unmarshal(data); err != nil {
		return errors.New(log, fmt.Sprintf("data len: %d, unMarshal error:%s", len(data), err.Error()))
	}

	t.Version = txPb.Payload.Version
	t.ChainID = common.NewHash(txPb.Payload.ChainID)
	t.Type = TxType(txPb.Payload.Type)
	t.From = common.AccountName(txPb.Payload.From)
	t.Permission = string(txPb.Payload.Permission)
	t.Addr = common.AccountName(txPb.Payload.Addr)
	t.Nonce = txPb.Payload.Nonce
	t.TimeStamp = txPb.Payload.Timestamp
	t.Receipt = TransactionReceipt{Hash: common.NewHash(txPb.Receipt.Hash), Cpu: txPb.Receipt.Cpu, Net: txPb.Receipt.Net, Result: common.CopyBytes(txPb.Receipt.Result)}
	if t.Payload == nil {
		switch t.Type {
		case TxTransfer:
			t.Payload = new(TransferInfo)
		case TxDeploy:
			t.Payload = new(DeployInfo)
		case TxInvoke:
			t.Payload = new(InvokeInfo)
		default:
			return errors.New(log, "the transaction's payload must not be nil")
		}
	}
	if err := t.Payload.Deserialize(txPb.Payload.Payload); err != nil {
		return err
	}
	for i := 0; i < len(txPb.Sign); i++ {
		sig := common.Signature{
			PubKey:  common.CopyBytes(txPb.Sign[i].PubKey),
			SigData: common.CopyBytes(txPb.Sign[i].SigData),
		}
		t.Signatures = append(t.Signatures, sig)
	}
	t.Hash = common.NewHash(txPb.Hash)

	return nil
}

func (t *Transaction) Show() {
	fmt.Println(t.JsonString())
	t.Payload.Show()
}

func (t *Transaction) show() {
	fmt.Println("\t---------------Transaction-------------")
	fmt.Println("\tVersion        :", t.Version)
	fmt.Println("\tFrom           :", t.From.String())
	fmt.Println("\tAddr           :", t.Addr.String())
	fmt.Println("\tTime           :", t.TimeStamp)
	fmt.Println("\tHash           :", t.Hash.HexString())
	fmt.Println("\tSig Len        :", len(t.Signatures))
	for i := 0; i < len(t.Signatures); i++ {
		fmt.Println("\tPublicKey      :", common.ToHex(t.Signatures[i].PubKey))
		fmt.Println("\tSigData        :", common.ToHex(t.Signatures[i].SigData))
	}
	t.Payload.Show()
}

func (t *Transaction) JsonString() string {
	data, _ := json.Marshal(struct {
		Version    uint32 `json:"version"`
		ChainID    string
		Type       string             `json:"type"`
		From       string             `json:"from"`
		Permission string             `json:"permission"`
		Addr       string             `json:"addr"`
		Nonce      uint64             `json:"nonce"`
		TimeStamp  int64              `json:"timeStamp"`
		Payload    string             `json:"payload"`
		Signatures []common.Signature `json:"signatures"`
		Hash       string             `json:"hash"`
		Receipt    TransactionReceipt `json:"receipt"`
	}{Version: t.Version, ChainID: t.ChainID.HexString(), Type: t.Type.String(), From: t.From.String(),
		Permission: t.Permission, Addr: t.Addr.String(), Nonce: t.Nonce,
		TimeStamp: t.TimeStamp, Payload: t.Payload.JsonString(), Signatures: t.Signatures,
		Hash: t.Hash.HexString(), Receipt: t.Receipt})
	//data, _ = json.Marshal(t)
	return string(data)
}

func charToSymbol(c byte) byte {
	if c >= 'a' && c <= 'z' {
		return c - 'a' + 6
	}
	if c >= '1' && c <= '5' {
		return c - '1' + 1
	}
	return 0
}

func stringToIndex(name string) TxType {
	var index uint64
	var i uint32
	sLen := uint32(len(name))
	for ; i <= 12; i++ {
		var c uint64
		if i < sLen {
			c = uint64(charToSymbol(name[i]))
		}
		if i < 12 {
			c &= 0x1f
			c <<= 64 - 5*(i+1)
		} else {
			c &= 0x0f
		}
		index |= c
	}
	return TxType(index)
}

func stringToBytes(name string) []byte {
	slen := uint32(len(name))
	var i uint32 = 0
	var a []byte
	for ; i < slen; i++ {
		a[i] = name[i]
	}
	return a
}

func (t *Transaction) StringJson(data string) {
	//tt := Transaction_invoke{}
	json.Unmarshal([]byte(data), t)

	/*t.Addr = common.NameToIndex(tt.Addr)
	t.From = common.NameToIndex(tt.From)
	t.Hash = common.HexToHash(tt.Hash)
	t.Nonce = tt.Nonce
	t.Permission = tt.Permission
	//t.Receipt = tt.Receipt
	t.Signatures = tt.Signatures
	t.TimeStamp = tt.TimeStamp
	t.Type = stringToIndex(tt.Type)
	t.Version = tt.Version
	t.Payload = &tt.Payload*/
}

func (t TxType) String() string {
	switch t {
	case TxDeploy:
		return "deploy transaction"
	case TxInvoke:
		return "invoke transaction"
	case TxTransfer:
		return "transfer transaction"
	default:
		return "unknown type"
	}
}
