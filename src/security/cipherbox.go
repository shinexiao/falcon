package security

import (
	"math/rand"
	"math"
)

const AES_KEY_LENGTH = 16

func RandomAESKey() []byte {

	data := make([]byte, AES_KEY_LENGTH)

	if _, err := rand.Read(data); err != nil {
		panic("rand.Read error")
	}
	return data
}

func RandomAESIV() []byte {

	data := make([]byte, AES_KEY_LENGTH)

	if _, err := rand.Read(data); err != nil {
		panic("rand.Read error")
	}

	return data
}

func MixKey(clientKey []byte, serverKey []byte) []byte {

	sessionKey := make([]byte, AES_KEY_LENGTH)

	for i := 0; i < AES_KEY_LENGTH; i++ {
		a := clientKey[i]
		b := serverKey[i]

		sum := math.Abs(float64(a + b));
		var c byte

		if int(sum)%2 == 0 {
			c = a ^ b
		} else {
			c = b ^ a
		}
		sessionKey[i] = c
	}

	return sessionKey;
}
