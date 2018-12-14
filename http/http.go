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
	"bytes"
	"log"
	"net/http"
	"os/exec"

	"strconv"

	"github.com/ecoball/eballscan/database"
	"github.com/gin-gonic/gin"
)

func StartHttpServer() (err error) {
	//get router instance
	router := gin.Default()

	//test
	//router.POST("/eballscan/add_transaction", addTransaction)
	/*router.POST("/eballscan/addCommittee_block", addCommittee_block)
	router.POST("/eballscan/addFinal_block", addFinal_block)
	router.POST("/eballscan/addMinor_block", addMinor_block)
	router.POST("/eballscan/addNode", addNode)
	router.POST("/eballscan/addViewchangeblock", addViewchangeblock)*/

	//register handle
	router.POST("/eballscan/getBlock", getBlock)
	router.POST("/eballscan/getBlockByHeight", getBlockByHeight)

	//transaction
	router.POST("/eballscan/getTransactionByHash", getTransactionByHash)
	router.POST("/eballscan/getTransactionByHeightAndShardId", getTransactionByHeightAndShardId)
	router.POST("/eballscan/getTransaction", getTransaction)
	router.POST("/eballscan/getTransactionsByAccountName", getTransactionsByAccountName)
	router.POST("/eballscan/getAccounts", getAccounts)
	router.POST("/eballscan/getAccountInfo", getAccountInfo)

	//committee block
	router.POST("/eballscan/getCommitteeBlock", getCommitteeBlock)
	router.POST("/eballscan/getCommitteeBlockByHeight", getCommitteeBlockByHeight)
	router.POST("/eballscan/getCommitteeBlockByHash", getCommitteeBlockByHash)
	router.POST("/eballscan/getNodesByHeight", getNodesByHeight)
	router.POST("/eballscan/getNodeByPubKey", getNodeByPubKey)
	router.POST("/eballscan/getNodes", getNodes)

	//final block
	router.POST("/eballscan/getFinalBlock", getFinalBlock)
	router.POST("/eballscan/getFinalBlockByHeight", getFinalBlockByHeight)

	//minor block
	router.GET("/eballscan/getMaxMinorBlockShardId", getMaxMinorBlockShardId)
	router.POST("/eballscan/getMinorBlockByShardId", getMinorBlockByShardId)
	router.POST("/eballscan/getMinorBlockByHeight", getMinorBlockByHeight)
	router.POST("/eballscan/getMinorBlockByHeightAndShardId", getMinorBlockByHeightAndShardId)

	//view change block
	router.POST("/eballscan/getViewChangeBlock", getViewChangeBlock)
	router.POST("/eballscan/getViewChangeBlockByHeight", getViewChangeBlockByHeight)
	router.POST("/eballscan/getViewChangeBlockByFinalBlockHeight", getViewChangeBlockByFinalBlockHeight)

	http.ListenAndServe(":20680", router)
	return nil
}

func getBlockByHeight(c *gin.Context) {
	height_str := c.PostForm("height")
	height, err := strconv.Atoi(height_str)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	info, max_height, err := database.QueryOneBlock(height)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"max_height": max_height, "block": info})
}

func getBlock(c *gin.Context) {
	num_str := c.PostForm("num")
	num, err := strconv.Atoi(num_str)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	index_str := c.PostForm("index")
	index, err := strconv.Atoi(index_str)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	info, pageNum, err := database.QueryBlock(index, num)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pageNum": pageNum, "blocks": info})
}

func getTransaction(c *gin.Context) {
	num_str := c.PostForm("num")
	num, err := strconv.Atoi(num_str)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	index_str := c.PostForm("index")
	index, err := strconv.Atoi(index_str)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	info, pageNum, err := database.QueryTransaction(index, num)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pageNum": pageNum, "transactions": info})
}

/*func addBlock(c *gin.Context) {
	height_str := c.PostForm("Height")
	height, err := strconv.Atoi(height_str)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	time_str := c.PostForm("TimeStamp")
	timeStamp, err := strconv.Atoi(time_str)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	countTxs_str := c.PostForm("CountTxs")
	countTxs, err := strconv.Atoi(countTxs_str)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	hash := c.PostForm("Hash")
	prevHash := c.PostForm("PrevHash")
	merkleHash := c.PostForm("MerkleHash")
	stateHash := c.PostForm("StateHash")
	errcode := database.AddBlock(height, countTxs, timeStamp, hash, prevHash, merkleHash, stateHash)
	if nil != errcode{
		c.JSON(http.StatusBadRequest, gin.H{"result": errcode.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}*/

