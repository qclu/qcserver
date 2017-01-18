package tcpnetwork

type IEventQueue interface {
	Push(*ConnEvent)
	Pop() *ConnEvent
}

type IStreamProtocol interface {
	//	Init
	Init()
	//	get the header length of the stream
	GetHeaderLength() int
	//	read the header length of the stream
	UnserializeHeader(interface{}) interface{}
	//	format header
	SerializeHeader(interface{}) []byte
}

type IEventHandler interface {
	OnConnected(evt *ConnEvent)
	OnDisconnected(evt *ConnEvent)
	OnRecv(evt *ConnEvent)
}

type IUnpacker interface {
	Unpack(*Connection, []byte) ([]byte, error)
}
