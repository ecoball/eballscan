package database

import (
	"encoding/json"
	"net"

	"github.com/ecoball/go-ecoball/spectator/info"
)

type BlockHight int

func (this *BlockHight) Serialize() ([]byte, error) {
	return json.Marshal(*this)
}

func (this *BlockHight) Deserialize(data []byte) error {
	return json.Unmarshal(data, this)
}

func SynBlocks(conn net.Conn) {
	hight := BlockHight(MaxHight)
	oneNotify, err := info.NewOneNotify(info.SynBlock, &hight)
	if nil != err {
		log.Error("SynBlocks newOneNotify error: ", err)
		return
	}

	info, err := oneNotify.Serialize()
	if nil != err {
		log.Error("SynBlocks Serialize error: ", err)
		return
	}

	if _, err := conn.Write(info); nil != err {
		log.Error("SynBlocks Write error: ", err)
	}
}
