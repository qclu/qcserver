package tcpnetwork

import (
	"bytes"
	"encoding/binary"
	"qcserver/util/log"
)

const (
	ProtocolMsgHeaderLength  = 64
	ProtocolMsgTailLength    = 32
	ProtocolMsgMaxDataLength = 4096
)

// Binary format : | 64 byte (total stream length) | data ... (total stream length - 2) | 32Bytes tail
//	stream protocol interface for 2 bytes header
type StreamProtocol struct {
	serializeBuf   *bytes.Buffer
	unserializeBuf *bytes.Buffer
	logger         *log.Log
}

func NewStreamProtocol() *StreamProtocol {
	protocol := &StreamProtocol{}
	protocol.logger = log.GetLog()
	return protocol
}

func (c *Connection) StreamProtocolInit() {
	c.streamProtocol.serializeBuf = new(bytes.Buffer)
	c.streamProtocol.unserializeBuf = new(bytes.Buffer)
}

func (c *Connection) GetHeaderLength() int {
	return ProtocolMsgHeaderLength
}

//func (c *Connection) UnserializeHeader(buf []byte) {
func (c *Connection) UnserializeHeader(buf []byte) *QcMessageHeader {
	msgheader := NewMsgHeader()
	s := c.streamProtocol
	s.unserializeBuf.Reset()
	s.unserializeBuf.Write(buf)
	binary.Read(s.unserializeBuf, binary.BigEndian, &msgheader.Magic)
	binary.Read(s.unserializeBuf, binary.BigEndian, &msgheader.TransactionId)
	binary.Read(s.unserializeBuf, binary.BigEndian, &msgheader.Type)
	binary.Read(s.unserializeBuf, binary.BigEndian, &msgheader.DataLen)
	binary.Read(s.unserializeBuf, binary.BigEndian, &msgheader.SrvId)
	binary.Read(s.unserializeBuf, binary.BigEndian, &msgheader.Status)
	//binary.Read(s.unserializeBuf, binary.BigEndian, &msgheader.Reserved)
	return msgheader
}

func (c *Connection) SerializeHeader(task *sendTask) []byte {
	msgheader := task.msg.MsgHeader
	s := c.streamProtocol
	if msgheader.DataLen > ProtocolMsgMaxDataLength {
		s.logger.LogError("stream is too long")
		return nil
	}

	//startpos := 0
	s.serializeBuf.Reset()
	binary.Write(s.serializeBuf, binary.BigEndian, &msgheader.Magic)
	//startpos += len(msgheader.Magic)
	binary.Write(s.serializeBuf, binary.BigEndian, &msgheader.TransactionId)
	//startpos += len(msgheader.TransactionId)
	binary.Write(s.serializeBuf, binary.BigEndian, &msgheader.Type)
	//startpos += len(msgheader.Type)
	binary.Write(s.serializeBuf, binary.BigEndian, &msgheader.DataLen)
	//startpos += len(msgheader.DataLen)
	binary.Write(s.serializeBuf, binary.BigEndian, &msgheader.SrvId)
	//startpos += len(msgheader.SrvId)
	binary.Write(s.serializeBuf, binary.BigEndian, &msgheader.Status)
	//startpos += len(msgheader.Status)
	binary.Write(s.serializeBuf, binary.BigEndian, &msgheader.Reserved)
	return s.serializeBuf.Bytes()
}

func (c *Connection) UnserializeTail(buf []byte) *QcMessageTail {
	var msgtail QcMessageTail
	s := c.streamProtocol
	s.unserializeBuf.Reset()
	s.unserializeBuf.Write(buf)
	binary.Read(s.unserializeBuf, binary.BigEndian, &msgtail.Crc)
	binary.Read(s.unserializeBuf, binary.BigEndian, &msgtail.Reserved)
	return &msgtail
}

func (c *Connection) SerializeTail(task *sendTask) []byte {
	s := c.streamProtocol
	msgtail := task.msg.MsgTail
	//startpos := 0
	s.serializeBuf.Reset()
	binary.Write(s.serializeBuf, binary.BigEndian, &msgtail.Crc)
	binary.Write(s.serializeBuf, binary.BigEndian, &msgtail.Reserved)
	return s.serializeBuf.Bytes()
}
