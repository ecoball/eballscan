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

	"github.com/ecoball/eballscan/database"
	"github.com/ecoball/go-ecoball/common/elog"
	"github.com/ecoball/go-ecoball/spectator/info"
)

var (
	log = elog.NewLogger("syn", elog.DebugLog)
)

type BlockHight int

func (this *BlockHight) Serialize() ([]byte, error) {
	return json.Marshal(*this)
}

func (this *BlockHight) Deserialize(data []byte) error {
	return json.Unmarshal(data, this)
}

func SynBlocks(conn net.Conn) {
	hight := BlockHight(database.MaxHight)
	oneNotify, err := info.NewOneNotify(info.SynBlock, &hight)
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
