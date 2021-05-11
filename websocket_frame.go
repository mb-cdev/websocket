package websocket

type opcode uint8

const (
	CONTINUATION_FRAME opcode = 0x0
	TEXT_FRAME         opcode = 0x1
	BINARY_FRAME       opcode = 0x2
	CONNECTION_CLOSE   opcode = 0x8
	PING               opcode = 0x9
	PONG               opcode = 0xA
)

type frame struct {
	FIN             bool
	RSV1            bool
	RSV2            bool
	RSV3            bool
	Opcode          opcode
	Mask            bool
	PayloadLength7  uint8
	PayloadLength64 uint64 // extendend payload length 126 or 127
	MaskingKey      [4]byte
	PayloadData     []byte
	payloadUnmasked bool
}

func (f *frame) getHeaderOffsetBytes() uint8 {
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

func (f *frame) getFrameLength() uint64 {
	return uint64(f.getHeaderOffsetBytes()) + f.PayloadLength64
}

func (f *frame) UnmaskPayload() bool {
	if !f.Mask || f.payloadUnmasked || len(f.PayloadData) == 0 {
		return false
	}

	for index := range f.PayloadData {
		j := index % 4
		f.PayloadData[index] = f.PayloadData[index] ^ f.MaskingKey[j]
	}

	return true
}

func (f *frame) IsPayloadUnmasked() bool {
	return !f.Mask || (f.Mask && f.payloadUnmasked)
}

func (f *frame) String() string {
	return string(f.PayloadData)
}
