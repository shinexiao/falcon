package server

import (
	"net"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBind(t *testing.T) {

	conn, err := net.Dial("tcp", "127.0.0.1:3000")
	assert.Nil(t, err)

	conn.Close()

}
