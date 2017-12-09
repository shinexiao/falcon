package connection

import (
	"falcon/src/security"
	"net"
	"time"
	"falcon/src/connector/protocol"
	"bytes"
	"fmt"
)

var RSA_CIPHER = security.Create()

/**
 * 连接封装，加入上下文
 */
type Conn struct {
	Channel       net.Conn
	Context       *SessionContext
	LastReadTime  int64
	LastWriteTime int64
}

/**
 * 创建连接
 */
func NewConn(conn net.Conn, security bool) *Conn {

	c := &Conn{}
	c.Channel = conn
	c.Context = NewSessionContext()
	if security {
		c.Context.Cipher = RSA_CIPHER
	}

	return c
}

/**
 * 判断是否读超时
 */
func (me *Conn) IsReadTimeout() bool {

	return time.Now().Unix()-me.LastReadTime > int64(me.Context.Heartbeat)+1000
}

/**
 * 判断是否写超时
 */
func (me *Conn) IsWriteTimeout() bool {

	return time.Now().Unix()-me.LastWriteTime > int64(me.Context.Heartbeat)-1000
}

/**
 * 更新读时间
 */
func (me *Conn) UpdateLastReadTime() {
	me.LastReadTime = time.Now().Unix()
}

/**
 * 更新写时间
 */
func (me *Conn) UpdateLastWriteTime() {
	me.LastWriteTime = time.Now().Unix()
}

/**
 * 关闭连接
 */
func (me *Conn) Close() {
	me.Channel.Close()
}

//发送消息
func (me *Conn) Send(packet *protocol.Packet) {
	me.SendWithListener(packet, nil)
}

//发送消息，带回调监听
func (me *Conn) SendWithListener(packet *protocol.Packet, listener func(conn *Conn, err error)) {
	buf := new(bytes.Buffer)
	protocol.EncodePacket(packet, buf)
	_, err := me.Channel.Write(buf.Bytes())

	fmt.Println("SendWithListener", packet.Cmd, packet.ToString())

	if listener != nil {
		listener(me, err)
	}
}
