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
	"github.com/ecoball/eballscan/database"
)

var (
	log  = elog.NewLogger("onlooker", elog.DebugLog)
	Conn net.Conn
)

func Bystander(address string) {
	//Connect to server node
	var err error
	Conn, err = net.Dial("tcp", address)
	if err != nil {
		log.Error("explorer server net.Dial "+address+" error: ", err)
		os.Exit(1)
	}

	//synchronous data
	go syn_data(Conn)

	//Get the notify data and process it
	for {
		buf, n, err := info.ReadData(Conn)
		if nil != err {
			log.Error("explorer server read data error: ", err)
			continue
		}

		one := info.OneNotify{info.InfoNil, []byte{}, 0}
		if err := one.Deserialize(buf[:n]); nil != err {
			log.Error("explorer server notify.Deserialize error: ", err)
			continue
		}
		notify.Dispatch(one)
	}
}

func syn_data(Conn net.Conn){
	heigt := syn.BlockHeight(database.MaxHeight)
	syn.SynBlocks(Conn, &heigt)

	committeeHeight := syn.CommitteeHeight(database.Max_Committee_Height)
	syn.SynBlocks(Conn, &committeeHeight)

	finalHeight := syn.FinalHeight(database.Max_Final_Height)
	syn.SynBlocks(Conn, &finalHeight)

	/*minorHeight := syn.MinorHeight(database.Max_Minor_Height)
	syn.SynBlocks(Conn, &minorHeight)*/

	viewChangeHeight := syn.ViewChangeHeight(database.Max_ViewChange_Height)
	syn.SynBlocks(Conn, &viewChangeHeight)
}
