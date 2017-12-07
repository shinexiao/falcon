package message

import (
	"bytes"
	"encoding/binary"
	"github.com/pkg/errors"
)

const INT16_MAX_VALUE = 32767

/**
 * 消息编码
 */
type ByteBufMessage struct {
	BaseMessage
}

//编码字符串
func (me *ByteBufMessage) EncodeString(buf *bytes.Buffer, filed string) {
	me.EncodeBytes(buf, []byte(filed))
}

//编码byte
func (me *ByteBufMessage) EncodeBytes(buf *bytes.Buffer, field []byte) {
	var length int16
	dataLength := len(field)
	var err error
	if nil == field || dataLength == 0 {
		length = 0
		err = binary.Write(buf, binary.LittleEndian, length)
		checkError(err, "EncodeBytes write int16")
	} else if dataLength < INT16_MAX_VALUE {
		length = int16(dataLength)
		err = binary.Write(buf, binary.BigEndian, length)
		checkError(err, "EncodeBytes write int16")
		err = binary.Write(buf, binary.BigEndian, field)
		checkError(err, "EncodeBytes write bytes")
	} else {
		err = binary.Write(buf, binary.BigEndian, int16(INT16_MAX_VALUE))
		checkError(err, "EncodeBytes write int16")
		err = binary.Write(buf, binary.BigEndian, int32(dataLength-INT16_MAX_VALUE))
		checkError(err, "EncodeBytes write int32")
		err = binary.Write(buf, binary.BigEndian, field)
		checkError(err, "EncodeBytes write bytes")
	}
}

//编码int8
func (me *ByteBufMessage) EncodeInt8(buf *bytes.Buffer, field int8) () {
	err := binary.Write(buf, binary.BigEndian, field)
	checkError(err, "EncodeInt8 write int8")
}

//编码int16
func (me *ByteBufMessage) EncodeInt16(buf *bytes.Buffer, field int16) () {
	err := binary.Write(buf, binary.BigEndian, field)
	checkError(err, "EncodeInt16 write int16")
}

//编码int32
func (me *ByteBufMessage) EncodeInt32(buf *bytes.Buffer, field int32) () {
	err := binary.Write(buf, binary.BigEndian, field)
	checkError(err, "EncodeInt32 write int32")
}

//编码int64
func (me *ByteBufMessage) EncodeInt64(buf *bytes.Buffer, field int64) () {
	err := binary.Write(buf, binary.BigEndian, field)
	checkError(err, "EncodeInt64 write int64")
}

//解码
func (me *ByteBufMessage) DecodeString(buf *bytes.Buffer) string {
	data := me.DecodeBytes(buf)
	if data == nil {
		return ""
	}
	return string(data)
}

//解码byte
func (me *ByteBufMessage) DecodeBytes(buf *bytes.Buffer) []byte {
	var length int16
	var total int
	var err error
	err = binary.Read(buf, binary.BigEndian, &length)
	checkError(err, "DecodeBytes read int16")
	if length == 0 {
		return nil
	}

	total = int(length)

	if total == INT16_MAX_VALUE {
		var length2 int32
		err = binary.Read(buf, binary.BigEndian, &length2)
		checkError(err, "DecodeBytes read int32")
		total = total + int(length2)
	}

	data := make([]byte, total)
	err = binary.Read(buf, binary.BigEndian, data)
	checkError(err, "DecodeBytes read bytes")

	return data
}

//解码int8
func (me *ByteBufMessage) DecodeInt8(buf *bytes.Buffer) int8 {
	var value int8
	err := binary.Read(buf, binary.BigEndian, &value)
	checkError(err, "DecodeInt8 read int8")
	return value
}

//解码int16
func (me *ByteBufMessage) DecodeInt16(buf *bytes.Buffer) int16 {
	var value int16
	err := binary.Read(buf, binary.BigEndian, &value)
	checkError(err, "DecodeInt16 read int16")
	return value
}

//解码int32
func (me *ByteBufMessage) DecodeInt32(buf *bytes.Buffer) int32 {
	var value int32
	err := binary.Read(buf, binary.BigEndian, &value)
	checkError(err, "DecodeInt32 read int32")
	return value
}

//解码int64
func (me *ByteBufMessage) DecodeInt64(buf *bytes.Buffer) int64 {
	var value int64
	err := binary.Read(buf, binary.BigEndian, &value)
	checkError(err, "DecodeInt64 read int64")
	return value
}

//错误终止
func checkError(err error, str string) {
	if err != nil {
		panic(errors.New("ByteBufMessage error " + err.Error() + " " + str))
	}
}
