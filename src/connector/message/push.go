package message

import (
	"connector/protocol"
	"connector/connection"
)

type PushMessage struct {
	BaseMessage
	Content []byte
}

/**
 * 创建推送消息
 */
func NewPushMessage(packet *protocol.Packet, conn *connection.Conn) *PushMessage {

	msg := &PushMessage{}
	msg.packet = packet
	msg.conn = conn
	msg.BaseMessage.Child = msg

	return msg
}

//解码消息
func (me *PushMessage) Decode(body []byte) {
	//输出消息看内容
	me.Content = body
}

//编码消息
func (me *PushMessage) Encode() []byte {

	return me.Content
}

//是否自动ack消息
func (me *PushMessage) AutoAck() bool {

	return me.Packet().HasFlag(protocol.FLAG_AUTO_ACK)
}

func (me *PushMessage) BizAck() bool {
	return me.Packet().HasFlag(protocol.FLAG_BIZ_ACK)
}

//添加标识
func (me *PushMessage) AddFlag(flag int8) *PushMessage {
	me.Packet().AddFlag(flag)
	return me
}
