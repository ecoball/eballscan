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
	"github.com/ecoball/go-ecoball/common/elog"
	"github.com/ecoball/go-ecoball/common/errors"
	"github.com/ecoball/go-ecoball/core/bloom"
	"github.com/ecoball/go-ecoball/core/pb"
	"github.com/ecoball/go-ecoball/crypto/secp256k1"
)

const VersionHeader = 1

type HeaderType uint32

const (
	HeMinorBlock HeaderType = 1
	HeCmBlock    HeaderType = 2
)

type Header struct {
	Type          HeaderType
	Payload       Payload
	Version       uint32
	ChainID       common.Hash
	TimeStamp     int64
	Height        uint64
	ConsensusData ConsensusData
	PrevHash      common.Hash
	MerkleHash    common.Hash
	StateHash     common.Hash
	Bloom         bloom.Bloom

	Receipt    BlockReceipt
	Signatures []common.Signature
	Hash       common.Hash
}

var log = elog.NewLogger("LedgerImpl", elog.DebugLog)

/**
 *  @brief create a new block header, the compute this header's hash
 *  @param version - the version of header, default 1
 *  @param height - the height of this block
 *  @param prevHash - the hash of perv block
 *  @param merkleHash - the merkle hash root of transactions' hash
 *  @param stateHash - the mpt hash root of state
 *  @param conData - the data of consensus module
 *  @param bloom - the bloom filter of transactions
 *  @param timeStamp - the timeStamp of block, unit is ns
 */
func NewHeader(payload Payload, version uint32, chainID common.Hash, height uint64, prevHash, merkleHash, stateHash common.Hash, conData ConsensusData, bloom bloom.Bloom, cpuLimit, netLimit float64, timeStamp int64) (*Header, error) {
	if version != VersionHeader {
		return nil, errors.New(log, "version mismatch")
	}
	if payload == nil {
		return nil, errors.New(log, "header's payload is nil")
	}
	if conData.Payload == nil {
		return nil, errors.New(log, "consensus' payload is nil")
	}
	header := Header{
		Type:          HeaderType(payload.Type()),
		Payload:       payload,
		ChainID:       chainID,
		Version:       version,
		TimeStamp:     timeStamp,
		Height:        height,
		ConsensusData: conData,
		PrevHash:      prevHash,
		MerkleHash:    merkleHash,
		StateHash:     stateHash,
		Bloom:         bloom,
		Receipt:       BlockReceipt{BlockCpu: cpuLimit, BlockNet: netLimit},
		Signatures:    nil,
		Hash:          common.Hash{},
	}
	data, err := header.unSignatureData()
	if err != nil {
		return nil, err
	}
	b, err := data.Marshal()
	if err != nil {
		return nil, err
	}
	header.Hash, err = common.DoubleHash(b)
	if err != nil {
		return nil, err
	}
	return &header, nil
}

func (h *Header) InitializeHash() error {
	if h.Version != VersionHeader {
		return errors.New(log, "version mismatch")
	}
	if h.ConsensusData.Payload == nil {
		return errors.New(log, "consensus' payload is nil")
	}
	payload, err := h.unSignatureData()
	if err != nil {
		return err
	}
	b, err := payload.Marshal()
	if err != nil {
		return err
	}
	h.Hash, err = common.DoubleHash(b)
	fmt.Println("New Header Hash:", h.Hash.HexString())
	if err != nil {
		return err
	}
	return nil
}

func (h *Header) SetSignature(account *account.Account) error {
	sigData, err := account.Sign(h.Hash.Bytes())
	if err != nil {
		return err
	}
	sig := common.Signature{}
	sig.SigData = common.CopyBytes(sigData)
	sig.PubKey = common.CopyBytes(account.PublicKey)
	h.Signatures = append(h.Signatures, sig)
	return nil
}

func (h *Header) VerifySignature() (bool, error) {
	for _, v := range h.Signatures {
		b, err := secp256k1.Verify(h.Hash.Bytes(), v.SigData, v.PubKey)
		if err != nil || b != true {
			return false, err
		}
	}
	return true, nil
}

/**
** Used to compute hash
 */
func (h *Header) unSignatureData() (*pb.Header, error) {
	if h.TimeStamp == 0 {
		return nil, errors.New(log, "this header struct is illegal")
	}
	pbCon, err := h.ConsensusData.ProtoBuf()
	if err != nil {
		return nil, err
	}
	return &pb.Header{
		Version:       h.Version,
		ChainID:       h.ChainID.Bytes(),
		Timestamp:     h.TimeStamp,
		Height:        h.Height,
		ConsensusData: pbCon,
		PrevHash:      h.PrevHash.Bytes(),
		MerkleHash:    h.MerkleHash.Bytes(),
		StateHash:     h.StateHash.Bytes(),
		Bloom:         h.Bloom.Bytes(),
	}, nil
}

