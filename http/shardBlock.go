package http

import (
	"net/http"
	//"bytes"
	//"log"
	//"os/exec" 

	"github.com/gin-gonic/gin"
	"github.com/ecoball/eballscan/database"
	"strconv"
)

func getCommitteeBlock(c *gin.Context){
	num_str := c.PostForm("num")
	num, err := strconv.Atoi(num_str)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return 
	}

	index_str := c.PostForm("index")
	index, err := strconv.Atoi(index_str)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	info, pageNum, err := database.QueryCommitteeBlock(index, num)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pageNum": pageNum, "blocks": info})
}

func getCommitteeBlockByHeight(c *gin.Context) {
	height_str := c.PostForm("height")
	height, err := strconv.Atoi(height_str)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	info, max_height, err := database.QueryOneCommitteeBlock(height)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"max_height": max_height, "block": info})
}

func getNodes(c *gin.Context){
	num_str := c.PostForm("num")
	num, err := strconv.Atoi(num_str)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return 
	}

	index_str := c.PostForm("index")
	index, err := strconv.Atoi(index_str)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	info, pageNum, err := database.QueryNodes(index, num)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pageNum": pageNum, "nodes": info})
}

func getNodesByHeight(c *gin.Context) {
	height_str := c.PostForm("height")
	height, err := strconv.Atoi(height_str)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	info, err := database.QueryNodesByHeight(height)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "success", "counts": len(info), "nodes": info})
}

func getNodeByPubKey(c *gin.Context) {
	PublicKey := c.PostForm("PublicKey")
	info, err := database.QueryOneNode(PublicKey)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"node": info})
}

func getFinalBlock(c *gin.Context){
	num_str := c.PostForm("num")
	num, err := strconv.Atoi(num_str)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return 
	}

	index_str := c.PostForm("index")
	index, err := strconv.Atoi(index_str)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	info, pageNum, err := database.QueryFinalBlock(index, num)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pageNum": pageNum, "blocks": info})
}

func getFinalBlockByHeight(c *gin.Context) {
	height_str := c.PostForm("height")
	height, err := strconv.Atoi(height_str)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	info, max_height, err := database.QueryOneFinalBlock(height)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"max_height": max_height, "block": info})
}

func getMinorBlock(c *gin.Context){
	num_str := c.PostForm("num")
	num, err := strconv.Atoi(num_str)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return 
	}

	index_str := c.PostForm("index")
	index, err := strconv.Atoi(index_str)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	info, pageNum, err := database.QueryMinorBlock(index, num)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pageNum": pageNum, "blocks": info})
}

func getMinorBlockByHeight(c *gin.Context) {
	height_str := c.PostForm("height")
	height, err := strconv.Atoi(height_str)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	info, max_height, err := database.QueryOneMinorBlock(height)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"max_height": max_height, "block": info})
}

func getViewChangeBlock(c *gin.Context){
	num_str := c.PostForm("num")
	num, err := strconv.Atoi(num_str)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return 
	}

	index_str := c.PostForm("index")
	index, err := strconv.Atoi(index_str)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	info, pageNum, err := database.QueryViewChangeBlock(index, num)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pageNum": pageNum, "blocks": info})
}

func getViewChangeBlockByHeight(c *gin.Context) {
	height_str := c.PostForm("height")
	height, err := strconv.Atoi(height_str)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	info, max_height, err := database.QueryOneViewChangeBlock(height)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"max_height": max_height, "block": info})
}

func getViewChangeBlockByFinalBlockHeight(c *gin.Context) {
	height_str := c.PostForm("FinalBlockHeigh")
	height, err := strconv.Atoi(height_str)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	info, err := database.QueryViewChangeBlockByFinalBlockHeight(height)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "success", "blocks": info})
}
