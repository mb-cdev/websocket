package websocket

import (
	"fmt"
	"testing"
)

func TestNewFrameFromPayloadUnmasked(t *testing.T) {
	//f := newFrameFromPayload([]byte{1, 2, 3, 4, 5}, false)
	//fmt.Printf("Unmasked:\n %#v\n\n", f)
}
func TestNewFrameFromPayloadMasked(t *testing.T) {
	/*dat := make([]byte, 0)
	for i := 1; i <= 6; i++ {
		dat = append(dat, byte(i))
	}
	f := newFrameFromPayload(dat, true)

	fmt.Printf("Masked:\n%#v\n", f)
	f.UnmaskPayload()
	fmt.Printf("Unmasked:\n%#v\n", f)

	f.FIN = true
	f.Opcode = TEXT_FRAME
	fmt.Printf("Unmasked to bytes %b", f.Bytes())*/
}

func TestNewFramesFromPayload(t *testing.T) {
	/*data := make([]byte, 0)
	for i := 0; i <= 4; i++ {
		data = append(data, byte(i))
	}

	fs, _ := newFramesFromPayloadBytes(data, true)
	for _, v := range fs.Frames {
		fmt.Printf("%#b\n", v.Bytes())
	}*/

}

func TestNewFramesFromEmptyPayload(t *testing.T) {
	fs, _ := newFramesFromPayloadBytes([]byte{}, false)
	fmt.Printf("%#v", fs.Frames[0])
}
