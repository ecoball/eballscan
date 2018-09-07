package types

import (
	"encoding/json"
	"fmt"
	"github.com/ecoball/go-ecoball/common"
	"github.com/ecoball/go-ecoball/common/errors"
	"github.com/ecoball/go-ecoball/core/pb"
)

type NodeInfo struct {
	PublicKey []byte
}

type NodeAddr struct {
	Address string
	Port    string
}

type Shard struct {
	Id         uint32
	Member     []NodeInfo
	MemberAddr []NodeAddr
}

type CMBlockHeader struct {
	LeaderPubKey    []byte
	CandidatePubKey []byte
	Nonce           uint32
	ShardsHash      common.Hash /*shards hash, not include node address*/
}

func (c *CMBlockHeader) Serialize() ([]byte, error) {
	pbHeader := pb.CMBlockHeader{
		LeaderPubKey:    common.CopyBytes(c.LeaderPubKey),
		CandidatePubKey: common.CopyBytes(c.CandidatePubKey),
		Nonce:           c.Nonce,
		ShardsHash:      c.ShardsHash.Bytes(),
	}
	data, err := pbHeader.Marshal()
	if err != nil {
		return nil, errors.New(log, fmt.Sprintf("ProtoBuf Marshal error:%s", err.Error()))
	}
	return data, nil
}

func (c *CMBlockHeader) Deserialize(data []byte) error {
	var pbHeader pb.CMBlockHeader
	if err := pbHeader.Unmarshal(data); err != nil {
		return err
	}
	c.LeaderPubKey = common.CopyBytes(pbHeader.LeaderPubKey)
	c.CandidatePubKey = common.CopyBytes(pbHeader.CandidatePubKey)
	c.Nonce = pbHeader.Nonce
	c.ShardsHash = common.NewHash(pbHeader.ShardsHash)
	return nil
}

func (c CMBlockHeader) GetObject() interface{} {
	return c
}

func (c *CMBlockHeader) JsonString() string {
	data, err := json.Marshal(struct {
		LeaderPubKey    string
		CandidatePubKey string
		Nonce           uint32
		ShardsHash      string
	}{
		LeaderPubKey:    common.ToHex(c.LeaderPubKey),
		CandidatePubKey: common.ToHex(c.CandidatePubKey),
		Nonce:           c.Nonce,
		ShardsHash:      c.ShardsHash.HexString(),
	})
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(data)
}

func (c *CMBlockHeader) Show() {
	log.Debug(c.JsonString())
}

func (c *CMBlockHeader) Type() uint32 {
	return uint32(HeCmBlock)
}
