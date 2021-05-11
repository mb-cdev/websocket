package websocket

import (
	"bytes"
	"strings"
)

type frames struct {
	Frames        []*frame
	HasCloseFrame bool
}

func newFramesContainer() frames {
	return frames{HasCloseFrame: false, Frames: make([]*frame, 0)}
}

func (f *frames) Append(fr *frame) {
	if fr.Opcode == CONNECTION_CLOSE {
		f.HasCloseFrame = true
	}

	f.Frames = append(f.Frames, fr)
}

func (f *frames) String() string {
	sb := strings.Builder{}
	for _, fr := range f.Frames {
		fr.UnmaskPayload()
		sb.WriteString(fr.String())
	}
	return sb.String()
}

func (f *frames) Bytes() []byte {
	buff := make([]byte, 0)
	b := bytes.NewBuffer(buff)

	for _, fr := range f.Frames {
		fr.UnmaskPayload()
		b.Write(fr.PayloadData)
	}

	return b.Bytes()
}
