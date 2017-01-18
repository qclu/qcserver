package main

import (
	"fmt"
	"os"
	"os/signal"
	"qcserver/tcpnetwork"
	"qcserver/util/log"
	"syscall"
)

var (
	kServerAddress  = "localhost:14444"
	serverConnected int32
	stopFlag        int32
)

func main() {
	logger, err := log.NewLog("/var/log/", "qcserver", 0)
	if err != nil {
		fmt.Println("Failed to init log module...")
		return
	}
	logger.LogInfo("log module start...")
	// create server
	server := tcpnetwork.NewTcpServer(kServerAddress)

	err = server.Listen()
	if err != nil {
		logger.LogError("Failed to start server listening routine, err:", err)
		return
	}
	server.Start()

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

	log.Println("Press enter to exit")
	server.Destroy()
}
