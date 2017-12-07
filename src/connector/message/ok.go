package message

import (
	"bytes"
	"connector/protocol"
	"connector/connection"
)

type OkMessage struct {
	ByteBufMessage
	Cmd  int8
	Code int8
	Data string
}

/**
 * 消息转换
 */
func OkMessageFrom(msg Message) *OkMessage {
	packet := msg.Packet()
	cmd := packet.Cmd

	packet.Cmd = protocol.CMD_OK

	return NewOkMessage(cmd, packet, msg.Conn())
}

/**
 * 创建ok消息
 */
func NewOkMessage(cmd int8, packet *protocol.Packet, conn *connection.Conn) *OkMessage {

	msg := &OkMessage{}
	msg.packet = packet
	msg.conn = conn
	msg.BaseMessage.Child = msg
	msg.Cmd = cmd

	return msg
}

//解码消息
func (me *OkMessage) Decode(body []byte) {
	//输出消息看内容
	buf := new(bytes.Buffer)
	buf.Write(body)

	me.Cmd = me.DecodeInt8(buf)
	me.Code = me.DecodeInt8(buf)
	me.Data = me.DecodeString(buf)
}

//编码消息
func (me *OkMessage) Encode() []byte {

	buf := new(bytes.Buffer)

	me.EncodeInt8(buf, me.Cmd)
	me.EncodeInt8(buf, me.Code)
	me.EncodeString(buf, me.Data)

	return buf.Bytes()
}
