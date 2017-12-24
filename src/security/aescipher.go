package security

import (
	"bytes"
	"falcon/src/common"
)

type AesCipher struct {
	key []byte
	iv  []byte
}

/**
 * 创建RSA加解密
 */
func NewAesCipher(key, iv []byte) *AesCipher {

	return &AesCipher{
		key: key,
		iv:  iv,
	}
}

//解密
func (me *AesCipher) Decrypt(data []byte) []byte {

	if decrypted, err := AesDecrypt(data, me.key, me.iv); err != nil {
		return nil
	} else {
		return decrypted
	}
}

//加密
func (me *AesCipher) Encrypt(data []byte) []byte {

	if encrypted, err := AesEncrypt(data, me.key, me.iv); err != nil {
		return nil
	} else {
		return encrypted
	}

}

//转换为string
func (me *AesCipher) ToString() string {

	buf := new(bytes.Buffer)
	buf.WriteString(common.ByteToHex(me.key))
	buf.WriteString(",")
	buf.WriteString(common.ByteToHex(me.iv))

	return buf.String()
}
