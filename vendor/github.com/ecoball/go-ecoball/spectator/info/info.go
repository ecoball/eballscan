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

package info

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
)

type NotifyInfo interface {
	Serialize() ([]byte, error)
	Deserialize(data []byte) error
}

type NotifyType int

const (
	InfoNil NotifyType = iota
	InfoBlock
	SynBlock
)

type OneNotify struct {
	InfoType NotifyType
	Info     []byte
}

func NewOneNotify(oneType NotifyType, message NotifyInfo) (*OneNotify, error) {
	oneMessage, err := message.Serialize()
	if nil != err {
		return nil, err
	}
	return &OneNotify{oneType, oneMessage}, nil
}

func (this *OneNotify) Serialize() ([]byte, error) {
	return json.Marshal(*this)
}

func (this *OneNotify) Deserialize(data []byte) error {
	return json.Unmarshal(data, this)
}

const (
	MESSAGE_SIZE int = 10
)

func MessageDecorate(dataBody []byte) []byte {
	data := []byte(fmt.Sprintf("%d", len(dataBody)))
	if length := len(data); length < MESSAGE_SIZE {
		for i := length + 1; i <= MESSAGE_SIZE; i++ {
			data = append(data, '#')
		}
	}
	data = append(data, dataBody...)
	return data
}

func ReadData(conn net.Conn) ([]byte, int, error) {
	bufSize := make([]byte, MESSAGE_SIZE)
	if n, err := conn.Read(bufSize); err != nil || MESSAGE_SIZE != n {
		return nil, 0, errors.New("server conn.Read read message head error")
	}

	var i int
	for i = 0; i < MESSAGE_SIZE; i++ {
		if bufSize[i] == '#' {
			break
		}
	}

	fmt.Println(i, string(bufSize))

	length, err := strconv.Atoi(string(bufSize[:i]))
	if nil != err {
		return nil, 0, err
	}

	buf := make([]byte, length)
	n, err := conn.Read(buf)
	if err != nil || length != n {
		return nil, 0, errors.New("server conn.Read read message body error")
	}

	return buf, n, nil
}
