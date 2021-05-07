package websocket

type Opcode uint

const (
	CONTINUATION_FRAME Opcode = 0x0
	TEXT_FRAME         Opcode = 0x1
	BINARY_FRAME       Opcode = 0x2
	CONNECTION_CLOSE   Opcode = 0x8
	PING               Opcode = 0x9
	PONG               Opcode = 0xA
)

type Frame struct {
	FIN             bool
	RSV1            bool
	RSV2            bool
	RSV3            bool
	Opcode          Opcode
	Mask            bool
	PayloadLength7  uint8
	PayloadLength64 uint64 // add payload length 126 or 127
	MaskingKey      [4]byte
	PayloadData     []byte
}

func (f *Frame) getHeaderOffsetBytes() uint8 {
	var offset uint8 = 2
	if f.Mask {
		offset += 4
	}

	if f.PayloadLength7 == 126 {
		offset += 2
	} else if f.PayloadLength7 == 127 {
		offset += 8
	}

	return offset
}

func (f *Frame) getFrameLength() uint64 {
	return uint64(f.getHeaderOffsetBytes()) + f.PayloadLength64
}
