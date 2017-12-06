package security

type Cipher interface {
	//解密
	Decrypt(data []byte) []byte
	//加密
	Encrypt(data []byte) []byte
	//转换为String
	ToString() string
}
