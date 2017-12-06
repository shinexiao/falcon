package security

import (
	"encoding/pem"
	"crypto/rsa"
	"crypto/x509"
	"crypto/rand"
	"errors"
	"crypto"
	"bytes"
)

// 生成RSA密钥
// @param	bits	密钥位数
// @return 	prikBlock	pem.Block
// @return 	pubkBlock	pem.Block
// @return 	prik	[]byte
// @return 	pubk	[]byte
// @return 	err		error
func GenerateRSASecretKey(bits int) (prikBlock *pem.Block, pubkBlock *pem.Block, prik []byte, pubk []byte, err error) {
	if bits < 1024 {
		bits = 1024
	}

	/*-------------------生成私钥------------------*/
	prikey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	prikByte := x509.MarshalPKCS1PrivateKey(prikey)

	/*-------------------得到公钥------------------*/
	pubkey := &prikey.PublicKey
	pubkByte, err := x509.MarshalPKIXPublicKey(pubkey)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	prikBlock = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: prikByte,
	}

	pubkBlock = &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubkByte,
	}

	prik = pem.EncodeToMemory(prikBlock)
	pubk = pem.EncodeToMemory(pubkBlock)
	return prikBlock, pubkBlock, prik, pubk, nil
}

// 解析公钥
func ParsePublicKey(publicKey []byte) (pubk *rsa.PublicKey, err error) {

	if nil == publicKey {
		return nil, errors.New("rsa public key nil!")
	}

	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("rsa public key error")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return pubInterface.(*rsa.PublicKey), nil
}

