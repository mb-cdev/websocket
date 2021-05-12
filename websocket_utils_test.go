package websocket

import (
	"testing"
)

func TestCuttingBytes(t *testing.T) {
	b := []byte{1, 2, 3, 4, 5}
	r := cutByteSlice(b, 3)

	if len(r[0]) != 3 && len(r[1]) != 2 {
		t.Errorf("Error while slicing!")
	}

}
func TestCuttingBytesEmpty(t *testing.T) {
	b := []byte{}
	r := cutByteSlice(b, 1)

	if len(r) != 0 {
		t.Error("Wrong creating bytes from empty")
	}
}