func (h *Header) protoBuf() (*pb.HeaderTx, error) {
	var sig []*pb.Signature
	for i := 0; i < len(h.Signatures); i++ {
		s := &pb.Signature{PubKey: h.Signatures[i].PubKey, SigData: h.Signatures[i].SigData}
		sig = append(sig, s)
	}
	pbCon, err := h.ConsensusData.ProtoBuf()
	if err != nil {
		return nil, err
	}
	if h.Payload == nil {
		return nil, errors.New(log, "header payload is nil")
	}
	payload, err := h.Payload.Serialize()
	if err != nil {
		return nil, err
	}
	return &pb.HeaderTx{
		Header: &pb.Header{
			Type:          uint32(h.Type),
			Payload:       payload,
			Version:       h.Version,
			ChainID:       h.ChainID.Bytes(),
			Timestamp:     h.TimeStamp,
			Height:        h.Height,
			ConsensusData: pbCon,
			PrevHash:      h.PrevHash.Bytes(),
			MerkleHash:    h.MerkleHash.Bytes(),
			StateHash:     h.StateHash.Bytes(),
			Bloom:         h.Bloom.Bytes(),
		},
		Receipt:   &pb.BlockReceipt{BlockCpu: h.Receipt.BlockCpu, BlockNet: h.Receipt.BlockNet},
		Sign:      sig,
		BlockHash: h.Hash.Bytes(),
	}, nil
}

/**
 *  @brief converts a structure into a sequence of characters
 *  @return []byte - a sequence of characters
 */
func (h *Header) Serialize() ([]byte, error) {
	p, err := h.protoBuf()
	if err != nil {
		return nil, err
	}
	data, err := p.Marshal()
	if err != nil {
		return nil, errors.New(log, fmt.Sprintf("ProtoBuf Marshal error:%s", err.Error()))
	}
	return data, nil
}

/**
 *  @brief converts a sequence of characters into a structure
 *  @param data - a sequence of characters
 */
func (h *Header) Deserialize(data []byte) error {
	if len(data) == 0 {
		return errors.New(log, "input data's length is zero")
	}
	var pbHeader pb.HeaderTx
	if err := pbHeader.Unmarshal(data); err != nil {
		return err
	}

	switch HeaderType(pbHeader.Header.Type) {
	case HeMinorBlock:
		h.Payload = new(MinorBlockHeader)
	case HeCmBlock:
		h.Payload = new(CMBlockHeader)
	default:
		return errors.New(log, "unknown header type")
	}
	h.Type = HeaderType(pbHeader.Header.Type)
	if err := h.Payload.Deserialize(pbHeader.Header.Payload); err != nil {
		return err
	}

	h.Version = pbHeader.Header.Version
	h.ChainID = common.NewHash(pbHeader.Header.ChainID)
	h.TimeStamp = pbHeader.Header.Timestamp
	h.Height = pbHeader.Header.Height
	h.PrevHash = common.NewHash(pbHeader.Header.PrevHash)
	h.MerkleHash = common.NewHash(pbHeader.Header.MerkleHash)
	for i := 0; i < len(pbHeader.Sign); i++ {
		sig := common.Signature{
			PubKey:  common.CopyBytes(pbHeader.Sign[i].PubKey),
			SigData: common.CopyBytes(pbHeader.Sign[i].SigData),
		}
		h.Signatures = append(h.Signatures, sig)
	}
	h.StateHash = common.NewHash(pbHeader.Header.StateHash)
	h.Hash = common.NewHash(pbHeader.BlockHash)
	h.Bloom = bloom.NewBloom(pbHeader.Header.Bloom)
	h.Receipt = BlockReceipt{BlockNet: pbHeader.Receipt.BlockNet, BlockCpu: pbHeader.Receipt.BlockCpu}

	dataCon, err := pbHeader.Header.ConsensusData.Marshal()
	if err != nil {
		return err
	}
	if err := h.ConsensusData.Deserialize(dataCon); err != nil {
		return err
	}

	return nil
}

func (h *Header) JsonString() string {
	data, err := json.Marshal(
		struct {
			ChainID       string
			Version       uint32
			TimeStamp     int64
			Height        uint64
			ConsensusData ConsensusData
			PrevHash      string
			MerkleHash    string
			StateHash     string
			bloom         bloom.Bloom
			Signatures    []common.Signature
			Hash          string
		}{
			ChainID:       h.ChainID.HexString(),
			Version:       h.Version,
			TimeStamp:     h.TimeStamp,
			Height:        h.Height,
			ConsensusData: h.ConsensusData,
			PrevHash:      h.PrevHash.HexString(),
			MerkleHash:    h.MerkleHash.HexString(),
			StateHash:     h.StateHash.HexString(),
			Signatures:    h.Signatures,
			Hash:          h.Hash.HexString(),
		})
	if err != nil {
		log.Error(err)
		return ""
	}
	return string(data)
}

func (h *Header) Show() {
	log.Debug("\t\tshow header:")
	log.Debug(h.JsonString())
}
