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

package syn

import (
	"encoding/json"
	"net"

	"github.com/ecoball/go-ecoball/common/elog"
	"github.com/ecoball/go-ecoball/common/message/mpb"
	"github.com/ecoball/go-ecoball/spectator/info"
)

var (
	log = elog.NewLogger("syn", elog.DebugLog)
)

type BlockHeight int

func (this *BlockHeight) Serialize() ([]byte, error) {
	return json.Marshal(*this)
}

func (this *BlockHeight) Deserialize(data []byte) error {
	return json.Unmarshal(data, this)
}

func (this *BlockHeight) Type() uint32 {
	return 0
}
func (this *BlockHeight) Identify() mpb.Identify {
	return mpb.Identify(0)
}

type CommitteeHeight int

func (this *CommitteeHeight) Serialize() ([]byte, error) {
	return json.Marshal(*this)
}

func (this *CommitteeHeight) Deserialize(data []byte) error {
	return json.Unmarshal(data, this)
}

func (this *CommitteeHeight) Type() uint32 {
	return 1
}
func (this *CommitteeHeight) Identify() mpb.Identify {
	return mpb.Identify(1)
}

type FinalHeight int

func (this *FinalHeight) Serialize() ([]byte, error) {
	return json.Marshal(*this)
}

func (this *FinalHeight) Deserialize(data []byte) error {
	return json.Unmarshal(data, this)
}

func (this *FinalHeight) Type() uint32 {
	return 2
}
func (this *FinalHeight) Identify() mpb.Identify {
	return mpb.Identify(2)
}

type MinorHeight int

func (this *MinorHeight) Serialize() ([]byte, error) {
	return json.Marshal(*this)
}

func (this *MinorHeight) Deserialize(data []byte) error {
	return json.Unmarshal(data, this)
}

func (this *MinorHeight) Type() uint32 {
	return 3
}
func (this *MinorHeight) Identify() mpb.Identify {
	return mpb.Identify(3)
}

type ViewChangeHeight int

func (this *ViewChangeHeight) Serialize() ([]byte, error) {
	return json.Marshal(*this)
}

func (this *ViewChangeHeight) Deserialize(data []byte) error {
	return json.Unmarshal(data, this)
}

func (this *ViewChangeHeight) Type() uint32 {
	return 4
}
func (this *ViewChangeHeight) Identify() mpb.Identify {
	return mpb.Identify(4)
}

func SynBlocks(conn net.Conn, message info.NotifyInfo) {
	//height := BlockHeight(database.MaxHeight)
	oneNotify, err := info.NewOneNotify(info.SynBlock, message)
	if nil != err {
		log.Error("SynBlocks newOneNotify error: ", err)
		return
	}

	one, err := oneNotify.Serialize()
	if nil != err {
		log.Error("SynBlocks Serialize error: ", err)
		return
	}

	one = info.MessageDecorate(one)

	if _, err := conn.Write(one); nil != err {
		log.Error("SynBlocks Write error: ", err)
	}
}
