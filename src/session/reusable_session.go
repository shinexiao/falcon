package session

import (
	"connector/connection"
	"bytes"
	"strings"
	"common"
	"security"
)

/**
 * 可重用session
 */
type ReusableSession struct {
	SessionId  string;
	ExpireTime int64
	Context    *connection.SessionContext
}

func NewReusableSession() *ReusableSession {

	return &ReusableSession{}
}

//编码ReusableSession
func Encode(context *connection.SessionContext) string {

	buf := new(bytes.Buffer)
	buf.WriteString(context.OsName)
	buf.WriteString(",")
	buf.WriteString(context.OsVersion)
	buf.WriteString(",")
	buf.WriteString(context.ClientVersion)
	buf.WriteString(",")
	buf.WriteString(context.DeviceId)
	buf.WriteString(",")
	buf.WriteString(context.Cipher.ToString())

	return buf.String()
}

//解码ReusableSession
func Decode(value string) *ReusableSession {

	array := strings.Split(value, ",")
	if len(array) != 6 {
		return nil
	}

	context := connection.NewSessionContext()
	context.OsName = array[0]
	context.OsVersion = array[1]
	context.ClientVersion = array[2]
	context.DeviceId = array[3]
	//加密密钥,加密方式
	key := common.HexToBye(array[4])
	iv := common.HexToBye(array[5])

	context.Cipher = security.NewAesCipher(key, iv)
	//恢复session
	session := NewReusableSession()
	session.Context = context

	return session
}
