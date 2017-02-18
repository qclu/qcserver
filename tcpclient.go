package main

import (
	"bufio"
	"encoding/json"
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
	kServerAddress  = "118.178.188.139:14444"
	serverConnected int32
	stopFlag        int32
)

// echo client routine
func echoClient() (*tcpnetwork.TCPNetwork, *tcpnetwork.Connection, error) {
	var err error
	client := tcpnetwork.NewTCPNetwork(1024, tcpnetwork.NewStreamProtocol())
	conn, err := client.Connect(kServerAddress)
	if nil != err {
		return nil, nil, err
	}

	return client, conn, nil
}

func routineEchoClient(client *tcpnetwork.TCPNetwork, wg *sync.WaitGroup, stopCh chan struct{}) {
	defer func() {
		log.Println("client done")
		wg.Done()
	}()

EVENTLOOP:
	for {
		select {
		case evt, ok := <-client.GetEventQueue():
			{
				if !ok {
					return
				}
				switch evt.EventType {
				case tcpnetwork.KConnEvent_Connected:
					{
						log.Println("Press any thing")
						atomic.StoreInt32(&serverConnected, 1)
					}
				case tcpnetwork.KConnEvent_Close:
					{
						log.Println("Disconnected from server")
						atomic.StoreInt32(&serverConnected, 0)
						break EVENTLOOP
					}
				case tcpnetwork.KConnEvent_Data:
					{
						text := string(evt.Msg.MsgData)
						log.Println(evt.Conn.GetRemoteAddress(), ":", text)
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

func routineInput(wg *sync.WaitGroup, clientConn *tcpnetwork.Connection) {
	defer func() {
		log.Println("input done")
		wg.Done()
	}()

	reader := bufio.NewReader(os.Stdin)
	header := tcpnetwork.NewMsgHeader()
	header.Type = tcpnetwork.MSG_STAT
	header.TransactionId = 0
	var msg tcpnetwork.QcMessage
	msg.MsgHeader = header
	msg.MsgTail = tcpnetwork.NewMsgTail()
	for {
		line, _, _ := reader.ReadLine()
		if atomic.LoadInt32(&stopFlag) != 0 {
			return
		}
		str := string(line)

		if str == "\n" {
			continue
		}

		if atomic.LoadInt32(&serverConnected) != 1 {
			log.Println("Not connected")
			continue
		}
		statinfo := &tcpnetwork.StatInfo{
			DevName:   "dev1",
			Hospital:  "协和医院",
			Location:  str,
			PostCode:  "100010",
			HwVersion: "HW1111",
			SwVersion: "SW0101",
			WorkCnt:   1000,
		}
		var err error
		msg.MsgData, err = json.Marshal(statinfo)
		if err != nil {
			fmt.Println("Failed to marshal statinfo, err: ", err)
			continue
		}
		msg.MsgHeader.DataLen = uint32(len(msg.MsgData))
		msg.MsgHeader.SrvId = 1
		log.Println("Send data[", str, "], tid: ", msg.MsgHeader.TransactionId, ", len: ", msg.MsgHeader.DataLen)
		clientConn.Send(&msg, 0)
		msg.MsgHeader.TransactionId++
	}
}

func main() {
	logger, err := log.NewLog("/var/log/", "qcserver", 0)
	if err != nil {
		fmt.Println("Failed to init log module...")
		return
	}
	logger.LogInfo("log module start...")
	// create client
	client, clientConn, err := echoClient()
	if nil != err {
		log.Println(err)
		return
	}

	stopCh := make(chan struct{})

	// process event
	var wg sync.WaitGroup
	wg.Add(1)
	go routineEchoClient(client, &wg, stopCh)

	// input event
	wg.Add(1)
	go routineInput(&wg, clientConn)

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
	clientConn.Close()
}