// 解析私钥
func ParsePrivateKey(privateKey []byte) (prik *rsa.PrivateKey, err error) {

	if nil == privateKey {
		return nil, errors.New("rsa private key nil!")
	}

	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("rsa private key error!")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

// 公钥加密
func RSAEncrypt(encryptText, pubk []byte) (encrypt []byte, err error) {

	pubkey, err := ParsePublicKey(pubk)
	if err != nil {
		return nil, err
	}

	return rsa.EncryptPKCS1v15(rand.Reader, pubkey, encryptText)
}

// 私钥解密
func RSADecrypt(decryptText []byte, key *rsa.PrivateKey) (decrypt []byte, err error) {

	//prikey, err := ParsePrivateKey(prik)
	//if err != nil {
	//	return nil, err
	//}

	return rsa.DecryptPKCS1v15(rand.Reader, key, decryptText)
}

// 私钥签名
func RSASign(data, prik []byte, hash crypto.Hash) (sign []byte, err error) {

	prikey, err := ParsePrivateKey(prik)
	if err != nil {
		return nil, err
	}

	var hashed []byte = nil
	switch hash {
	case crypto.MD5:
		hashed = MD5For32(data, false)
	case crypto.SHA1:
		hashed = SHAX(data, crypto.SHA1, false)
	case crypto.SHA256:
		hashed = SHAX(data, crypto.SHA256, false)
	case crypto.SHA512:
		hashed = SHAX(data, crypto.SHA512, false)
	default:
		return nil, errors.New("the hash is not support")
	}

	if hashed == nil {
		return nil, errors.New("hashed error occurred")
	}

	return rsa.SignPKCS1v15(rand.Reader, prikey, hash, hashed)
}

// 公钥验签
func RSAVerify(data, sign, pubk []byte, hash crypto.Hash) (err error) {

	pub, err := ParsePublicKey(pubk)
	if err != nil {
		return err
	}

	var hashed []byte = nil
	switch hash {
	case crypto.MD5:
		hashed = MD5For32(data, false)
	case crypto.SHA1:
		hashed = SHAX(data, crypto.SHA1, false)
	case crypto.SHA256:
		hashed = SHAX(data, crypto.SHA256, false)
	case crypto.SHA512:
		hashed = SHAX(data, crypto.SHA512, false)
	default:
		return errors.New("the hash is not support")
	}

	if hashed == nil {
		return errors.New("hashed error occurred")
	}

	return rsa.VerifyPKCS1v15(pub, hash, hashed, sign)
}

// 格式话密钥文件，主要是读取的时候不换行的话，会在pem.Decode()的时候出现密钥不通过
// @param	header	头部内容
// @param	src		读取的原来的密钥
// @return 	cert	格式化换行后的密钥
func RsaCertFormat(header, src string) (cert string) {
	bits := 64
	offset := 0
	src_len := len(src)

	buffer := new(bytes.Buffer)
	buffer.Write([]byte("-----BEGIN " + header + "-----\n"))
	for i := 0; src_len-offset > 0; offset = i * bits {
		if src_len-offset > bits {
			buffer.Write([]byte(src[offset:offset+bits] + "\n"))
		} else {
			buffer.Write([]byte(src[offset:src_len] + "\n"))
		}
		i ++
	}
	buffer.Write([]byte("-----END " + header + "-----\n"))
	return buffer.String()
}

// 为密钥文件内容添加头尾
// @param 	header	头部内容，不要包括'BEGIN''END'
// @param 	src		原密钥内容
// @return 	cert	添加后的密钥
func RsaCertAddHeader(header, src string) (cert string) {
	buffer := new(bytes.Buffer)
	buffer.Write([]byte("-----BEGIN " + header + "-----\n"))
	buffer.Write([]byte(src))
	buffer.Write([]byte("\n-----END " + header + "-----\n"))
	return buffer.String()
}

// 公钥分段加密
func RSASubEncrypt(encryptText, pubk []byte, bits int) (encrypt []byte, err error) {

	// 1024位的证书，加密时最大支持117个字节，解密时为128
	// 2048位的证书，加密时最大支持245个字节，解密时为256。
	// 加密时支持的最大字节数：证书位数/8 -11（比如：2048位的证书，支持的最大加密字节数：2048/8 - 11 = 245）
	// 解密时支持的最大字节数：证书位数/8（比如：2048位的证书，支持的最大加密字节数：2048/8  = 256）

	// 解析证书最大加密长度
	maxEnLen := bits/8 - 11
	// 获取待加密数据直接长度
	textLen := len(encryptText)
	// 获取分段的最大次数
	maxCount := 0
	if textLen%maxEnLen > 0 {
		maxCount = textLen/maxEnLen + 1
	} else {
		maxCount = textLen / maxEnLen
	}
	// 缓存每次加密的数据
	buffer := &bytes.Buffer{}

	// 开始分段加密
	for i := 0; i < maxCount; i ++ {

		var encrypt []byte = nil

		// 如果当前的加密下标小于文本长度,则表示可以正常截取,否则就只能截取当前的文本长度
		if textLen-maxEnLen > 0 {

			// 获取本次要加密的文本
			countText := []byte(string(encryptText)[:maxEnLen])
			// 加密该段文本
			encrypt, err = RSAEncrypt(countText, pubk)
			if err != nil {
				return nil, err
			}
			// 将加密的文本剔除,留下未加密的文本
			encryptText = []byte(string(encryptText)[maxEnLen:])
		} else {

			encrypt, err = RSAEncrypt(encryptText, pubk)
			if err != nil {
				return nil, err
			}
		}

		// 重新定义offet,文本长度
		textLen = len(encryptText)
		buffer.Write(encrypt)
	}

	return buffer.Bytes(), nil
}

// 私钥分段解密
func RSASubDecrypt(decryptText []byte, prik *rsa.PrivateKey, bits int) (decrypt []byte, err error) {

	// 1024位的证书，加密时最大支持117个字节，解密时为128
	// 2048位的证书，加密时最大支持245个字节，解密时为256。
	// 加密时支持的最大字节数：证书位数/8 -11（比如：2048位的证书，支持的最大加密字节数：2048/8 - 11 = 245）
	// 解密时支持的最大字节数：证书位数/8（比如：2048位的证书，支持的最大加密字节数：2048/8  = 256）

	// 解析证书最大解密长度
	maxDnLen := bits / 8
	// 获取待解密数据直接长度
	textLen := len(decryptText)
	// 获取分段的最大次数
	maxCount := 0
	if textLen%maxDnLen > 0 {
		maxCount = textLen/maxDnLen + 1
	} else {
		maxCount = textLen / maxDnLen
	}
	// 缓存每次解密的数据
	buffer := &bytes.Buffer{}

	// 开始分段解密
	for i := 0; i < maxCount; i ++ {

		var decrypt []byte = nil

		// 如果当前的解密下标小于文本长度,则表示可以正常截取,否则就只能截取当前的文本长度
		if textLen-maxDnLen > 0 {

			// 获取本次要解密的文本
			countText := []byte(string(decryptText)[:maxDnLen])
			// 将解密的文本剔除,留下未解密的文本
			decryptText = []byte(string(decryptText)[maxDnLen:])
			// 解密该段文本
			decrypt, err = RSADecrypt(countText, prik)
			if err != nil {
				return nil, err
			}
		} else {

			decrypt, err = RSADecrypt([]byte(string(decryptText)), prik)
			if err != nil {
				return nil, err
			}
		}

		// 重新定义offet,文本长度
		textLen = len(decryptText)
		buffer.Write(decrypt)
	}

	return buffer.Bytes(), nil
}
