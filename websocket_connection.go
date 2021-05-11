package websocket

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

type webSocketConnection struct {
	c net.Conn
	m *WebSocketMux

	//client handshake
	ch *http.Request
	//server handshake
	sh *http.Response
}

func handleWebSocketConnection(c net.Conn, m *WebSocketMux) error {
	conn := webSocketConnection{c: c, m: m}

	if err := conn.waitAndParseClientHandshake(); err != nil {
		return err
	}

	if err := conn.createServerHandshake(); err != nil {
		return err
	}

	conn.sendHTTPResponse(conn.sh)

	return nil
}

func (conn *webSocketConnection) waitAndParseClientHandshake() error {
	reqBytes := make([]byte, 0)
	buff := make([]byte, 1024)
	var err2 error

	conn.c.SetReadDeadline(time.Now().Add(time.Millisecond * 500))

	for {
		n, err := conn.c.Read(buff)
		if err != nil || n == 0 {
			log.Default().Println(err)
			break
		}
		reqBytes = append(reqBytes, buff[:n]...)
	}

	conn.c.SetReadDeadline(time.Time{})

	conn.ch, err2 = http.ReadRequest(bufio.NewReader(bytes.NewBuffer(reqBytes)))
	if err2 != nil {
		conn.throw404HTTPError()
		return err2
	}

	return nil
}

func (conn *webSocketConnection) createServerHandshake() error {
	key, err := getWebSocketAcceptKey(conn.ch)
	var statusLine string

	if err != nil {
		log.Default().Println("Bad client handshake")
		statusLine = "HTTP/1.1 400 Bad Request\r\n\r\n"
	} else {
		statusLine = "HTTP/1.1 101 Switching Protocols\r\n\r\n"
	}

	buff := bufio.NewReader(strings.NewReader(statusLine))
	conn.sh, err = http.ReadResponse(buff, conn.ch)

	if err != nil {
		return err
	}

	conn.sh.Header.Add("Sec-WebSocket-Accept", key)
	conn.sh.Header.Add("Upgrade", "websocket")
	conn.sh.Header.Add("Connection", "Upgrade")

	return nil
}

func (conn *webSocketConnection) throw404HTTPError() error {
	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n\r\n404 page not found", http.StatusNotFound, http.StatusText(http.StatusNotFound))

	buff := bufio.NewReader(strings.NewReader(statusLine))
	resp, err := http.ReadResponse(buff, nil)

	if err != nil {
		log.Default().Panic("Error in throwing 404", err, "\n")
	}
	errWrite := conn.sendHTTPResponse(resp)

	return errWrite
}

func (conn *webSocketConnection) sendHTTPResponse(resp *http.Response) error {
	errWrite := resp.Write(conn.c)
	return errWrite
}

func (conn *webSocketConnection) close() {
	conn.c.Close()
}
