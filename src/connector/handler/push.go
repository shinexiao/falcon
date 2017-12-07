package handler

import (
	"connector/protocol"
	"connector/connection"
	"connector/message"
	"core"
	"github.com/gogap/logrus"
	"fmt"
)

type PushHandler struct {
}

func NewPushHandler(context core.FalconContext) *PushHandler {

	return &PushHandler{}
}

func (me *PushHandler) Handle(packet *protocol.Packet, conn *connection.Conn) {

	msg := me.decode(packet, conn)
	me.info(fmt.Sprintf("recive push message=%v", msg))

	if msg.AutoAck() {
		message.AckMessageFrom(msg).Send()
		me.info(fmt.Sprintf("send ack for push message=%v", msg))
	}

	//业务逻辑
	fmt.Println("push", string(msg.Content))
}

//解码出消息
func (me *PushHandler) decode(packet *protocol.Packet, conn *connection.Conn) *message.PushMessage {

	msg := message.NewPushMessage(packet, conn)
	msg.DecodeBody()

	return msg
}

//输出日志
func (me *PushHandler) info(str string) {
	logrus.
		WithField("system", "falcon").
		WithField("struct", "PushHandler").Info(str)
}

//输出日志
func (me *PushHandler) error(str string) {
	logrus.
		WithField("system", "falcon").
		WithField("struct", "PushHandler").Error(str)
}
