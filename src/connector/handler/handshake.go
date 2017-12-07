package handler

import (
	"connector/protocol"
	"connector/message"
	"connector/connection"
	"security"
	"session"
	"github.com/gogap/logrus"
	"core"
)

/**
 * 握手认证
 */
type HandshakeHandler struct {
	reusableSessionManager *session.ReusableSessionManager
}

func NewHandshakeHandler(context core.FalconContext) *HandshakeHandler {

	return &HandshakeHandler{
		reusableSessionManager: context.ReusableSessionManager(),
	}
}

func (me *HandshakeHandler) Handle(packet *protocol.Packet, conn *connection.Conn) {

	msg := me.decode(packet, conn)

	if msg.Conn().Context.IsSecurity() {
		me.doSecurity(msg)
	} else {
		me.doInsecurity(msg)
	}

}

//解码出消息
func (me *HandshakeHandler) decode(packet *protocol.Packet, conn *connection.Conn) *message.HandshakeMessage {

	msg := message.NewHandshakeMessagePacket(packet, conn)
	msg.DecodeBody()

	return msg
}

//安全方式认证
func (me *HandshakeHandler) doSecurity(msg *message.HandshakeMessage) {
	iv := msg.Iv                                        //AES密钥向量16位
	clientKey := msg.ClientKey                          //客户端随机数16位
	serverKey := security.RandomAESKey()                //服务端随机数16位
	sessionKey := security.MixKey(clientKey, serverKey) //会话密钥16位
	ctx := msg.Conn().Context

	//1.校验客户端消息字段
	if len(msg.DeviceId) == 0 || len(iv) != security.AES_KEY_LENGTH || len(clientKey) != security.AES_KEY_LENGTH {
		errMsg := message.ErrorMessageFrom(msg)
		errMsg.Reson = "Param invalid"
		errMsg.Close()
		me.error("handshake in security Param invalid.")
		return
	}

	//2.重复握手判断
	if msg.DeviceId == ctx.DeviceId {
		errMsg := message.ErrorMessageFrom(msg)
		errMsg.Reson = "repeat handshake"
		errMsg.Send()
		me.error("handshake in security repeat handshake.")
		return
	}

	//3.更换会话密钥RSA=>AES(clientKey)
	ctx.Cipher = security.NewAesCipher(clientKey, iv)

	//4.生成可复用session, 用于快速重连
	sess := me.reusableSessionManager.GenSession(ctx)

	//5.计算心跳时间
	heartbeat := msg.MaxHeartbeat //todo 这里需要与服务器配置进行对比

	//6.响应握手成功消息
	okMsg := message.NewHandshakeOkMessage(msg.Packet(), msg.Conn())
	okMsg.ServerKey = serverKey
	okMsg.Heartbeat = heartbeat
	okMsg.SessionId = sess.SessionId
	okMsg.ExpireTime = sess.ExpireTime //传输给客户端的需要乘以1000,服务器为秒，客户端为毫秒
	okMsg.Send()

	//7.更换会话密钥AES(clientKey)=>AES(sessionKey)
	ctx.Cipher = security.NewAesCipher(sessionKey, iv)

	//8.保存client信息到当前连接
	ctx.OsName = msg.OsName
	ctx.OsVersion = msg.OsVersion
	ctx.ClientVersion = msg.ClientVersion
	ctx.DeviceId = msg.DeviceId
	ctx.Heartbeat = heartbeat

	//9.保存可复用session到Redis, 用于快速重连
	me.reusableSessionManager.CacheSession(sess)

	me.info("handshake in security success.")
}

//非安全方式认证
func (me *HandshakeHandler) doInsecurity(msg *message.HandshakeMessage) {
	//1.校验客户端消息字段
	ctx := msg.Conn().Context

	if len(msg.DeviceId) == 0 {
		errMsg := message.ErrorMessageFrom(msg)
		errMsg.Reson = "Param invalid"
		errMsg.Close()
		me.error("handshake not security Param invalid.")
		return
	}
	//2.重复握手判断
	if msg.DeviceId == ctx.DeviceId {
		errMsg := message.ErrorMessageFrom(msg)
		errMsg.Reson = "repeat handshake"
		errMsg.Send()
		me.error("handshake not security repeat handshake.")
		return
	}

	//6.响应握手成功消息
	message.NewHandshakeOkMessage(msg.Packet(), msg.Conn()).Send()

	//8.保存client信息到当前连接
	ctx.OsName = msg.OsName
	ctx.OsVersion = msg.OsVersion
	ctx.ClientVersion = msg.ClientVersion
	ctx.DeviceId = msg.DeviceId
	ctx.Heartbeat = 0x7fffffff

	me.info("handshake not security success.")
}

//输出日志
func (me *HandshakeHandler) info(str string) {
	logrus.
		WithField("system", "falcon").
		WithField("struct", "HandshakeHandler").Info(str)
}

//输出日志
func (me *HandshakeHandler) error(str string) {
	logrus.
		WithField("system", "falcon").
		WithField("struct", "HandshakeHandler").Error(str)
}
