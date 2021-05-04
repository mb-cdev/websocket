package websocket

import (
	"fmt"
	"log"
	"net"
	"time"
)

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

func NewFrame() *Frame {
	return &Frame{PayloadData: make([]byte, 0)}
}

func (f *Frame) ReadFromConnection(c *net.Conn) error {
	fmt.Println(c, "inside")
	frame := make([]byte, 0)
	var buff []byte
	for {
		buff = make([]byte, 0)
		n, err := (*c).Read(buff)

		if err != nil {
			log.Default().Println("error while reading from conn", err)
			return err
		}

		if n == 0 {
			fmt.Println(*c)
			break
		}
		frame = append(frame, buff[:n-1]...)
	}

	f.processFrame(frame)
	return nil
}

func (f *Frame) processFrame(frame []byte) {
	time.Sleep(time.Millisecond * 2000)
}
