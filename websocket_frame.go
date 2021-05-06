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
	FIN           bool
	RSV1          bool
	RSV2          bool
	RSV3          bool
	Opcode        uint8
	Mask          bool
	PayloadLength uint64
	MaskingKey    [4]byte
	PayloadData   []byte
}

func NewFrameFromBytes(b []byte) Frame {
	f := Frame{}

	//First byte - FIN, RSV1,2,3, OPCODE
	f.FIN = (b[0]&(1<<7) == 1<<7)
	f.RSV1 = (b[0]&(1<<6) == 1<<6)
	f.RSV2 = (b[0]&(1<<5) == 1<<5)
	f.RSV3 = (b[0]&(1<<4) == 1<<4)
	f.Opcode = (b[0] & 0xF)

	//Second byte - isMask and Payload len
	f.Mask = (b[1]&(1<<7) == 1<<7)

	var payloadLength uint64 = uint64((b[1] & 0x7F))

	if payloadLength == 126 {
		payloadLength = 0
		payloadLength = (uint64(b[2]) << 8)
		payloadLength |= uint64(b[3])
	}
	if payloadLength == 127 {
		payloadLength = 0

		payloadLength |= (uint64(b[2]) << 56)
		payloadLength |= (uint64(b[3]) << 48)
		payloadLength |= (uint64(b[4]) << 40)
		payloadLength |= (uint64(b[5]) << 32)
		payloadLength |= (uint64(b[6]) << 24)
		payloadLength |= (uint64(b[7]) << 16)
		payloadLength |= (uint64(b[8]) << 8)
		payloadLength |= uint64(b[9])
	}

	f.PayloadLength = payloadLength
	return f
}
