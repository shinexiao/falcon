package message

import (
	"falcon/src/connector/protocol"
	"falcon/src/connector/connection"
	"bytes"
)

/**
 * 错误消息
 */
type ErrorMessage struct {
	ByteBufMessage
	Cmd   int8
	Code  int8
	Reson string
	Data  string
}

/**
 * 从消息转换
 */
func ErrorMessageFrom(msg Message) *ErrorMessage {

	packet := &protocol.Packet{
		Cmd:       protocol.CMD_ERROR,
		SessionId: msg.Packet().SessionId,
	}

	errMsg := &ErrorMessage{}
	errMsg.conn = msg.Conn()
	errMsg.Cmd = msg.Packet().Cmd
	errMsg.packet = packet
	errMsg.BaseMessage.Child = errMsg // 非常重要
	return errMsg
}

//创建错误消息
func NewErrorMessage(cmd int8, packet *protocol.Packet, conn *connection.Conn) *ErrorMessage {
	msg := &ErrorMessage{}
	msg.packet = packet
	msg.conn = conn
	msg.Cmd = cmd
	msg.BaseMessage.Child = msg // 非常重要
	return msg
}

//解码消息
func (me *ErrorMessage) Decode(body []byte) {
	//输出消息看内容
	buf := new(bytes.Buffer)
	buf.Write(body)

	me.Cmd = me.DecodeInt8(buf)
	me.Code = me.DecodeInt8(buf)
	me.Reson = me.DecodeString(buf)
	me.Data = me.DecodeString(buf)
}

//编码消息
func (me *ErrorMessage) Encode() []byte {

	buf := new(bytes.Buffer)
	me.EncodeInt8(buf, me.Cmd)
	me.EncodeInt8(buf, me.Code)
	me.EncodeString(buf, me.Reson)
	me.EncodeString(buf, me.Data)

	return buf.Bytes()
}
