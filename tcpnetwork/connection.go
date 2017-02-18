package tcpnetwork

import (
	"errors"
	//"github.com/golang/protobuf/proto"
	"net"
	"qcserver/util/log"
	"sync/atomic"
	"time"
)

const (
	kConnStatus_None = iota
	kConnStatus_Connected
	kConnStatus_Disconnected
)

// All connection event
const (
	KConnEvent_None = iota
	KConnEvent_Connected
	KConnEvent_Disconnected
	KConnEvent_Data
	KConnEvent_Pb
	KConnEvent_Close
	KConnEvent_Total
)

const (
	kConnConf_DefaultSendTimeoutSec = 5
	kConnConf_MaxReadBufferLength   = 0xffff // 0xffff
)

// Send method flag
const (
	KConnFlag_CopySendBuffer = 1 << iota // do not copy the send buffer
	KConnFlag_NoHeader                   // do not append stream header
)

// send task
type sendTask struct {
	msg  *QcMessage
	flag int64
}

// Connection is a wrap for net.Conn and process read and write task of the conn
// When event occurs, it will call the eventQueue to dispatch event
type Connection struct {
	conn                net.Conn
	status              int32
	connId              int
	sendMsgQueue        chan *sendTask
	sendTimeoutSec      int
	eventQueue          IEventQueue
	streamProtocol      *StreamProtocol
	maxReadBufferLength uint32
	userdata            interface{}
	from                int
	readTimeoutSec      int
	fnSyncExecute       FuncSyncExecute
	unpacker            IUnpacker
	disableSend         int32
	localAddr           string
	remoteAddr          string
	logger              *log.Log
}

func newConnection(c net.Conn, sendBufferSize int, eq IEventQueue) *Connection {
	return &Connection{
		conn:                c,
		status:              kConnStatus_None,
		connId:              0,
		sendMsgQueue:        make(chan *sendTask, sendBufferSize),
		sendTimeoutSec:      kConnConf_DefaultSendTimeoutSec,
		maxReadBufferLength: kConnConf_MaxReadBufferLength,
		eventQueue:          eq,
		logger:              log.GetLog(),
	}
}

// ConnEvent represents a event occurs on a connection, such as connected, disconnected or data arrived
type ConnEvent struct {
	EventType int
	Conn      *Connection
	Msg       *QcMessage
	Userdata  interface{}
	//PbM       proto.Message
}

func newConnEvent(et int, c *Connection, msg *QcMessage) *ConnEvent {
	return &ConnEvent{
		EventType: et,
		Conn:      c,
		Msg:       msg,
	}
}

func (c *Connection) init() {
	c.localAddr = c.conn.LocalAddr().String()
	c.remoteAddr = c.conn.RemoteAddr().String()
}

//	directly close, packages in queue will not be sent
func (c *Connection) close() {
	//	set the disconnected status, use atomic operation to avoid close twice
	if atomic.CompareAndSwapInt32(&c.status, kConnStatus_Connected, kConnStatus_Disconnected) {
		c.conn.Close()
	}
}

// Close the connection, routine safe, send task in the queue will be sent before closing the connection
func (c *Connection) Close() {
	if atomic.LoadInt32(&c.status) != kConnStatus_Connected {
		return
	}

	select {
	case c.sendMsgQueue <- nil:
		{
			//	nothing
		}
	case <-time.After(time.Duration(c.sendTimeoutSec)):
		{
			//	timeout, close the connection
			c.close()
		}
	}

	//	disable send
	atomic.StoreInt32(&c.disableSend, 1)
}

//	When don't need conection to send any thing, free it, DO NOT call it on multi routines
func (c *Connection) Free() {
	if nil != c.sendMsgQueue {
		close(c.sendMsgQueue)
		c.sendMsgQueue = nil
	}
}

func (c *Connection) syncExecuteEvent(evt *ConnEvent) bool {
	if nil == c.fnSyncExecute {
		return false
	}

	return c.fnSyncExecute(evt)
}

func (c *Connection) pushEvent(et int, msg *QcMessage) {
	//	this is for sync execute
	evt := newConnEvent(et, c, msg)
	if c.syncExecuteEvent(evt) {
		return
	}

	if nil == c.eventQueue {
		panic("Nil event queue")
		return
	}
	c.eventQueue.Push(evt)
}

