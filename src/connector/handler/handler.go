package handler

import (
	"falcon/src/connector/protocol"
	"falcon/src/connector/connection"
)
type MessageHandler interface {
	Handle(packet *protocol.Packet, conn *connection.Conn) //处理消息
}
