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

package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ecoball/eballscan/onlooker"
	"github.com/ecoball/eballscan/http"
	"github.com/ecoball/go-ecoball/common/elog"
)

var log = elog.NewLogger("wallet", elog.DebugLog)

func main() {
	go onlooker.Bystander()
	go http.StartHttpServer()

	wait()
}

//capture single
func wait() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer signal.Stop(interrupt)
	sig := <-interrupt
	log.Info("ecoscan received signal:", sig)
}

