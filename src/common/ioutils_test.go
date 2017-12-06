package common

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCompress(t *testing.T) {
	data := Compress([]byte("Hello world"))
	buff := []byte{120, 156, 242, 72, 205, 201, 201, 87, 40, 207, 47, 202, 73, 1, 4, 0, 0, 255, 255, 24, 171, 4, 61}
	assert.Equal(t, data, buff)
}

func TestUnCompress(t *testing.T) {

	buff := []byte{120, 156, 242, 72, 205, 201, 201, 87, 40, 207, 47, 202, 73, 1, 4, 0, 0, 255, 255, 24, 171, 4, 61}

	assert.Equal(t, string(UnCompress(buff)), "Hello world")
}
