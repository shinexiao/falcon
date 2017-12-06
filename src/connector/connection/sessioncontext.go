package connection

import (
	"common"
	"security"
)

type SessionContext struct {
	OsName        string
	OsVersion     string
	ClientVersion string
	DeviceId      string
	UserId        string
	Tags          string
	Heartbeat     int32
	Cipher        security.Cipher
	ClientType    int8
}

/**
 * 创建Session上下文
 */
func NewSessionContext() *SessionContext {

	return &SessionContext{
		Heartbeat: 10000,
	}
}

/**
 * 判断是否认证成功
 */
func (me *SessionContext) HandshakeOk() bool {

	return len(me.DeviceId) > 0
}

/**
 * 获取客户端类型
 */
func (me *SessionContext) GetClientType() int8 {

	if me.ClientType == 0 {
		me.ClientType = common.FindClientType(me.OsName)
	}

	return me.ClientType
}

/**
 * 判断是否使用安全模式
 */
func (me *SessionContext) IsSecurity() bool {

	return me.Cipher != nil
}
