package types

import "github.com/ecoball/go-ecoball/common"

const VirtualBlockCpuLimit float64 = 200000000.0
const VirtualBlockNetLimit float64 = 1048576000.0
const BlockCpuLimit float64 = 200000.0
const BlockNetLimit float64 = 1048576.0

type TransactionReceipt struct {
	Hash   common.Hash
	Cpu    float64
	Net    float64
	Result []byte
}

type BlockReceipt struct {
	BlockCpu float64
	BlockNet float64
}