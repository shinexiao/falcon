package handler

import (
	"connector/protocol"
	"connector/connection"
	"github.com/gogap/logrus"
	"connector/message"
	"core"
	"fmt"
)

type FastConnectHandler struct {
	ctx core.FalconContext
}

func NewFastConnectHandler(ctx core.FalconContext) *FastConnectHandler {

	return &FastConnectHandler{
		ctx: ctx,
	}
}

func (me *FastConnectHandler) Handle(packet *protocol.Packet, conn *connection.Conn) {

	msg := me.decode(packet, conn)

	sess := me.ctx.ReusableSessionManager().QuerySession(msg.SessionId)

	if sess == nil {
		//没有查询到session说明session已经失效
		errMsg := message.ErrorMessageFrom(msg)
		errMsg.Reson = "session expired"
		errMsg.Send()
		me.warn(fmt.Sprintf("fast connect failure, session is expired, sessionId=%s, deviceId=%s, conn=%v",
			msg.SessionId, msg.DeviceId, conn))
	} else if sess.Context.DeviceId != msg.DeviceId {
		//非法的设备，当前设备不是上次生成session时的设备
		errMsg := message.ErrorMessageFrom(msg)
		errMsg.Reson = "invalid device"
		errMsg.Send()
		me.warn(fmt.Sprintf("fast connect failure, not the same device, deviceId=%s, session=%v, conn=%v",
			msg.DeviceId, sess.Context, conn))
	} else {
		//校验成功，重新计算心跳，完成快速重连
		heartbeat := msg.MaxHeartbeat
		sess.Context.Heartbeat = heartbeat
		conn.Context = sess.Context

		//发送消息
		okMsg := message.FastConnectOkMessageFrom(msg)
		okMsg.Heartbeat = heartbeat
		okMsg.Send()
	}

}

//解码出消息
func (me *FastConnectHandler) decode(packet *protocol.Packet, conn *connection.Conn) *message.FastConnectMessage {

	msg := message.NewFastConnectMessage(packet,conn)
	msg.DecodeBody()

	return msg
}

//输出日志
func (me *FastConnectHandler) debug(str string) {
	logrus.
		WithField("system", "falcon").
		WithField("struct", "FastConnectHandler").Debug(str)
}

//输出日志
func (me *FastConnectHandler) info(str string) {
	logrus.
		WithField("system", "falcon").
		WithField("struct", "FastConnectHandler").Info(str)
}

//输出日志
func (me *FastConnectHandler) warn(str string) {
	logrus.
		WithField("system", "falcon").
		WithField("struct", "FastConnectHandler").Warn(str)
}
