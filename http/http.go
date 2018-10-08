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
	router.POST("/eballscan/getBlock", getBlock)
	router.POST("/eballscan/getBlockByHeight", getBlockByHeight)
	router.POST("/eballscan/add_block", addBlock)
	router.POST("/eballscan/getTransactionByHash", getTransactionByHash)
	router.POST("/eballscan/add_transaction", addTransaction)
	router.POST("/eballscan/getTransactionByHight", getTransactionByHight)
	router.POST("/eballscan/getTransaction", getTransaction)

	http.ListenAndServe(":20680", router)
	return nil
}

func getBlockByHeight(c *gin.Context) {
	height_str := c.PostForm("hight")
	height, err := strconv.Atoi(height_str)
	if nil != err{
		panic(err) 
	}
	info, err := database.QueryOneBlock(height)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"CountTxs": info.CountTxs, "StateHash": info.StateHash, "hash": info.Hash, "MerkleHash": info.MerkleHash, "PrevHash": info.PrevHash,
			"timeStamp": info.Timestamp})
}

func getBlock(c *gin.Context) {
	num_str := c.PostForm("num")
	num, err := strconv.Atoi(num_str)
	if nil != err{
		panic(err) 
	}

	index_str := c.PostForm("index")
	index, err := strconv.Atoi(index_str)
	if nil != err{
		panic(err) 
	}

	info, pageNum, err := database.QueryBlock(index, num)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"pageNum": pageNum, "blocks": info})
}

func getTransaction(c *gin.Context) {
	num_str := c.PostForm("num")
	num, err := strconv.Atoi(num_str)
	if nil != err{
		panic(err) 
	}

	index_str := c.PostForm("index")
	index, err := strconv.Atoi(index_str)
	if nil != err{
		panic(err) 
	}

	info, pageNum, err := database.QueryTransaction(index, num)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"pageNum": pageNum, "transactions": info})
}

func addBlock(c *gin.Context) {
	hight_str := c.PostForm("hight")
	hight, err := strconv.Atoi(hight_str)
	if nil != err{
		panic(err) 
	}

	time_str := c.PostForm("timeStamp")
	timeStamp, err := strconv.Atoi(time_str)
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
	errcode := database.AddBlock(hight, countTxs, timeStamp, hash, prevHash, merkleHash, stateHash)
	if nil != errcode{
		c.JSON(http.StatusBadRequest, gin.H{"result": errcode.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

func getTransactionByHash(c *gin.Context) {
	hash := c.PostForm("hash")
	info, err := database.QueryOneTransaction(hash)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"transaction": info})
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

func getTransactionByHight(c *gin.Context) {
	height_str := c.PostForm("blockHight")
	blockHight, err := strconv.Atoi(height_str)
	if nil != err{
		panic(err) 
	}
	
	info, err := database.QueryTransactionsByHight(blockHight)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
	}

	/*datas := []data.TransactionInfoH{}
	for _, v := range info {
		datas = append(datas, *v)
	}*/

	c.JSON(http.StatusOK, gin.H{"counts": len(info), "transactions": info})
}