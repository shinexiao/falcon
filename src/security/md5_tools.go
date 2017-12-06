package security

import (
	"crypto/md5"
	"encoding/hex"
)

// MD532位加密
func MD5For32(data []byte, isHex bool) (md []byte) {
	hash := md5.New()
	hash.Write(data)
	hashed := hash.Sum(nil)

	if isHex {
		return []byte(hex.EncodeToString(hashed))
	}
	return hashed
}

// MD516位加密
func MD5For16(data []byte, isHex bool)  (md []byte) {
	return MD5For32(data, isHex)[8:24]
}