func getTransactionByHash(c *gin.Context) {
	hash := c.PostForm("hash")
	info, err := database.QueryOneTransaction(hash)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transaction": info})
}

/*func addTransaction(c *gin.Context) {
	txType_str := c.PostForm("TxType")
	txType, err := strconv.Atoi(txType_str)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	timeStamp_str := c.PostForm("TimeStamp")
	timeStamp, err := strconv.Atoi(timeStamp_str)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	blockHeight_str := c.PostForm("BlockHeight")
	blockHeight, err := strconv.Atoi(blockHeight_str)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	hash := c.PostForm("Hash")
	permission := c.PostForm("Permission")
	txFrom := c.PostForm("TxFrom")
	address := c.PostForm("Address")
	errcode := database.AddTransaction(txType, timeStamp, blockHeight, hash, permission, txFrom, address)
	if nil != errcode{
		c.JSON(http.StatusBadRequest, gin.H{"result": errcode.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}*/

func getTransactionByHeightAndShardId(c *gin.Context) {
	height_str := c.PostForm("height")
	blockHeight, err := strconv.Atoi(height_str)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	shardId_str := c.PostForm("shardId")
	shardId, err := strconv.Atoi(shardId_str)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	info, err := database.QueryTransactionsByHeightAndShardId(blockHeight, shardId)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"counts": len(info), "transactions": info})
}

func getTransactionsByAccountName(c *gin.Context) {
	name := c.PostForm("name")

	num_str := c.PostForm("num")
	num, err := strconv.Atoi(num_str)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	index_str := c.PostForm("index")
	index, err := strconv.Atoi(index_str)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	info, pageNum, counts, err := database.QueryTransactionsByAccountName(num, index, name)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pageNums": pageNum, "counts": counts, "transactions": info})
}

func getAccounts(c *gin.Context) {
	num_str := c.PostForm("num")
	num, err := strconv.Atoi(num_str)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	index_str := c.PostForm("index")
	index, err := strconv.Atoi(index_str)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	info, pageNum, err := database.QueryAccounts(num, index)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pageNum": pageNum, "accounts": info})
}

func getAccountInfo(c *gin.Context) {
	name := c.PostForm("name")

	data, err := database.QueryOneAccount(name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"account": data})
}

