package main

import (
	"fmt"
	"os"
	"os/signal"
	"qcserver/tcpnetwork"
	"qcserver/util/log"
	"sync"
	"sync/atomic"
	"syscall"
)

var (
	kServerAddress  = "localhost:14444"
	serverConnected int32
	stopFlag        int32
)

// echo server routine
func qcMsgServer() (*tcpnetwork.TCPNetwork, error) {
	var err error
	server := tcpnetwork.NewTCPNetwork(1024, tcpnetwork.NewStreamProtocol())
	err = server.Listen(kServerAddress)
	if nil != err {
		return nil, err
	}

	return server, nil
}

func routineQcMsgServer(server *tcpnetwork.TCPNetwork, wg *sync.WaitGroup, stopCh chan struct{}) {
	defer func() {
		log.Println("QcMsg server done")
		wg.Done()
	}()

	for {
		select {
		case evt, ok := <-server.GetEventQueue():
			{
				if !ok {
					return
				}

				switch evt.EventType {
				case tcpnetwork.KConnEvent_Connected:
					{
						log.Println("Client ", evt.Conn.GetRemoteAddress(), " connected")
					}
				case tcpnetwork.KConnEvent_Close:
					{
						log.Println("Client ", evt.Conn.GetRemoteAddress(), " disconnected")
					}
				case tcpnetwork.KConnEvent_Data:
					{
						evt.Conn.Send(evt.Msg, 0)
					}
				}
			}
		case <-stopCh:
			{
				return
			}
		}
	}
}

func main() {
	logger, err := log.NewLog("/var/log/", "qcserver", 0)
	if err != nil {
		fmt.Println("Failed to init log module...")
		return
	}
	logger.LogInfo("log module start...")
	// create server
	server, err := qcMsgServer()
	if nil != err {
		log.Println(err)
		return
	}

	stopCh := make(chan struct{})

	// process event
	var wg sync.WaitGroup
	wg.Add(1)
	go routineQcMsgServer(server, &wg, stopCh)

	// wait
	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

MAINLOOP:
	for {
		select {
		case <-sc:
			{
				//	app cancelled by user , do clean up work
				log.Println("Terminating ...")
				break MAINLOOP
			}
		}
	}

	atomic.StoreInt32(&stopFlag, 1)
	log.Println("Press enter to exit")
	close(stopCh)
	wg.Wait()
	server.Shutdown()
}
