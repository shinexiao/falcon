package handler

import (
	"connector/protocol"
	"connector/connection"
	"github.com/gogap/logrus"
)

/**
 * 心跳检测
 */
type HeartBeatHandler struct {
}

func NewHeartBeatHandler() *HeartBeatHandler {

	return &HeartBeatHandler{}
}

func (me *HeartBeatHandler) Handle(packet *protocol.Packet, conn *connection.Conn) {
	conn.Send(packet) //ping -> pong
	me.debug("ping -> pong")
}

//输出日志
func (me *HeartBeatHandler) debug(str string) {
	logrus.
		WithField("system", "falcon").
		WithField("struct", "HeartBeatHandler").Debug(str)
}
