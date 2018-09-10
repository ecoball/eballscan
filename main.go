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
	"github.com/ecoball/eballscan/data"
	"github.com/ecoball/eballscan/onlooker"
	"github.com/kataras/iris"
)

func main() {
	go onlooker.Bystander()
	app := iris.New()

	app.OnErrorCode(iris.StatusInternalServerError, func(ctx iris.Context) {

		errMessage := ctx.Values().GetString("error")
		if errMessage != "" {
			ctx.Writef("Internal server error: %s", errMessage)
			return
		}

		ctx.Writef("(Unexpected) internal server error")
	})

	app.Use(func(ctx iris.Context) {
		ctx.Application().Logger().Infof("Begin request for path: %s", ctx.Path())
		ctx.Next()
	})
	app.Get("/b", func(ctx iris.Context) {
	
		ctx.HTML(data.PrintBlock())
	

	})
	app.Get("/t", func(ctx iris.Context) {

		ctx.HTML(data.PrintTransaction())

	})

	app.Run(iris.Addr(":8080"), iris.WithCharset("UTF-8"), iris.WithoutVersionChecker)
}
