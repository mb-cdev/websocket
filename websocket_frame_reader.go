package websocket

import (
	"log"
)

func NewFrameFromBytes(b []byte) []Frame {
	sf := []Frame{}

	var offset uint64 = 0
	for {
		f, r, err := createFrameFromBytes(b[offset:])
		if err != nil {
			log.Default().Println(err)
			break
		}

		sf = append(sf, f)

		if f.FIN {
			break
		}
		offset += r
	}

	return sf
}

func createFrameFromBytes(b []byte) (Frame, uint64, error) {
	f := Frame{}

	//First byte - FIN, RSV1,2,3, OPCODE
	f.FIN = (b[0]&(1<<7) == 1<<7)
	f.RSV1 = (b[0]&(1<<6) == 1<<6)
	f.RSV2 = (b[0]&(1<<5) == 1<<5)
	f.RSV3 = (b[0]&(1<<4) == 1<<4)
	f.Opcode = Opcode(b[0] & 0xF)

	//Second byte - isMask and Payload len
	f.Mask = (b[1]&(1<<7) == 1<<7)

	var payloadLength uint64 = uint64(b[1] & 0x7F)
	f.PayloadLength7 = uint8(payloadLength)

	if payloadLength == 126 {
		payloadLength = 0

		payloadLength = (uint64(b[2]) << 8)
		payloadLength |= uint64(b[3])
	}
	if payloadLength == 127 {
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

	copied := copy(f.PayloadData, b[copyIndexStart:copyIndexEnd])

	if copied != int(payloadLength) {
		log.Fatal("Not copied")
	}

	return f, f.getFrameLength(), nil
}
