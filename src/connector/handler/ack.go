package handler

import (
	"connector/protocol"
	"connector/connection"
	"connector/message"
	"ack"
	"core"
	"github.com/gogap/logrus"
	"fmt"
)

/**
 * ack处理
 */
type AckHandler struct {
	ackTaskQueue *ack.AckTaskQueue
}

func NewAckHandler(ctx core.FalconContext) *AckHandler {

	ackHandler := &AckHandler{}
	ackHandler.ackTaskQueue = ctx.AckTaskQueue()

	return ackHandler
}

func (me *AckHandler) Handle(packet *protocol.Packet, conn *connection.Conn) {

	task := me.ackTaskQueue.GetAndRemove(int(packet.SessionId))

	if task == nil { //ack超时
		me.warn(fmt.Sprintf("ack time out , sessionId=%d", packet.SessionId))
		return
	}

	//直接返回ack确认消息
	message.AckMessageFrom(me.decode(packet, conn)).Send()

	me.info(fmt.Sprintf("ack success , sessionId=%d", packet.SessionId))

	//task.OnResponse() //成功收到客户的ACK响应
}

//解码出消息
func (me *AckHandler) decode(packet *protocol.Packet, conn *connection.Conn) *message.AckMessage {

	msg := message.NewAckMessage(packet, conn)

	return msg
}

//输出日志
func (me *AckHandler) debug(str string) {
	logrus.
		WithField("system", "falcon").
		WithField("struct", "AckHandler").Debug(str)
}

//输出日志
func (me *AckHandler) info(str string) {
	logrus.
		WithField("system", "falcon").
		WithField("struct", "AckHandler").Info(str)
}

//输出日志
func (me *AckHandler) warn(str string) {
	logrus.
		WithField("system", "falcon").
		WithField("struct", "AckHandler").Warn(str)
}
