package server

import (
	"net"
	"github.com/gogap/logrus"
	"fmt"
	"connector/protocol"
	"bytes"
	"encoding/binary"
	"io"
	"connector/connection"
	"connector/dispatcher"
	"connector/handler"
	"strings"
	"core"
	"runtime"
)

const MaxPacketSize = 1024 * 1024

/**
 * 建立服务连接
 * 解析消息
 * 分发消息
 */
type TcpConnector struct {
	listeners         []*net.TCPListener
	Option            ConnectorOptions
	messageDispatcher *dispatcher.MessageDispatcher
	falcontext        core.FalconContext
}

func NewTcpConnector(context core.FalconContext, options ConnectorOptions) *TcpConnector {

	messageDispatcher := dispatcher.NewMessageDispatcher()
	messageDispatcher.Register(protocol.CMD_HANDSHAKE, handler.NewHandshakeHandler(context))      // 握手
	messageDispatcher.Register(protocol.CMD_HEARTBEAT, handler.NewHeartBeatHandler())             // 心跳
	messageDispatcher.Register(protocol.CMD_BIND, handler.NewBindUserHandler())                   // 绑定
	messageDispatcher.Register(protocol.CMD_UNBIND, handler.NewBindUserHandler())                 // 解绑
	messageDispatcher.Register(protocol.CMD_ACK, handler.NewAckHandler(context))                  //消息ack
	messageDispatcher.Register(protocol.CMD_FAST_CONNECT, handler.NewFastConnectHandler(context)) //快速重连
	messageDispatcher.Register(protocol.CMD_PUSH, handler.NewPushHandler(context))                //上行消息

	return &TcpConnector{
		messageDispatcher: messageDispatcher,
		Option:            options,
		falcontext:        context,
	}
}

//绑定服务端口
func (me *TcpConnector) Start(addrs []string) error {

	for _, bind := range addrs {

		addr, err := net.ResolveTCPAddr("tcp4", bind)
		if err != nil {
			me.error(fmt.Sprintf("net resolveTCPAddr tcp4 %s error %v", bind, err))
			return err
		}

		listener, err := net.ListenTCP("tcp4", addr);
		if err != nil {
			me.error(fmt.Sprintf("net ListenTCP tcp4 %s error %v", bind, err))
			return err
		}

		me.info(fmt.Sprintf("start tcp listen: %s", bind))

		for i := 0; i < runtime.NumCPU(); i++ {
			go me.acceptTcp(listener)
		}

	}

	return nil
}

//接收tcp连接
func (me *TcpConnector) acceptTcp(listener *net.TCPListener) {
	var (
		conn *net.TCPConn
		err  error
		r    int
	)
	for {
		if conn, err = listener.AcceptTCP(); err != nil {
			me.error(fmt.Sprintf("listener.Accept(%s) error(%v)", listener.Addr().String(), err))
			return
		}
		if err = conn.SetKeepAlive(me.Option.TCPKeepalive); err != nil {
			me.error(fmt.Sprintf("conn.SetKeepAlive() error(%v)", err))
			return
		}
		if err = conn.SetReadBuffer(me.Option.TCPRcvbuf); err != nil {
			me.error(fmt.Sprintf("conn.SetReadBuffer() error(%v)", err))
			return
		}
		if err = conn.SetWriteBuffer(me.Option.TCPSndbuf); err != nil {
			me.error(fmt.Sprintf("conn.SetWriteBuffer() error(%v)", err))
			return
		}
		//开启处理
		c := connection.NewConn(conn, true)
		go me.serve(c)
		//负载均衡
		if r++; r == 1000000 {
			r = 0
		}
		//todo 其他逻辑
		//开启心跳检测
	}
}

//关闭服务
func (me *TcpConnector) Close() {
	for _, listener := range me.listeners {
		listener.Close()
	}
}

//服务
func (me *TcpConnector) serve(conn *connection.Conn) {
	me.info(fmt.Sprintf("recive connection from :%s", conn.Channel.RemoteAddr()))

	//声明一个临时缓冲区，用来存储被截断的数据
	tempBuf := make([]byte, 0)
	buffer := make([]byte, 2048) //参数可配置

	readChan := make(chan *protocol.Packet)

	go func(ch chan *protocol.Packet) {
		for {
			select {
			case packet := <-ch:
				me.messageDispatcher.OnReceive(packet, conn)
			}
		}
	}(readChan)

	//声明一个管道用于接收解包的数据
	for {
		if n, err := conn.Channel.Read(buffer); err != nil {
			if err == io.EOF {
				conn.Close()
				break
			}
			if strings.Contains(err.Error(), "use of closed network connection") {
				conn.Close()
				break
			}
			me.error(fmt.Sprintf("connection from: %s , error: %v", conn.Channel.RemoteAddr(), err))
		} else {
			//循环读取数据
			tempBuf = append(tempBuf, buffer[:n]...)
			//解码心跳
			//tempBuf = me.decodeHeartbeat(tempBuf, readChan)
			//解码数据包
			tempBuf = me.decodeFrames(tempBuf, readChan)
		}
	}

}

//解包心跳
func (me *TcpConnector) decodeHeartbeat(buf []byte, ch chan *protocol.Packet) []byte {
	var hb_packet int8

	buffer := new(bytes.Buffer)
	buffer.Write(buf)

	binary.Read(buffer, binary.BigEndian, &hb_packet)
	if hb_packet == protocol.HB_PACKET_BYTE {
		ch <- protocol.NewPacket(protocol.CMD_HEARTBEAT)
		return buf[1:]
	}

	return buf
}

//解包消息
func (me *TcpConnector) decodeFrames(buf []byte, ch chan *protocol.Packet) []byte {

	bufSize := len(buf)

	if bufSize < protocol.PACKET_HEADER_LEN { //当前内容长度
		return buf
	}

	buffer := new(bytes.Buffer)
	buffer.Write(buf)

	var bodyLength int32
	binary.Read(buffer, binary.BigEndian, &bodyLength)

	total := protocol.PACKET_HEADER_LEN + int(bodyLength)

	if total > bufSize { // 消息内容不够
		return buf
	}

	//if bodyLength > MaxPacketSize { //todo 消息内容太长
	//	return
	//}
	//获取命令
	var cmd int8
	binary.Read(buffer, binary.BigEndian, &cmd)
	fmt.Println("cmd", cmd)
	//解码
	packet := protocol.NewPacket(cmd)
	protocol.DecodePacket(packet, buffer, int(bodyLength))

	ch <- packet

	return buf[total:]
}

//输出日志
func (me *TcpConnector) info(str string) {
	logrus.
		WithField("system", "falcon").
		WithField("struct", "TcpConnector").Info(str)
}

//输出日志
func (me *TcpConnector) error(str string) {
	logrus.
		WithField("system", "falcon").
		WithField("struct", "TcpConnector").Error(str)
}
