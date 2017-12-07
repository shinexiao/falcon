package message

import (
	"connector/protocol"
	"common"

	"sync/atomic"
	"connector/connection"
)

const (
	STATUS_DECODED byte = 1
	STATUS_ENCODED byte = 2
)

const CompressLimit = 1024

//全局sid增长序列
var SID_SEQ int32

type Message interface {
	Conn() *connection.Conn //获取连接
	DecodeBody()            //解码消息体
	EncodeBody()            //编码消息体
	Send()                  //发送
	//	SendRaw()                 //发送原始消息
	Packet() *protocol.Packet //数据包
}

//内部消息
type InternalMessage interface {
	Decode(body []byte) //解码
	Encode() []byte     //编码
}

/**
 * 基础消息
 */
type BaseMessage struct {
	packet *protocol.Packet
	conn   *connection.Conn
	Status byte
	Child  InternalMessage
}

//获取连接
func (me *BaseMessage) Conn() *connection.Conn {
	return me.conn
}

func (me *BaseMessage) Packet() *protocol.Packet {

	return me.packet
}

//decode 消息体
func (me *BaseMessage) DecodeBody() {
	if (me.Status & STATUS_DECODED) != 0 {
		return
	} else {
		me.Status |= STATUS_DECODED
	}

	if me.packet.Body != nil && len(me.packet.Body) > 0 {

		tmp := me.packet.Body

		ctx := me.conn.Context
		//1、解密
		if me.packet.HasFlag(protocol.FLAG_CRYPTO) {
			if ctx.Cipher != nil {
				tmp = ctx.Cipher.Decrypt(tmp)
			}
		}
		//2、解压
		if me.packet.HasFlag(protocol.FLAG_COMPRESS) {
			tmp = common.UnCompress(tmp)
		}

		if len(tmp) == 0 {
			panic("message decode ex");
		}

		me.packet.Body = tmp

		//调用子类
		me.Child.Decode(me.packet.Body)
	}

}

//encode 消息体
func (me *BaseMessage) EncodeBody() {
	if (me.Status & STATUS_ENCODED) != 0 {
		return
	} else {
		me.Status |= STATUS_ENCODED
	}
	tmp := me.Child.Encode()

	if tmp != nil && len(tmp) > 0 {
		//1.压缩
		if len(tmp) > CompressLimit {
			result := common.Compress(tmp)
			if len(result) > 0 {
				tmp = result
				me.packet.AddFlag(protocol.FLAG_COMPRESS)
			}
		}

		//2.加密
		ctx := me.conn.Context
		if ctx.Cipher != nil {
			result := ctx.Cipher.Encrypt(tmp)
			if len(result) > 0 {
				tmp = result
				me.packet.AddFlag(protocol.FLAG_CRYPTO)
			}
		}

		me.packet.Body = tmp
	}
}

//生成sessionId
func (me *BaseMessage) GenSessionId() int32 {

	SID_SEQ = atomic.AddInt32(&SID_SEQ, 1)

	return SID_SEQ
}

//发送消息
func (me *BaseMessage) Send() {
	me.EncodeBody()
	me.conn.Send(me.packet)
}

//发送并关闭通道
func (me *BaseMessage) Close() {
	me.EncodeBody()
	me.conn.SendWithListener(me.packet, func(conn *connection.Conn, err error) {
		conn.Close()
	})
}