//func (c *Connection) pushPbEvent(pb proto.Message) {
//	if nil == c.eventQueue {
//		panic("Nil event queue")
//		return
//	}
//	c.eventQueue.Push(&ConnEvent{
//		EventType: KConnEvent_Data,
//		Conn:      c,
//		PbM:       pb,
//	})
//}

// SetSyncExecuteFunc , you can set a callback that you can synchoronously process the event in every connection's event routine
// If the callback function return true, the event will not be dispatched
func (c *Connection) SetSyncExecuteFunc(fn FuncSyncExecute) FuncSyncExecute {
	prevFn := c.fnSyncExecute
	c.fnSyncExecute = fn
	return prevFn
}

// GetStatus get the connection's status
func (c *Connection) GetStatus() int32 {
	return c.status
}

func (c *Connection) setStatus(stat int) {
	c.status = int32(stat)
}

// GetConnId get the connection's id
func (c *Connection) GetConnId() int {
	return c.connId
}

// SetConnId set the connection's id
func (c *Connection) SetConnId(id int) {
	c.connId = id
}

// GetConn get the raw net.Conn interface
func (c *Connection) GetConn() net.Conn {
	return c.conn
}

// GetUserdata get the userdata you set
func (c *Connection) GetUserdata() interface{} {
	return c.userdata
}

// SetUserdata set the userdata you need
func (c *Connection) SetUserdata(ud interface{}) {
	c.userdata = ud
}

// SetReadTimeoutSec set the read deadline for the connection
func (c *Connection) SetReadTimeoutSec(sec int) {
	c.readTimeoutSec = sec
}

//  GetReadTimeoutSec get the read deadline for the connection
func (c *Connection) GetReadTimeoutSec() int {
	return c.readTimeoutSec
}

// GetRemoteAddress return the remote address of the connection
func (c *Connection) GetRemoteAddress() string {
	return c.remoteAddr
}

// GetLocalAddress return the local address of the connection
func (c *Connection) GetLocalAddress() string {
	return c.localAddr
}

// SetUnpacker you can set a custom binary stream unpacker on the connection
//func (c *Connection) SetUnpacker(unpacker IUnpacker) {
//	c.unpacker = unpacker
//}

// GetUnpacker you can get the unpacker you set
//func (c *Connection) GetUnpacker() IUnpacker {
//	return c.unpacker
//}

func (c *Connection) setStreamProtocol(sp *StreamProtocol) {
	c.streamProtocol = sp
}

func (c *Connection) sendRaw(task *sendTask) error {
	if atomic.LoadInt32(&c.disableSend) != 0 {
		return ErrConnIsClosed
	}
	if atomic.LoadInt32(&c.status) != kConnStatus_Connected {
		return ErrConnIsClosed
	}

	select {
	case c.sendMsgQueue <- task:
		{
			//	nothing
		}
	case <-time.After(time.Duration(c.sendTimeoutSec)):
		{
			//	timeout, close the connection
			c.close()
			return ErrConnSendTimeout
		}
	}

	return nil
}

// ApplyReadDealine
func (c *Connection) ApplyReadDeadline() {
	if 0 != c.readTimeoutSec {
		c.conn.SetReadDeadline(time.Now().Add(time.Duration(c.readTimeoutSec) * time.Second))
	}
}

// ResetReadDeadline
func (c *Connection) ResetReadDeadline() {
	c.conn.SetReadDeadline(time.Time{})
}

// Send the buffer
func (c *Connection) Send(msg *QcMessage, f int64) error {
	task := &sendTask{
		msg:  msg,
		flag: f,
	}
	//buf := msg.MsgData
	//log.Println(evt.Conn.GetRemoteAddress(), ":", text)
	//	copy send buffer
	//if 0 != f&KConnFlag_CopySendBuffer {
	//	msgCopy := make([]byte, len(msg.MsgData))
	//	copy(msgCopy, msg.MsgData)
	//	buf = msgCopy
	//	task.data = buf
	//}

	return c.sendRaw(task)
}

// SendPb send protocol buffer message
//func (c *Connection) SendPb(pb proto.Message, f int64) error {
//	var pbd []byte
//	var err error
//
//	if pbd, err = proto.Marshal(pb); nil != err {
//		return err
//	}
//
//	task := sendTask{
//		data: pbd,
//		flag: f,
//	}
//
//	//	send protobuf
//	return c.sendRaw(&task)
//}

//	run a routine to process the connection
func (c *Connection) run() {
	go c.routineMain()
}

