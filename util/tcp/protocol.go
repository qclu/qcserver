package tcp

import (
	"strconv"
)

type QcMessageHeader struct {
	Magic         uint32
	TransactionId uint64
	Type          uint16
	DataLen       uint32
	SrvId         uint64
	Status        int8
	Reserved      [24]byte
}

func NewMsgHeader() {
	buf := []byte("NJCJ")
	magic := binary.BigEndian.Uint32(buf)
	header := &QcMessageHeader{
		Magic: magic,
	}
	return header
}