func exec_shell(s string) {
	cmd := exec.Command("/bin/bash", "-c", s)
	var out bytes.Buffer

	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func Compile() {
	exec_shell("cd ../c2wasm-compiler/;./compile.sh api.c api;cd ../eballscan/")
}

func addCommittee_block(c *gin.Context) {
	Nonce_str := c.PostForm("Nonce")
	Nonce, err := strconv.Atoi(Nonce_str)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	timeStamp_str := c.PostForm("TimeStamp")
	timeStamp, err := strconv.Atoi(timeStamp_str)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	Height_str := c.PostForm("Height")
	Height, err := strconv.Atoi(Height_str)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	NodeCounts_str := c.PostForm("NodeCounts")
	NodeCounts, err := strconv.Atoi(NodeCounts_str)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	hash := c.PostForm("Hash")
	PrevHash := c.PostForm("PrevHash")
	ShardsHash := c.PostForm("ShardsHash")
	LeaderPubKey := c.PostForm("LeaderPubKey")
	CandidatePublicKey := c.PostForm("CandidatePublicKey")
	CandidateAddress := c.PostForm("CandidateAddress")
	CandidatePort := c.PostForm("CandidatePort")

	errcode := database.AddCommittee_block(Height, Nonce, timeStamp, NodeCounts, hash, PrevHash, ShardsHash, LeaderPubKey, CandidatePort, CandidateAddress, CandidatePublicKey)
	if nil != errcode {
		c.JSON(http.StatusBadRequest, gin.H{"result": errcode.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

func addFinal_block(c *gin.Context) {
	TrxCount_str := c.PostForm("TrxCount")
	TrxCount, err := strconv.Atoi(TrxCount_str)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	timeStamp_str := c.PostForm("TimeStamp")
	timeStamp, err := strconv.Atoi(timeStamp_str)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	Height_str := c.PostForm("Height")
	Height, err := strconv.Atoi(Height_str)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	EpochNo_str := c.PostForm("EpochNo")
	EpochNo, err := strconv.Atoi(EpochNo_str)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	hash := c.PostForm("Hash")
	PrevHash := c.PostForm("PrevHash")
	CMBlockHash := c.PostForm("CMBlockHash")
	TrxRootHash := c.PostForm("TrxRootHash")
	StateDeltaRootHash := c.PostForm("StateDeltaRootHash")
	MinorBlocksHash := c.PostForm("MinorBlocksHash")
	StateHashRoot := c.PostForm("StateHashRoot")
	ProposalPubKey := c.PostForm("ProposalPubKey")

	errcode := database.AddFinal_block(Height, timeStamp, -1, TrxCount, EpochNo, hash, PrevHash, CMBlockHash, TrxRootHash, StateDeltaRootHash, MinorBlocksHash,
		StateHashRoot, ProposalPubKey)
	if nil != errcode {
		c.JSON(http.StatusBadRequest, gin.H{"result": errcode.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

func addMinor_block(c *gin.Context) {
	ShardId_str := c.PostForm("ShardId")
	ShardId, err := strconv.Atoi(ShardId_str)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	timeStamp_str := c.PostForm("TimeStamp")
	timeStamp, err := strconv.Atoi(timeStamp_str)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	Height_str := c.PostForm("Height")
	Height, err := strconv.Atoi(Height_str)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	CMEpochNo_str := c.PostForm("CMEpochNo")
	CMEpochNo, err := strconv.Atoi(CMEpochNo_str)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	hash := c.PostForm("Hash")
	PrevHash := c.PostForm("PrevHash")
	TrxHashRoot := c.PostForm("TrxHashRoot")
	StateDeltaHash := c.PostForm("StateDeltaHash")
	CMBlockHash := c.PostForm("CMBlockHash")
	ProposalPublicKey := c.PostForm("ProposalPublicKey")
	counts := 1

	errcode := database.AddMinor_block(Height, timeStamp, ShardId, -1, CMEpochNo, counts, hash, PrevHash, TrxHashRoot, StateDeltaHash, CMBlockHash, ProposalPublicKey)

	if nil != errcode {
		c.JSON(http.StatusBadRequest, gin.H{"result": errcode.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

func addNode(c *gin.Context) {
	Committee_blockHeight_str := c.PostForm("Committee_blockHeight")
	Committee_blockHeight, err := strconv.Atoi(Committee_blockHeight_str)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	PublicKey := c.PostForm("PublicKey")
	Address := c.PostForm("Address")
	Port := c.PostForm("Port")

	errcode := database.AddNode(PublicKey, Port, Address, Committee_blockHeight)
	if nil != errcode {
		c.JSON(http.StatusBadRequest, gin.H{"result": errcode.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

func addViewchangeblock(c *gin.Context) {
	Round_str := c.PostForm("Round")
	Round, err := strconv.Atoi(Round_str)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	timeStamp_str := c.PostForm("TimeStamp")
	timeStamp, err := strconv.Atoi(timeStamp_str)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	Height_str := c.PostForm("Height")
	Height, err := strconv.Atoi(Height_str)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	CMEpochNo_str := c.PostForm("CMEpochNo")
	CMEpochNo, err := strconv.Atoi(CMEpochNo_str)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	FinalBlockHeight_str := c.PostForm("FinalBlockHeight")
	FinalBlockHeight, err := strconv.Atoi(FinalBlockHeight_str)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	hash := c.PostForm("Hash")
	PrevHash := c.PostForm("PrevHash")
	CandidatePublicKey := c.PostForm("CandidatePublicKey")
	CandidateAddress := c.PostForm("CandidateAddress")
	CandidatePort := c.PostForm("CandidatePort")

	errcode := database.AddViewchangeblock(Height, timeStamp, Round, CMEpochNo, FinalBlockHeight, hash, PrevHash, CandidatePort, CandidateAddress, CandidatePublicKey)

	if nil != errcode {
		c.JSON(http.StatusBadRequest, gin.H{"result": errcode.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}
