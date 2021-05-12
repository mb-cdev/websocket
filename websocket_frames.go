package websocket

import (
	"bytes"
	"errors"
	"strings"
)

const MAX_PAYLOAD_LENGTH int = 3

var errEmptyPayloads error = errors.New("empty payloads")

type frames struct {
	Frames        []*frame
	HasCloseFrame bool
}

func newFramesClosingConnection() (*frames, error) {
	fs, err := newFramesFromPayloadBytes([]byte{}, false)
	fs.Frames[0].Opcode = CONNECTION_CLOSE
	return fs, err
}

func newFramesFromPayloadBytes(payload []byte, mask bool) (*frames, error) {
	payloads := cutByteSlice(payload, MAX_PAYLOAD_LENGTH)

	if len(payloads) == 0 {
		payloads = append(payloads, []byte{})
	}

	frs := &frames{}
	finIndex := len(payloads) - 1
	for i, p := range payloads {
		f := newFrameFromPayload(p, mask)

		if i == 0 {
			f.Opcode = TEXT_FRAME
		} else {
			f.Opcode = CONTINUATION_FRAME
		}

		if i == finIndex {
			f.FIN = true
		}

		frs.Append(f)
	}
	return frs, nil
}

func newFramesFromPayloadString(payload string, mask bool) (*frames, error) {
	return newFramesFromPayloadBytes([]byte(payload), mask)
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
