package websocket

import (
	"crypto"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
)

const wsGUID string = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

func getWebSocketAcceptKey(req *http.Request) (string, error) {
	wsKey := req.Header.Get("Sec-WebSocket-Key")

	if len(wsKey) == 0 {
		return "", errors.New("empty sec-websocket-key header")
	}

	key := []byte(strings.Join([]string{wsKey, wsGUID}, ""))
	sha1 := crypto.SHA1.New()
	sha1.Write(key)

	return base64.StdEncoding.EncodeToString(sha1.Sum(nil)), nil
}
