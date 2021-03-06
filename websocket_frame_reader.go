package websocket

import (
	"errors"
	"log"
)

var errBadFrameBytes = errors.New("bad frame bytes")

func newFramesFromBytes(b []byte) frames {
	sf := newFramesContainer()

	var offset uint64 = 0
	for {
		f, r, err := createFrameFromBytes(b[offset:])
		if err != nil {
			log.Default().Println(err)
			break
		}

		sf.Append(f)

		if f.FIN {
			break
		}
		offset += r
	}

	return sf
}

func createFrameFromBytes(b []byte) (*frame, uint64, error) {
	f := &frame{}

	if len(b) < 2 {
		return nil, 0, errBadFrameBytes
	}

	//First byte - FIN, RSV1,2,3, OPCODE
	f.FIN = (b[0]&(1<<7) == 1<<7)
	f.RSV1 = (b[0]&(1<<6) == 1<<6)
	f.RSV2 = (b[0]&(1<<5) == 1<<5)
	f.RSV3 = (b[0]&(1<<4) == 1<<4)
	f.Opcode = opcode(b[0] & 0xF)

	//Second byte - isMask and Payload len
	f.Mask = (b[1]&(1<<7) == 1<<7)

	var payloadLength uint64 = uint64(b[1] & 0x7F)
	f.PayloadLength7 = uint8(payloadLength)

	if payloadLength == 126 {

		if len(b) < 4 {
			return nil, 0, errBadFrameBytes
		}

		payloadLength = 0

		payloadLength = (uint64(b[2]) << 8)
		payloadLength |= uint64(b[3])
	}
	if payloadLength == 127 {
		if len(b) < 10 {
			return nil, 0, errBadFrameBytes
		}
		payloadLength = 0

		for i := 2; i <= 9; i++ {
			payloadLength |= uint64(b[i]) << (56 - (i-2)*8)
		}
	}

	f.PayloadLength64 = payloadLength

	//Payload data
	headerOffset := f.getHeaderOffsetBytes()
	f.PayloadData = make([]byte, f.PayloadLength64)

	copyIndexStart := uint64(headerOffset)
	copyIndexEnd := f.PayloadLength64 + copyIndexStart

	if len(b) < int(copyIndexEnd) {
		return nil, 0, errBadFrameBytes
	}

	copied := copy(f.PayloadData, b[copyIndexStart:copyIndexEnd])

	if copied != int(payloadLength) {
		log.Fatal("Not copied")
	}

	//unmask if masked
	if f.Mask {
		setMask(f, b)
	}

	return f, f.getFrameLength(), nil
}

func setMask(f *frame, b []byte) error {
	//the mask is at end of frame header
	//but before payload data

	headerOffset := f.getHeaderOffsetBytes()
	//mask is always 4 bytes len
	maskStartIndex := headerOffset - 4

	f.MaskingKey[0] = b[maskStartIndex]
	f.MaskingKey[1] = b[maskStartIndex+1]
	f.MaskingKey[2] = b[maskStartIndex+2]
	f.MaskingKey[3] = b[maskStartIndex+3]

	return nil
}
