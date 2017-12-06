package message

import (
	"bytes"
	"connector/protocol"
	"connector/connection"
)

/**
 * 握手成功消息
 */
type HandshakeOkMessage struct {
	ByteBufMessage
	ServerKey  []byte
	Heartbeat  int32
	SessionId  string
	ExpireTime int64
}

func NewHandshakeOkMessage(packet *protocol.Packet, conn *connection.Conn) *HandshakeOkMessage {

	resp := &protocol.Packet{}
	resp.Cmd = protocol.CMD_HANDSHAKE
	resp.SessionId = packet.SessionId

	msg := &HandshakeOkMessage{}
	msg.packet = packet
	msg.conn = conn
	msg.BaseMessage.Child = msg

	return msg
}

//解码消息
func (me *HandshakeOkMessage) Decode(body []byte) {
	//输出消息看内容
	buf := new(bytes.Buffer)
	buf.Write(body)

	me.ServerKey = me.DecodeBytes(buf)
	me.Heartbeat = me.DecodeInt32(buf)
	me.SessionId = me.DecodeString(buf)
	me.ExpireTime = me.DecodeInt64(buf)
}

//编码消息
func (me *HandshakeOkMessage) Encode() []byte {

	buf := new(bytes.Buffer)

	me.EncodeBytes(buf, me.ServerKey)
	me.EncodeInt32(buf, me.Heartbeat)
	me.EncodeString(buf, me.SessionId)
	me.EncodeInt64(buf, me.ExpireTime)

	return buf.Bytes()
}
