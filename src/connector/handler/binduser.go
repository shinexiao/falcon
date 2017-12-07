package handler

import (
	"connector/protocol"
	"connector/connection"
	"github.com/gogap/logrus"
	"connector/message"
	"fmt"
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
		me.info("recive bind user message.")
		me.bind(msg)
	} else {
		me.info("recive unbind user message.")
		me.unbind(msg)
	}

}

//解码出消息
func (me *BindUserHandler) decode(packet *protocol.Packet, conn *connection.Conn) *message.BindUserMessage {

	msg := message.NewBindUserMessage(packet, conn)
	msg.DecodeBody()

	return msg
}

//绑定消息
func (me *BindUserHandler) bind(bindUser *message.BindUserMessage) {

	if len(bindUser.UserId) == 0 {
		errMsg := message.ErrorMessageFrom(bindUser)
		errMsg.Reson = "invalid param"
		errMsg.Close()
		me.error(fmt.Sprintf("bind user failure for invalid param, conn=%v", bindUser.Conn()))
		return
	}

	//检查是否握手成功
	ctx := bindUser.Conn().Context
	if ctx.HandshakeOk() {
		//处理重复绑定
		if len(ctx.UserId) > 0 {
			if bindUser.UserId == ctx.UserId {
				ctx.Tags = bindUser.Tags
				okMsg := message.OkMessageFrom(bindUser)
				okMsg.Data = "bind success"
				okMsg.Send() //原来是sendRaw
				me.info(fmt.Sprintf("rebind user success, userId=%s, session=%s", bindUser.UserId, ctx))
			} else {
				me.unbind(bindUser)
			}
		}

		//todo 验证用户身份
		//todo 验证成功注册到路由上

		var success = true

		if success {
			ctx.UserId = bindUser.UserId
			ctx.Tags = bindUser.Tags
			//发送注册成功事件
			okMsg := message.OkMessageFrom(bindUser)
			okMsg.Data = "bind success"
			okMsg.Send() //原来是sendRaw
			me.info(fmt.Sprintf("bind user success, userId=%s, session=%s", bindUser.UserId, ctx))
		} else {
			//注册失败
			errMsg := message.ErrorMessageFrom(bindUser)
			errMsg.Reson = "bind faild"
			errMsg.Close()
			me.error(fmt.Sprintf("bind user failure, userId=%s, session=%v", bindUser.UserId, ctx))
		}

	} else {
		errMsg := message.ErrorMessageFrom(bindUser)
		errMsg.Reson = "no handshake"
		errMsg.Close()
		me.error(fmt.Sprintf("bind user failure not handshake, userId=%s, conn=%v", bindUser.UserId, bindUser.Conn()))
	}

}

//解绑消息
func (me *BindUserHandler) unbind(bindUser *message.BindUserMessage) {

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
