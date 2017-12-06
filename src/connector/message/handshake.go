package message

import (
	"connector/protocol"
	"bytes"
	"connector/connection"
)

/**
 * 握手消息
 */
type HandshakeMessage struct {
	ByteBufMessage
	DeviceId      string
	OsName        string
	OsVersion     string
	ClientVersion string
	Iv            []byte
	ClientKey     []byte
	MinHeartbeat  int32
	MaxHeartbeat  int32
	Timestamp     int64
}

//创建消息
func NewHandshakeMessage(conn *connection.Conn) *HandshakeMessage {

	msg := &HandshakeMessage{}

	packet := protocol.NewPacket(protocol.CMD_HANDSHAKE)
	packet.SessionId = msg.GenSessionId()

	msg.packet = packet
	msg.conn = conn
	msg.BaseMessage.Child = msg
	
	return msg
}

//创建消息
func NewHandshakeMessagePacket(packet *protocol.Packet, conn *connection.Conn) *HandshakeMessage {
	msg := &HandshakeMessage{}
	msg.packet = packet
	msg.conn = conn
	msg.BaseMessage.Child = msg

	return msg
}

//解码消息
func (me *HandshakeMessage) Decode(body []byte) {
	//输出消息看内容
	buf := new(bytes.Buffer)
	buf.Write(body)

	me.DeviceId = me.DecodeString(buf)
	me.OsName = me.DecodeString(buf)
	me.OsVersion = me.DecodeString(buf)
	me.ClientVersion = me.DecodeString(buf)
	me.Iv = me.DecodeBytes(buf)
	me.ClientKey = me.DecodeBytes(buf)
	me.MinHeartbeat = me.DecodeInt32(buf)
	me.MaxHeartbeat = me.DecodeInt32(buf)
	me.Timestamp = me.DecodeInt64(buf)
}

//编码消息
func (me *HandshakeMessage) Encode() []byte {

	buf := new(bytes.Buffer)
	me.EncodeString(buf, me.DeviceId)
	me.EncodeString(buf, me.OsName)
	me.EncodeString(buf, me.OsVersion)
	me.EncodeString(buf, me.ClientVersion)
	me.EncodeBytes(buf, me.Iv)
	me.EncodeBytes(buf, me.ClientKey)
	me.EncodeInt32(buf, me.MinHeartbeat)
	me.EncodeInt32(buf, me.MaxHeartbeat)
	me.EncodeInt64(buf, me.Timestamp)

	return buf.Bytes()
}
