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
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ecoball/eballscan/http"
	"github.com/ecoball/eballscan/onlooker"
	"github.com/ecoball/go-ecoball/common/elog"
	"github.com/urfave/cli"
)

var (
	log          = elog.NewLogger("eballscan", elog.DebugLog)
	startCommand = cli.Command{
		Name:   "start",
		Usage:  "start eballscan service",
		Action: startServive,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "ecoball-ip, i",
				Usage: "ecoball full node iP address",
				Value: "127.0.0.1",
			},
			cli.IntFlag{
				Name:  "ecoball-bystander-port, p",
				Usage: "ecoball listen port",
				Value: 9001,
			},
		},
	}
)

func main() {
	app := cli.NewApp()

	//set attribute of EcoBall
	app.Name = "eballscan"
	app.HelpName = "eballscan"
	app.Usage = "Blockchain browser from QuakerChain Technology"
	app.UsageText = "Eballscan is a high concurrency and fast response blockchain browser"
	app.Copyright = "2018 ecoball. All rights reserved"
	app.Author = "EcoBall"
	app.Email = "service@ecoball.org"
	app.HideHelp = true
	app.HideVersion = true

	//commands
	app.Commands = []cli.Command{
		startCommand,
	}

	//run
	app.Run(os.Args)
}

func startServive(c *cli.Context) error {
	ip := c.String("ecoball-ip")
	port := c.Int("ecoball-bystander-port")
	address := fmt.Sprintf(ip+":%d", port)

	fmt.Println(address)
	go onlooker.Bystander(address)
	go http.StartHttpServer()

	wait()
	return nil
}

//capture single
func wait() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer signal.Stop(interrupt)
	sig := <-interrupt
	log.Info("eballscan received signal:", sig)
}
