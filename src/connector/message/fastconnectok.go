package message

import (
	"bytes"
	"connector/protocol"
)

type FastConnectOkMessage struct {
	ByteBufMessage
	Heartbeat int32
}

func FastConnectOkMessageFrom(src Message) *FastConnectOkMessage {

	packet := src.Packet()
	packet.Cmd = protocol.CMD_FAST_CONNECT

	msg := &FastConnectOkMessage{}
	msg.BaseMessage.Child = msg
	msg.conn = src.Conn()
	msg.packet = packet
	
	return msg
}

//解码消息
func (me *FastConnectOkMessage) Decode(body []byte) {
	//输出消息看内容
	buf := new(bytes.Buffer)
	buf.Write(body)
	me.Heartbeat = me.DecodeInt32(buf)
}

//编码消息
func (me *FastConnectOkMessage) Encode() []byte {

	buf := new(bytes.Buffer)
	me.EncodeInt32(buf, me.Heartbeat)
	return buf.Bytes()
}
