package protocol

import (
	"bytes"
	"encoding/binary"
	"github.com/pkg/errors"
)

//编码消息体
func EncodePacket(packet *Packet, out *bytes.Buffer) {
	var err error
	if packet.Cmd == HB_PACKET_BYTE {
		binary.Write(out, binary.BigEndian, HB_PACKET_BYTE)
	} else {
		err = binary.Write(out, binary.BigEndian, int32(packet.BodyLength()))
		checkError(err)
		err = binary.Write(out, binary.BigEndian, packet.Cmd)
		checkError(err)
		err = binary.Write(out, binary.BigEndian, packet.CC)
		checkError(err)
		err = binary.Write(out, binary.BigEndian, packet.Flags)
		checkError(err)
		err = binary.Write(out, binary.BigEndian, packet.SessionId)
		checkError(err)
		err = binary.Write(out, binary.BigEndian, packet.Lrc)
		checkError(err)
		if packet.BodyLength() > 0 {
			err = binary.Write(out, binary.BigEndian, packet.Body)
			checkError(err)
		}
	}
}

//错误终止
func checkError(err error) {
	if err != nil {
		panic(errors.New("encoder error " + err.Error()))
	}
}

//解码消息体
func DecodePacket(packet *Packet, in *bytes.Buffer, bodyLength int) {
	var err error
	err = binary.Read(in, binary.BigEndian, &packet.CC)
	checkError(err)
	err = binary.Read(in, binary.BigEndian, &packet.Flags)
	checkError(err)
	err = binary.Read(in, binary.BigEndian, &packet.SessionId)
	checkError(err)
	err = binary.Read(in, binary.BigEndian, &packet.Lrc)
	checkError(err)
	if bodyLength > 0 {
		packet.Body = make([]byte, bodyLength)
		err = binary.Read(in, binary.BigEndian, packet.Body)
		checkError(err)
	}
}
