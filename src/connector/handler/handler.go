package handler

import (
	"connector/protocol"
	"connector/connection"
)

type MessageHandler interface {
	Handle(packet *protocol.Packet, conn *connection.Conn) //处理消息
}
