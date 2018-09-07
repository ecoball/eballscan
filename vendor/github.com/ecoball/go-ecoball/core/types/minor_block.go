package types

import (
	"github.com/ecoball/go-ecoball/common"
	"github.com/ecoball/go-ecoball/core/pb"
	"github.com/ecoball/go-ecoball/common/errors"
	"fmt"
	"encoding/json"
)

type MinorBlockHeader struct {
	ProposalPublicKey []byte
	StateChangeHash   common.Hash
	ShardId           uint16
	CMEpochNo         uint64
	CmBlockHash       common.Hash
}

func (m *MinorBlockHeader) Serialize() ([]byte, error) {
	protoHeader := pb.MinorBlockHeader{
		ProposalPublicKey: common.CopyBytes(m.ProposalPublicKey),
		StateChangeHash:   m.StateChangeHash.Bytes(),
		ShardId:           uint32(m.ShardId),
		CMEpochNo:         m.CMEpochNo,
		CmBlockHash:       m.CmBlockHash.Bytes(),
	}
	data, err := protoHeader.Marshal()
	if err != nil {
		return nil, errors.New(log, fmt.Sprintf("ProtoBuf Marshal error:%s", err.Error()))
	}
	return data, nil
}

func (m *MinorBlockHeader) Deserialize(data []byte) error {
	var pbHeader pb.MinorBlockHeader
	if err := pbHeader.Unmarshal(data); err != nil {
		return err
	}
	m.ProposalPublicKey = common.CopyBytes(pbHeader.ProposalPublicKey)
	m.StateChangeHash = common.NewHash(pbHeader.StateChangeHash)
	m.ShardId = uint16(pbHeader.ShardId)
	m.CMEpochNo = pbHeader.CMEpochNo
	m.CmBlockHash = common.NewHash(pbHeader.CmBlockHash)

	return nil
}

func (m MinorBlockHeader) GetObject() interface{} {
	return m
}

func (m *MinorBlockHeader) JsonString() string {
	data, err := json.Marshal(struct {
		ProposalPublicKey string
		StateChangeHash   string
		ShardId           uint16
		CMEpochNo         uint64
		CmBlockHash       string
	}{
		ProposalPublicKey: common.ToHex(m.ProposalPublicKey),
		StateChangeHash:   m.StateChangeHash.HexString(),
		ShardId:           m.ShardId,
		CMEpochNo:         m.CMEpochNo,
		CmBlockHash:       m.CmBlockHash.HexString(),
	})
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(data)
}

func (m *MinorBlockHeader) Show() {
	log.Debug(m.JsonString())
}

func (m *MinorBlockHeader) Type() uint32 {
	return uint32(HeMinorBlock)
}