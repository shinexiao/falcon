package security

import (
	"crypto/sha1"
	"encoding/hex"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"crypto"
)

// type
const (
	SHA_1 = "SHA1"
	SHA_256 = "SHA_256"
	SHA_512 = "SHA_512"
)

// SHA加密[1, 256, 512]
func SHAX(data []byte, shaType crypto.Hash, isHex bool) (sh []byte) {

	var hashs hash.Hash = nil
	switch shaType {
	case crypto.SHA1:
		hashs = sha1.New()
	case crypto.SHA256:
		hashs = sha256.New()
	case crypto.SHA3_512:
		hashs = sha512.New()
	default:
		return nil
	}

	hashs.Write(data)
	hashed := hashs.Sum(nil)

	if isHex {
		return []byte(hex.EncodeToString(hashed))
	}
	return hashed
}
