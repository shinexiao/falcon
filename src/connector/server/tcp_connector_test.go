package server

import (
	"testing"
	"net"
)

func TestBind(t *testing.T) {

	conn, err := net.Dial("tcp", "127.0.0.1:3000")
	assert.Nil(t, err)

	conn.Close()

}
