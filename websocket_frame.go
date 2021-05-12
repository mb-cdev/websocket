package websocket

import (
	"bytes"
	"encoding/binary"
	"math"
	"math/rand"
	"time"
)

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

//default all frames has FIN false
//and opcode CONTINUATION_FRAME
//change this values after return
func newFrameFromPayload(payload []byte, mask bool) *frame {
	len := len(payload)
	f := newFrame()

	var payloadLen uint64
	if len > 126 && len <= math.MaxUint16 {
		f.PayloadLength7 = 126
	} else if len >= math.MaxUint16 {
		f.PayloadLength7 = 127
	} else {
		f.PayloadLength7 = uint8(len)
		payloadLen = uint64(f.PayloadLength7)
	}

	if f.PayloadLength7 >= 126 {
		f.PayloadLength64 = uint64(len)
		payloadLen = f.PayloadLength64
	}

	f.PayloadData = make([]byte, payloadLen)
	copy(f.PayloadData, payload)

	if mask {
		f.Mask = true
		//generate masking key
		rnd := rand.New(rand.NewSource(time.Now().Unix() + int64(-23438472)))
		for i := 0; i <= 3; i++ {
			f.MaskingKey[i] = byte(rnd.Intn(math.MaxUint8))
		}
		//mask payload data
		for i := range f.PayloadData {
			maskKey := f.MaskingKey[i%4]
			f.PayloadData[i] = f.PayloadData[i] ^ maskKey
		}
	}

	return f
}

func newFrame() *frame {
	return &frame{PayloadData: make([]byte, 0)}
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
	f.payloadUnmasked = true
	return true
}

func (f *frame) IsPayloadUnmasked() bool {
	return !f.Mask || (f.Mask && f.payloadUnmasked)
}

func (f *frame) String() string {
	return string(f.PayloadData)
}

func (f *frame) Bytes() []byte {
	d := make([]byte, 2)

	if f.FIN {
		d[0] |= 1 << 7
	}
	if f.RSV1 {
		d[0] |= 1 << 6
	}
	if f.RSV2 {
		d[0] |= 1 << 5
	}
	if f.RSV3 {
		d[0] |= 1 << 4
	}
	if f.Opcode > 0 {
		d[0] |= byte(f.Opcode)
	}

	d[1] |= f.PayloadLength7

	//convert length to slice of bytes
	b := bytes.NewBuffer([]byte{})
	binary.Write(b, binary.BigEndian, f.PayloadLength64)
	payloadLenBytes := b.Bytes()

	if f.PayloadLength7 == 126 {
		d = append(d, payloadLenBytes[6:]...)
	} else if f.PayloadLength7 == 127 {
		d = append(d, payloadLenBytes...)
	}

	// masking key if any
	if f.Mask {
		d[1] |= 1 << 7
		d = append(d, f.MaskingKey[:]...)
	}

	// payload data
	d = append(d, f.PayloadData...)

	return d
}
