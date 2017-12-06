package handler

import (
	"connector/protocol"
	"connector/connection"
	"github.com/gogap/logrus"
	"connector/message"
)

/**
 * 用户绑定
 */
type BindUserHandler struct {
	// 增加验证机制
}

func NewBindUserHandler() *BindUserHandler {

	return &BindUserHandler{}
}

func (me *BindUserHandler) Handle(packet *protocol.Packet, conn *connection.Conn) {

	msg := me.decode(packet, conn) //消息

	if packet.Cmd == protocol.CMD_BIND {
		me.bind(msg)
	} else {
		me.unbind(msg)
	}

}

//解码出消息
func (me *BindUserHandler) decode(packet *protocol.Packet, conn *connection.Conn) *message.BindUserMessage {

	msg := message.NewBindUserMessage(conn)
	msg.DecodeBody()

	return msg
}

//绑定消息
func (me *BindUserHandler) bind(bindUser *message.BindUserMessage)  {

}

//解绑消息
func (me *BindUserHandler) unbind(bindUser *message.BindUserMessage)  {

}

//输出日志
func (me *BindUserHandler) debug(str string) {
	logrus.
		WithField("system", "falcon").
		WithField("struct", "BindUserHandler").Debug(str)
}

//输出日志
func (me *BindUserHandler) info(str string) {
	logrus.
		WithField("system", "falcon").
		WithField("struct", "BindUserHandler").Info(str)
}

//输出日志
func (me *BindUserHandler) error(str string) {
	logrus.
		WithField("system", "falcon").
		WithField("struct", "BindUserHandler").Error(str)
}
