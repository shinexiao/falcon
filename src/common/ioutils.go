package common

import (
	"compress/zlib"
	"bytes"
	"io"
)

// zip 压缩数据
func Compress(data []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write(data)
	w.Close()
	return in.Bytes()
}

// zip 解压数据
func UnCompress(data []byte) []byte {
	b := bytes.NewReader(data)
	var out bytes.Buffer
	r, _ := zlib.NewReader(b)
	io.Copy(&out, r)
	return out.Bytes()
}
