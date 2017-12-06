package server

import "time"

//连接选项
type ConnectorOptions struct {
	CliProto         int
	SvrProto         int
	HandshakeTimeout time.Duration //握手超时时间
	TCPKeepalive     bool          //tcp keepalive
	TCPRcvbuf        int           //tcp 接收buf大小
	TCPSndbuf        int           //tcp 发送buf大小
}

type Connector interface {
	Start(addrs []string, option ConnectorOptions) error //绑定服务
	Close()                                             //关闭服务
}
