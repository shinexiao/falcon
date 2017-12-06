package protocol

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestPacket_HasFlag(t *testing.T) {

	packet := NewPacket(HB_PACKET_BYTE)
	assert.NotNil(t, packet)

	packet.AddFlag(FLAG_CRYPTO)
	assert.Equal(t, packet.HasFlag(FLAG_CRYPTO), true)

	packet.AddFlag(FLAG_COMPRESS)
	assert.Equal(t, packet.HasFlag(FLAG_COMPRESS), true)
}
