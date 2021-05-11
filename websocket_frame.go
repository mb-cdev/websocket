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
	payloadUnmasked bool
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

func (f *Frame) UnmaskPayload() bool {
	if !f.Mask || f.payloadUnmasked || len(f.PayloadData) == 0 {
		return false
	}

	for index := range f.PayloadData {
		j := index % 4
		f.PayloadData[index] = f.PayloadData[index] ^ f.MaskingKey[j]
	}

	return true
}

func (f *Frame) IsPayloadUnmasked() bool {
	return !f.Mask || (f.Mask && f.payloadUnmasked)
}

func (f *Frame) String() string {
	return string(f.PayloadData)
}
