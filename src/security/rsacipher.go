package security

import (
	"crypto/rsa"
)

type RsaCipher struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

/**
 * 创建RSA加解密
 */
func NewRsaCipher(publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) *RsaCipher {

	return &RsaCipher{
		publicKey:  publicKey,
		privateKey: privateKey,
	}
}

//解密
func (me *RsaCipher) Decrypt(data []byte) []byte {

	if data, err := RSADecrypt(data, me.privateKey); err != nil {
		panic(err)
	} else {
		return data
	}

}

//加密
func (me *RsaCipher) Encrypt(data []byte) []byte {
	//Todo
	return nil
}

//转换为string
func (me *RsaCipher) ToString() string {
	//Todo
	return ""
}