func (c *Connection) routineMain() {
	defer func() {
		//routine end
		e := recover()
		if e != nil {
			c.logger.LogError("Panic", e)
		}

		//	close the connection
		c.close()

		//	free channel
		//	FIXED : consumers need free it, not producer

		//	post event
		c.pushEvent(KConnEvent_Disconnected, nil)
	}()

	if nil == c.streamProtocol {
		panic("Nil stream protocol")
		return
	}
	c.StreamProtocolInit()

	//	connected
	c.pushEvent(KConnEvent_Connected, nil)
	atomic.StoreInt32(&c.status, kConnStatus_Connected)

	go c.routineSend()
	c.routineRead()
}

func (c *Connection) routineSend() error {
	defer func() {
		e := recover()
		if nil != e {
			//	panic
			c.logger.LogError("Panic", e)
		}
	}()

	for {
		select {
		case evt, ok := <-c.sendMsgQueue:
			{
				if !ok {
					//	channel closed, quit
					return nil
				}

				if nil == evt {
					c.close()
					return nil
				}

				var err error

				//if 0 == evt.flag&KConnFlag_NoHeader {
				headerBytes := c.SerializeHeader(evt)
				if nil != headerBytes {
					//	write header first
					if len(headerBytes) != 0 {
						_, err = c.conn.Write(headerBytes)
						if err != nil {
							return err
						}
					}
				} else {
					//	invalid packet
					panic("Failed to serialize header")
					break
				}
				//}

				_, err = c.conn.Write(evt.msg.MsgData)
				if err != nil {
					return err
				}

				tailBytes := c.SerializeTail(evt)
				if nil != tailBytes {
					//	write header first
					if len(tailBytes) != 0 {
						_, err = c.conn.Write(tailBytes)
						if err != nil {
							return err
						}
					}
				} else {
					//	invalid packet
					panic("Failed to serialize tail")
					break
				}
			}
		}
	}

	return nil
}

func (c *Connection) routineRead() error {
	for {
		msg, err := c.unpack()
		if err != nil {
			return err
		}

		c.logger.LogInfo("receive msg header[tid: ", msg.MsgHeader.TransactionId, ", type: ", msg.MsgHeader.Type, ", DataLen: ", msg.MsgHeader.DataLen, "]")
		c.logger.LogInfo("receive msg data[", string(msg.MsgData), "]")
		//	only push event when the connection is connected
		if atomic.LoadInt32(&c.status) == kConnStatus_Connected {
			//	try to unserialize to pb. do it in each routine to reduce the pressure of worker routine
			c.pushEvent(KConnEvent_Data, msg)
		}
	}

	return nil
}

func (c *Connection) unpack() (msg *QcMessage, err error) {
	//	read head
	c.ApplyReadDeadline()
	msg = NewMessage()
	headBuf := make([]byte, ProtocolMsgHeaderLength)
	rc, err := c.conn.Read(headBuf)
	if err != nil {
		return nil, err
	}
	if rc != ProtocolMsgHeaderLength {
		c.logger.LogError("Failed to read msg header, readlen: ", rc)
		return nil, errors.New("Failed to read msg header")
	}
	//	check length
	msg.MsgHeader = c.UnserializeHeader(headBuf)
	c.logger.LogInfo("header datalen: ", msg.MsgHeader.DataLen)
	packetLength := msg.MsgHeader.DataLen + ProtocolMsgTailLength + ProtocolMsgHeaderLength
	if packetLength > c.maxReadBufferLength {
		c.logger.LogError("The stream data is too long")
		return nil, errors.New("The stream data is too long")
	}

	//	read body
	c.ApplyReadDeadline()
	bodyLength := msg.MsgHeader.DataLen
	msg.MsgData = make([]byte, bodyLength)
	_, err = c.conn.Read(msg.MsgData)
	if err != nil {
		return nil, err
	}

	// read tail
	c.ApplyReadDeadline()
	tailbuf := make([]byte, ProtocolMsgTailLength)
	rc, err = c.conn.Read(tailbuf)
	if err != nil {
		return nil, err
	}
	if rc != ProtocolMsgTailLength {
		c.logger.LogError("Failed to read msg tail, readlen: ", rc)
		return nil, errors.New("Failed to read msg tail")
	}
	msg.MsgTail = c.UnserializeTail(tailbuf)

	c.ResetReadDeadline()
	return msg, nil
}
