package websocket

import (
	"bufio"
	"crypto"
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"strings"
)

const wsGUID string = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

func newHandshakeResponse(clientHandshakeRequest *http.Request) (*http.Response, error) {
	key, err := getWebSocketAcceptKey(clientHandshakeRequest)
	var statusLine string

	if err != nil {
		log.Default().Println("Bad client handshake")
		statusLine = "HTTP/1.1 400 Bad Request\r\n\r\n"
	} else {
		statusLine = "HTTP/1.1 101 Switching Protocols\r\n\r\n"
	}

	buff := bufio.NewReader(strings.NewReader(statusLine))
	resp, err := http.ReadResponse(buff, clientHandshakeRequest)

	if err != nil {
		return nil, err
	}

	resp.Header.Add("Sec-WebSocket-Accept", key)
	resp.Header.Add("Upgrade", "websocket")
	resp.Header.Add("Connection", "Upgrade")

	return resp, nil
}

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
