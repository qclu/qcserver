package tcpnetwork

import (
	"encoding/binary"
)

//message type
const (
	MSG_BYE         = 201
	MSG_IDEL        = MSG_BYE + 1
	MSG_DON         = MSG_IDEL + 1
	MSG_DOFF        = MSG_DON + 1
	MSG_LOG         = MSG_DOFF + 1
	MSG_FAULT       = MSG_LOG + 1
	MSG_RON         = MSG_FAULT + 1
	MSG_ROFF        = MSG_RON + 1
	MSG_STAT        = MSG_ROFF + 1
	MSG_TEST_DOA    = MSG_STAT + 1
	MSG_QUPDATE_DOA = MSG_TEST_DOA + 1
	MSG_AUPDATE_DOA = MSG_QUPDATE_DOA + 1
)

type QcMessage struct {
	MsgHeader *QcMessageHeader
	MsgTail   *QcMessageTail
	MsgData   []byte
}

func NewMessage() *QcMessage {
	return &QcMessage{}
}

//64Bytes
type QcMessageHeader struct {
	Magic         uint32 //4
	TransactionId uint64 //8
	Type          uint16 //2
	DataLen       uint32 //4
	SrvId         uint64 //8
	Status        int8   //1
	Reserved      [37]byte
}

func NewMsgHeader() *QcMessageHeader {
	buf := []byte("NJCJ")
	magic := binary.BigEndian.Uint32(buf)
	header := &QcMessageHeader{
		Magic: magic,
		SrvId: 0,
	}
	return header
}

type QcMessageTail struct {
	Crc      [8]byte
	Reserved [24]byte
}

func NewMsgTail() *QcMessageTail {
	tail := &QcMessageTail{}
	return tail
}
