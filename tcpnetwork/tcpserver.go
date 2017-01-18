package tcpnetwork

import (
	"qcserver/util/log"
)

type QcTcpServer struct {
	srvAddr  string
	network  *TCPNetwork
	stopFlag int32
	logger   *log.Log
}

// echo server routine
func NewTcpServer(addr string) *QcTcpServer {
	tcpSrv := &QcTcpServer{
		srvAddr:  addr,
		network:  NewTCPNetwork(1024, NewStreamProtocol()),
		stopFlag: 0,
		logger:   log.GetLog(),
	}
	return tcpSrv
}

func (srv *QcTcpServer) Listen() error {
	err := srv.network.Listen(srv.srvAddr)
	return err
}

func (srv *QcTcpServer) Start() {
	srv.logger.LogInfo("TcpServer start ...")
	//srv.network.ServeWithHandler(
}

func (srv *QcTcpServer) Destroy() {
	srv.network.Shutdown()
}

type MsgHandler struct {
}

func (handler *MsgHandler) OnRecv(evt *ConnEvent) {
}

func (handler *MsgHandler) OnConnected(evt *ConnEvent) {
}

func (handler *MsgHandler) OnDisconnected(evt *ConnEvent) {
}
