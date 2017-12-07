package message

import (
	"bytes"
	"connector/connection"
)

type FastConnectMessage struct {
	ByteBufMessage
	SessionId    string
	DeviceId     string
	MinHeartbeat int32
	MaxHeartbeat int32
}

func NewFastConnectMessage(conn *connection.Conn) *FastConnectMessage {

	msg := &FastConnectMessage{}
	msg.conn = conn
	msg.BaseMessage.Child = msg

	return msg
}

//解码消息
func (me *FastConnectMessage) Decode(body []byte) {
	//输出消息看内容
	buf := new(bytes.Buffer)
	buf.Write(body)

	me.SessionId = me.DecodeString(buf)
	me.DeviceId = me.DecodeString(buf)
	me.MinHeartbeat = me.DecodeInt32(buf)
	me.MaxHeartbeat = me.DecodeInt32(buf)
}

//编码消息
func (me *FastConnectMessage) Encode() []byte {

	buf := new(bytes.Buffer)

	me.EncodeString(buf, me.SessionId)
	me.EncodeString(buf, me.DeviceId)
	me.EncodeInt32(buf, me.MinHeartbeat)
	me.EncodeInt32(buf, me.MaxHeartbeat)

	return buf.Bytes()
}
