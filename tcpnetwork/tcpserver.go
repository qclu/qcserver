package tcpnetwork

import (
	"encoding/json"
	"errors"
	"qcserver/models"
	"qcserver/util/log"
	"strings"
	"time"
)

type QcTcpServer struct {
	srvAddr  string
	network  *TCPNetwork
	stopFlag int32
	logger   *log.Log
	dbSync   *models.DBSync
}

// echo server routine
func NewTcpServer(addr string) *QcTcpServer {
	tcpSrv := &QcTcpServer{
		srvAddr:  addr,
		network:  NewTCPNetwork(1024, NewStreamProtocol()),
		stopFlag: 0,
		logger:   log.GetLog(),
		dbSync:   models.GetDBSync(),
	}
	return tcpSrv
}

func (srv *QcTcpServer) Listen() error {
	err := srv.network.Listen(srv.srvAddr)
	return err
}

func (srv *QcTcpServer) Start() {
	srv.logger.LogInfo("TcpServer start ...")
	srv.network.ServeWithHandler(srv)
}

func (srv *QcTcpServer) Destroy() {
	srv.network.Shutdown()
}

func (srv *QcTcpServer) OnRecv(evt *ConnEvent) {
	srv.logger.LogInfo("receive msg from ", evt.Conn.GetRemoteAddress(),
		"msgheader[type: ", evt.Msg.MsgHeader.Type, ", tid: ", evt.Msg.MsgHeader.TransactionId,
		", datalen: ", evt.Msg.MsgHeader.DataLen, "]")
	err := srv.logHandler(evt.Msg)
	if err != nil {
		srv.logger.LogError("Failed to handle log process, err: ", err)
		return
	}
	switch evt.Msg.MsgHeader.Type {
	case MSG_STAT:
		srv.logger.LogInfo("MSG_STAT received[datalen: ", evt.Msg.MsgHeader.DataLen, ", tid: ",
			evt.Msg.MsgHeader.TransactionId, ", data: ", string(evt.Msg.MsgData), "]")
		srv.devstatHandler(evt)
	default:
		srv.logger.LogError("Invalid msg type: ", evt.Msg.MsgHeader.Type)
	}
	return
}

func (srv *QcTcpServer) devstatHandler(evt *ConnEvent) error {
	devrel, err := srv.dbSync.GetQcDevRelWithId(evt.Msg.MsgHeader.SrvId)
	if err != nil {
		srv.logger.LogError("Failed to get devrel with id ", evt.Msg.MsgHeader.SrvId, ", err: ", err)
		return err
	}
	var stat StatInfo
	err = json.Unmarshal(evt.Msg.MsgData, &stat)
	if err != nil {
		srv.logger.LogError("Failed to unmarshal stat info from data: ", string(evt.Msg.MsgData), ", err: ", err)
		return err
	}
	srv.logger.LogInfo("location: ", stat.Location)
	gis := strings.Split(stat.Location, ",")
	if len(gis) != 2 {
		srv.logger.LogError("Invalid gis info: ", stat.Location)
		return errors.New("Invalid location info of dev")
	}
	srv.logger.LogInfo("Longitude: ", gis[0], ", Longitude: ", gis[1])
	devstat := &models.QcDevStat{
		Dev:       devrel,
		Latitude:  gis[0],
		Longitude: gis[1],
		WorkCnt:   stat.WorkCnt,
		Reported:  time.Now().Format(models.TIME_FMT),
	}
	err = srv.dbSync.InsertQcDevStat(devstat)
	if err != nil {
		srv.logger.LogError("Failed to insert devstat into database, err: ", err)
		return err
	}
	return nil
}

func (srv *QcTcpServer) logHandler(msg *QcMessage) (err error) {
	devrel, err := srv.dbSync.GetQcDevRelWithId(msg.MsgHeader.SrvId)
	if err != nil {
		srv.logger.LogError("Failed to get devrel with id: ", msg.MsgHeader.SrvId, ", err: ", err)
		return err
	}
	_, err = models.CreateQcDevLog(srv.dbSync, msg.MsgHeader.Type, string(msg.MsgData), devrel)
	if err != nil {
		srv.logger.LogError("Failed to create devlog, err: ", err)
		return err
	}
	return nil
}

func (srv *QcTcpServer) OnConnected(evt *ConnEvent) {
	srv.logger.LogInfo("Connect from ", evt.Conn.GetRemoteAddress())
}

func (srv *QcTcpServer) OnDisconnected(evt *ConnEvent) {
	srv.logger.LogInfo("Disconnect from ", evt.Conn.GetRemoteAddress())
}
