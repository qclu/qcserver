package tcpnetwork

import (
	"net"
	"qcserver/util/log"
	"sync/atomic"
	"time"
)

const (
	kServerConf_SendBufferSize = 1024
	kServerConn                = 0
	kClientConn                = 1
)

type TCPNetworkConf struct {
	SendBufferSize int
}

type TCPNetwork struct {
	streamProtocol  *StreamProtocol
	eventQueue      chan *ConnEvent
	rplQueue        chan *ConnEvent
	listener        net.Listener
	Conf            TCPNetworkConf
	connIdForServer int
	connIdForClient int
	connsForServer  map[int]*Connection
	connsForClient  map[int]*Connection
	shutdownFlag    int32
	readTimeoutSec  int
	logger          *log.Log
}

func NewTCPNetwork(eventQueueSize int, sp *StreamProtocol) *TCPNetwork {
	s := &TCPNetwork{}
	s.eventQueue = make(chan *ConnEvent, eventQueueSize)
	s.rplQueue = make(chan *ConnEvent, eventQueueSize)
	s.streamProtocol = sp
	s.connsForServer = make(map[int]*Connection)
	s.connsForClient = make(map[int]*Connection)
	s.shutdownFlag = 0
	//	default config
	s.Conf.SendBufferSize = kServerConf_SendBufferSize
	s.logger = log.GetLog()
	return s
}

// Push implements the IEventQueue interface
func (t *TCPNetwork) PushRpl(evt *ConnEvent) {
	if nil == t.rplQueue {
		return
	}

	//	push timeout
	select {
	case t.rplQueue <- evt:
		{

		}
	case <-time.After(time.Second * 5):
		{
			evt.Conn.close()
		}
	}

}

// Push implements the IEventQueue interface
func (t *TCPNetwork) Push(evt *ConnEvent) {
	if nil == t.eventQueue {
		return
	}

	//	push timeout
	select {
	case t.eventQueue <- evt:
		{

		}
	case <-time.After(time.Second * 5):
		{
			evt.Conn.close()
		}
	}

}

// Pop the event in event queue
func (t *TCPNetwork) PopRpl() *ConnEvent {
	evt, ok := <-t.rplQueue
	if !ok {
		//	event queue already closed
		return nil
	}

	return evt
}

// Pop the event in event queue
func (t *TCPNetwork) Pop() *ConnEvent {
	evt, ok := <-t.eventQueue
	if !ok {
		//	event queue already closed
		return nil
	}

	return evt
}

// GetEventQueue get the event queue channel
func (t *TCPNetwork) GetRplQueue() <-chan *ConnEvent {
	return t.eventQueue
}

// GetEventQueue get the event queue channel
func (t *TCPNetwork) GetEventQueue() <-chan *ConnEvent {
	return t.eventQueue
}

// Listen an address to accept client connection
func (t *TCPNetwork) Listen(addr string) error {
	t.logger.LogInfo("Server start listening...")
	ls, err := net.Listen("tcp", addr)
	if nil != err {
		return err
	}

	//	accept
	t.listener = ls
	go t.acceptRoutine()
	return nil
}

// Connect the remote server
func (t *TCPNetwork) Connect(addr string) (*Connection, error) {
	conn, err := net.Dial("tcp", addr)
	if nil != err {
		return nil, err
	}

	connection := t.createConn(conn)
	connection.from = 1
	connection.run()
	connection.init()

	return connection, nil
}

func (t *TCPNetwork) GetStreamProtocol() *StreamProtocol {
	return t.streamProtocol
}

func (t *TCPNetwork) SetStreamProtocol(sp *StreamProtocol) {
	t.streamProtocol = sp
}

func (t *TCPNetwork) GetReadTimeoutSec() int {
	return t.readTimeoutSec
}

func (t *TCPNetwork) SetReadTimeoutSec(sec int) {
	t.readTimeoutSec = sec
}

func (t *TCPNetwork) DisconnectAllConnectionsServer() {
	for k, c := range t.connsForServer {
		c.Close()
		delete(t.connsForServer, k)
	}
}

func (t *TCPNetwork) DisconnectAllConnectionsClient() {
	for k, c := range t.connsForClient {
		c.Close()
		delete(t.connsForClient, k)
	}
}

// Shutdown frees all connection and stop the listener
func (t *TCPNetwork) Shutdown() {
	if !atomic.CompareAndSwapInt32(&t.shutdownFlag, 0, 1) {
		return
	}

	//	stop accept routine
	if nil != t.listener {
		t.listener.Close()
	}

	//	close all connections
	t.DisconnectAllConnectionsClient()
	t.DisconnectAllConnectionsServer()
}

func (t *TCPNetwork) createConn(c net.Conn) *Connection {
	conn := newConnection(c, t.Conf.SendBufferSize, t)
	conn.setStreamProtocol(t.streamProtocol)
	return conn
}

// ServeWithHandler process all events in the event queue and dispatch to the IEventHandler
func (t *TCPNetwork) ServeWithHandler(handler IEventHandler) {
SERVE_LOOP:
	for {
		select {
		case evt, ok := <-t.eventQueue:
			{
				if !ok {
					//	channel closed or shutdown
					break SERVE_LOOP
				}

				t.handleEvent(evt, handler)
			}
		}
	}
}

func (t *TCPNetwork) acceptRoutine() {
	// after accept temporary failure, enter sleep and try again
	var tempDelay time.Duration

	t.logger.LogInfo("Start listening routine...")
	time.Sleep(10 * time.Second)
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			// check if the error is an temporary error
			if acceptErr, ok := err.(net.Error); ok && acceptErr.Temporary() {
				if 0 == tempDelay {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}

				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}

				t.logger.LogWarn("Accept error %s , retry after %d ms", acceptErr.Error(), tempDelay)
				time.Sleep(tempDelay)
				continue
			}

			t.logger.LogError("accept routine quit.error:%s", err.Error())
			t.listener = nil
			return
		}

		//	process conn event
		connection := t.createConn(conn)
		connection.SetReadTimeoutSec(t.readTimeoutSec)
		connection.from = kServerConn
		connection.init()
		connection.run()
	}
}

func (t *TCPNetwork) handleEvent(evt *ConnEvent, handler IEventHandler) {
	switch evt.EventType {
	case KConnEvent_Connected:
		{
			//	add to connection map
			connId := 0
			if kServerConn == evt.Conn.from {
				connId = t.connIdForServer + 1
				t.connIdForServer = connId
				t.connsForServer[connId] = evt.Conn
			} else {
				connId = t.connIdForClient + 1
				t.connIdForClient = connId
				t.connsForClient[connId] = evt.Conn
			}
			evt.Conn.connId = connId

			handler.OnConnected(evt)
		}
	case KConnEvent_Disconnected:
		{
			handler.OnDisconnected(evt)

			//	remove from connection map
			if kServerConn == evt.Conn.from {
				delete(t.connsForServer, evt.Conn.connId)
			} else {
				delete(t.connsForClient, evt.Conn.connId)
			}
		}
	case KConnEvent_Data:
		{
			handler.OnRecv(evt)
		}
	}
}
