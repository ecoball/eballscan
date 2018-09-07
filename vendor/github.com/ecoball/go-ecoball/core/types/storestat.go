// Copyright 2018 The go-ecoball Authors
// This file is part of the go-ecoball.
//
// The go-ecoball is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ecoball is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ecoball. If not, see <http://www.gnu.org/licenses/>.

package types

import (
	"errors"
	"github.com/ecoball/go-ecoball/core/pb"
	"gx/ipfs/QmYVNvtQkeZ6AKSwDrjQTs432QtL6umrrK41EBq3cu7iSP/go-cid"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
)

type ShardInfo struct {
	Cid     *cid.Cid
	RawSize uint64
}

type StoreRepoStat struct {
	Peer         peer.ID
	ChainID      uint32
	RepoSize     uint64
	StorageMax   uint64
	NumObjects   uint64
	RandBlkInfo  []*ShardInfo
	BadBlocks    []*cid.Cid
}

func (srs *StoreRepoStat) Serialize() ([]byte, error) {
	p := &pb.StoreRepoStatMsg{
		PeerHash:   []byte(srs.Peer),
		ChainId:    srs.ChainID,
		RepoSize:   srs.RepoSize,
		StorageMax: srs.StorageMax,
		NumObjects: srs.NumObjects,
	}

	var pb_badcids []*pb.Cid
	for _, cid := range srs.BadBlocks {
		prefix := cid.Prefix()
		pb_cid := &pb.Cid{
			Version:  prefix.Version,
			Codec:    prefix.Codec,
			Hash:     cid.Hash(),
		}
		pb_badcids = append(pb_badcids, pb_cid)
	}
	p.BadBlocks = pb_badcids

	var sis []*pb.ShardInfo
	for _, shard := range srs.RandBlkInfo {
		prefix :=shard.Cid.Prefix()
		pb_cid := &pb.Cid{
			Version:  prefix.Version,
			Codec:    prefix.Codec,
			Hash:     shard.Cid.Hash(),
		}
		si := &pb.ShardInfo{pb_cid, shard.RawSize}
		sis = append(sis, si)
	}
	p.RandBlkInfo = sis

	b, err := p.Marshal()
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (srs *StoreRepoStat) Deserialize(data []byte) error {
	if len(data) == 0 {
		return errors.New("input data's length is zero")
	}
	var pb_srs pb.StoreRepoStatMsg
	if err := pb_srs.Unmarshal(data); err != nil {
		return err
	}
	srs.Peer = peer.ID(pb_srs.PeerHash)
	srs.ChainID = pb_srs.ChainId
	srs.RepoSize = pb_srs.RepoSize
	srs.StorageMax = pb_srs.StorageMax
	srs.NumObjects = pb_srs.NumObjects

	var badcids []*cid.Cid
	for _, pb_cid := range pb_srs.BadBlocks {
		var newCid *cid.Cid
		switch pb_cid.Version {
		case 0:
			newCid = cid.NewCidV0(pb_cid.Hash)
		case 1:
			newCid = cid.NewCidV1(pb_cid.Codec, pb_cid.Hash)
		default:
			return errors.New("error for decoding proof message")
		}
		badcids = append(badcids, newCid)
	}
	srs.BadBlocks = badcids

	var randBlksInfo []*ShardInfo
	for _, shard := range pb_srs.RandBlkInfo {
		var newCid *cid.Cid
		switch shard.ShardCid.Version {
		case 0:
			newCid = cid.NewCidV0(shard.ShardCid.Hash)
		case 1:
			newCid = cid.NewCidV1(shard.ShardCid.Codec, shard.ShardCid.Hash)
		default:
			return errors.New("error for decoding proof message")
		}
		randBlk := &ShardInfo{newCid, shard.ShardSize}
		randBlksInfo = append(randBlksInfo, randBlk)
	}
	srs.RandBlkInfo = randBlksInfo

	return nil
}