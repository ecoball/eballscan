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

package onlooker

import (
	"net"
	"os"

	"github.com/ecoball/eballscan/notify"
	"github.com/ecoball/eballscan/syn"
	"github.com/ecoball/go-ecoball/common/elog"
	"github.com/ecoball/go-ecoball/spectator/info"
)

var (
	log  = elog.NewLogger("onlooker", elog.DebugLog)
	Conn net.Conn
)

func Bystander() {
	//Connect to server node
	var err error
	Conn, err = net.Dial("tcp", "127.0.0.1:9000")
	if err != nil {
		log.Error("explorer server net.Dial error: ", err)
		os.Exit(1)
	}

	//synchronous data
	go syn.SynBlocks(Conn)

	//Get the notify data and process it
	for {
		buf, n, err := info.ReadData(Conn)
		if nil != err {
			log.Error("explorer server read data error: ", err)
			continue
		}

		one := info.OneNotify{info.InfoNil, []byte{}}
		if err := one.Deserialize(buf[:n]); nil != err {
			log.Error("explorer server notify.Deserialize error: ", err)
			continue
		}
		go notify.Dispatch(one)
	}
}
