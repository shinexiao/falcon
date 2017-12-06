package protocol

import (
	"testing"
	"bytes"
	"fmt"
)

func TestEncodePacket(t *testing.T) {

	packet := NewPacket(HB_PACKET_BYTE)

	buf := new(bytes.Buffer)
	EncodePacket(packet,buf)

	fmt.Println(buf.Len())
	fmt.Println(buf.Cap())
}
