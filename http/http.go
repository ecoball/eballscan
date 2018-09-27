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

package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ecoball/eballscan/database"
	"strconv"
)

func StartHttpServer() (err error) {
	//get router instance
	router := gin.Default()

	//register handle
	router.POST("/eballscan/get_block", getBlock)
	router.POST("/eballscan/add_block", addBlock)
	router.POST("/eballscan/get_transaction", getTransaction)
	router.POST("/eballscan/add_transaction", addTransaction)

	http.ListenAndServe(":20680", router)
	return nil
}

func getBlock(c *gin.Context) {
	height_str := c.PostForm("height")
	height, err := strconv.Atoi(height_str)
	if nil != err{
		panic(err) 
	}
	info, err := database.QueryOneBlock(height)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"CountTxs": info.CountTxs, "StateHash": info.StateHash, "hash": info.Hash, "MerkleHash": info.MerkleHash, "PrevHash": info.PrevHash})
}

func addBlock(c *gin.Context) {
	height_str := c.PostForm("height")
	height, err := strconv.Atoi(height_str)
	if nil != err{
		panic(err) 
	}
	
	countTxs_str := c.PostForm("countTxs")
	countTxs, err := strconv.Atoi(countTxs_str)
	if nil != err{
		panic(err) 
	}

	hash := c.PostForm("hash")
	prevHash := c.PostForm("prevHash")
	merkleHash := c.PostForm("merkleHash")
	stateHash := c.PostForm("stateHash")
	errcode := database.AddBlock(height, countTxs, hash, prevHash, merkleHash, stateHash)
	if nil != errcode{
		c.JSON(http.StatusBadRequest, gin.H{"result": errcode.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

func getTransaction(c *gin.Context) {
	hash := c.PostForm("hash")
	info, err := database.QueryOneTransaction(hash)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"Address": info.Address, "BlockHight": info.BlockHight, "Permission": info.Permission, "TimeStamp": info.TimeStamp, "TxFrom":info.TxFrom, 
	"TxType": info.TxType})
}

func addTransaction(c *gin.Context) {
	txType_str := c.PostForm("txType")
	txType, err := strconv.Atoi(txType_str)
	if nil != err{
		panic(err) 
	}
	
	timeStamp_str := c.PostForm("timeStamp")
	timeStamp, err := strconv.Atoi(timeStamp_str)
	if nil != err{
		panic(err) 
	}

	blockHight_str := c.PostForm("blockHight")
	blockHight, err := strconv.Atoi(blockHight_str)
	if nil != err{
		panic(err) 
	}

	hash := c.PostForm("hash")
	permission := c.PostForm("permission")
	txFrom := c.PostForm("txFrom")
	address := c.PostForm("address")
	errcode := database.AddTransaction(txType, timeStamp, blockHight, hash, permission, txFrom, address)
	if nil != errcode{
		c.JSON(http.StatusBadRequest, gin.H{"result": errcode.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}