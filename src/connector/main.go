package main

import (
	"connector/server"
	"common"
	"github.com/gogap/logrus"
)

func main() {

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)

	options := server.ConnectorOptions{
		TCPKeepalive: true,
		TCPRcvbuf:    2048,
		TCPSndbuf:    2048,
	}

	falconServer := server.NewFalconServer()

	srv := server.NewTcpConnector(falconServer, options)

	srv.Start([]string{":3000"})

	//增加信号监听
	common.InitSignal()
}
