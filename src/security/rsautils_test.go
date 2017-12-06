package security

import (
	"testing"
	"fmt"
)

func TestParsePrivateKey(t *testing.T) {

	privateKey := `-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQCrfrw/ErEziTHilOhFyd98+J6suaKPB5kAFh+jVEKUM86d+Wx1
r1e3RmvzImrNTpZ7LQtrT9GRX8vwyacsvRu/xBqq30PgS8/aikSdIiMXU1iOQErv
VHOumIPi0g4hKlglQKTQ9U1E4UEoMokHWNUtbYimwDpmUQDk9NwCUf7iDQIDAQAB
AoGANChLYHNy6VWkkmDvc6o+Cmgi+i1LP2z0H46a+LW7ug83m9wsHG7Dor4MPtoM
2Xw5UCUXAAA6oJgeEpGCAp1RPrYUJOODs6T+DGzrp0wM6wCx9cozh/a69UNx5JAQ
pwJP9bwGndM7QXk6GrjcEJ7dZWGLtAyPx6eNPySj3ErmnpECQQDT2Q0jdet+UgMj
AV6iN7djnD+oKQ33NOHxX3qUNNx0Tcpr/HpxxpcKBwOEL8CTPlHOLXmzMZGRyLLC
6KeVYaD7AkEAzzyyi3QPttBHkUlAjvgt4HcFo+FTcDdmzCSPrBVJMalBjfYllNNS
Lp7w6J+Jug4yt1WOP5bhFeCEKjSw6HNqlwJAVtgZrLnAai5Qnt8G3lUc1rbM2bDK
ytZg8UQEyhDJdtwU6SO9Rjr02+V4KY4x0aqwembmBvGBDVRLA9/AI1q8VQJAQFCT
DKppUhATlehI69XjzvzBOFnuni3jbkmOeRZmD856dMdGZIiswaE8HMWeZaqQXMtl
iSCXHEYAXmTZ3lorYwJBAL4tFju9cX7uKzEPGg9CpvbDJmBnwZo1qp1WI495mKh2
jpZVEn4Y8liRAR95fQt6FNn1xDm2AKoKH/3WlItA06c=
-----END RSA PRIVATE KEY-----`

	key, err := ParsePrivateKey([]byte(privateKey))

	fmt.Println(key, err)

}
