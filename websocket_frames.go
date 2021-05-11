package websocket

import (
	"bytes"
	"strings"
)

type Frames []Frame

func (f *Frames) String() string {
	sb := strings.Builder{}
	for _, fr := range *f {
		fr.UnmaskPayload()
		sb.WriteString(fr.String())
	}
	return sb.String()
}

func (f *Frames) Bytes() []byte {
	buff := make([]byte, 0)
	b := bytes.NewBuffer(buff)

	for _, fr := range *f {
		fr.UnmaskPayload()
		b.Write(fr.PayloadData)
	}

	return b.Bytes()
}
