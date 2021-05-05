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
	Opcode        [4]byte
	Mask          bool
	PayloadLenght uint64
	MaskingKey    [4]byte
	PayloadData   []byte
}

func FromWebSocketFrame(f []byte) Frame {
	return Frame{}
}
