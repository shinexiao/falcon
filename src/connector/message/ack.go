package message

import (
	"connector/protocol"
	"connector/connection"
)

/**
 * ack 消息确认
 */
type AckMessage struct {
	ByteBufMessage
}

func AckMessageFrom(src Message) *AckMessage {

	packet := &protocol.Packet{
		Cmd:       protocol.CMD_ACK,
		SessionId: src.Packet().SessionId,
	}

	msg := &AckMessage{}
	msg.packet = packet
	msg.BaseMessage.Child = msg

	return msg
}

func NewAckMessage(packet *protocol.Packet, conn *connection.Conn) *AckMessage {

	data := &protocol.Packet{
		Cmd:       protocol.CMD_ACK,
		SessionId: packet.SessionId,
	}

	msg := &AckMessage{}
	msg.packet = data
	msg.BaseMessage.Child = msg
	msg.conn = conn

	return msg
}

//解码消息
func (me *AckMessage) Decode(body []byte) {

}

//编码消息
func (me *AckMessage) Encode() []byte {

	return nil
}
