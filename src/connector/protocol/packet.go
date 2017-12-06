package protocol

import "encoding/json"

const (
	PACKET_HEADER_LEN int  = 13 //头长度
	FLAG_CRYPTO       int8 = 1  //加密标识
	FLAG_COMPRESS     int8 = 2  //压缩标识
	FLAG_BIZ_ACK      int8 = 4  //ack确认标识
	FLAG_AUTO_ACK     int8 = 8  //自动ack标识
	FLAG_JSON_BODY    int8 = 16 //json消息体

	HB_PACKET_BYTE int8 = -33 //心跳消息
)

//通信协议
//length(4)+cmd(1)+cc(2)+flags(1)+sessionId(4)+lrc(1)+body(n)
type Packet struct {
	Cmd       int8   //命令
	CC        int16  //校验码 暂时未用
	Flags     int8   //特性，如是否加密，是否压缩等
	SessionId int32  // 会话id。客户端生成。
	Lrc       int8   // 校验，纵向冗余校验。只校验head
	Body      []byte //消息体
}

//创建消息
func NewPacket(cmd int8) *Packet {

	return &Packet{
		Cmd: cmd,
	}
}

//获取消息体长度
func (me *Packet) BodyLength() int {

	if nil != me.Body {
		return len(me.Body)
	}

	return 0
}

//添加特性标记
func (me *Packet) AddFlag(flag int8) {
	me.Flags |= flag
}

//检测是否存在指定特性标记
func (me *Packet) HasFlag(flag int8) bool {

	return (me.Flags & flag) != 0
}

//方便调试打印
func (me *Packet) ToString() string {

	data, _ := json.Marshal(me)

	return string(data)
}
