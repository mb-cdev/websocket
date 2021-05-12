package websocket

import (
	"bufio"
	"bytes"
	"errors"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

var errHandlerNotFound = errors.New("handler not found for address")
var errUnconfirmedServerHandshake = errors.New("handler did not confirm the server handshake")

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

	//find handler
	handler := m.Dispatch(conn.ch.URL)
	if handler == nil {
		return errHandlerNotFound
	}

	if ok := handler.ConfirmHandshake(conn.ch, conn.sh); !ok {
		return errUnconfirmedServerHandshake
	}

	//send handshake to client
	if err := conn.sendHTTPResponse(conn.sh); err != nil {
		return err
	}

	//from client to server
	in := make(chan string, 10)
	//from server to client
	out := make(chan string, 10)
	//exit signal
	disconnect := make(chan bool, 1)

	go func(in chan string, out chan string, disconnect chan bool, conn *webSocketConnection) {
		for {
			select {
			case o := <-out:
				fs, err := newFramesFromPayloadString(o, false)
				if err != nil {
					log.Default().Println(err)
				}
				conn.sendFrames(fs)
			case <-disconnect:
				conn.close()
				return
			default:
				conn.listenForNewFrames(in, disconnect)
			}
		}
	}(in, out, disconnect, &conn)

	//
	go func(handler Handler, in chan string, out chan string, disconnect chan bool) {
		defer func(disconnect chan bool) {
			disconnect <- true
		}(disconnect)

		handler.ServeConnection(in, out, disconnect)
	}(handler, in, out, disconnect)

	return nil
}

func (conn *webSocketConnection) sendFrames(fs *frames) (int, error) {
	sent := 0
	n := 0
	var err error

	for _, f := range fs.Frames {
		n, err = conn.c.Write(f.Bytes())
		sent += n
		if err != nil {
			break
		}
	}

	return sent, err
}

func (conn *webSocketConnection) listenForNewFrames(in chan string, disconnect chan bool) error {
	frameBytes := make([]byte, 0)
	for {
		buff := make([]byte, 1024)
		conn.c.SetReadDeadline(time.Now().Add(time.Millisecond))

		n, err := conn.c.Read(buff)

		if n > 0 {
			frameBytes = append(frameBytes, buff[:n]...)
		}

		if n == 0 && len(frameBytes) > 0 {
			fs := newFramesFromBytes(frameBytes)
			if fs.HasCloseFrame {
				disconnect <- true
				break
			}
			in <- fs.String()
		}

		if err != nil {
			break
		}
	}
	conn.c.SetReadDeadline(time.Time{})
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
		return err2
	}

	return nil
}

func (conn *webSocketConnection) createServerHandshake() error {
	key, err := getWebSocketAcceptKey(conn.ch)

	if err != nil {
		return err
	}

	statusLine := "HTTP/1.1 101 Switching Protocols\r\n\r\n"

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

func (conn *webSocketConnection) sendHTTPResponse(resp *http.Response) error {
	errWrite := resp.Write(conn.c)
	return errWrite
}

func (conn *webSocketConnection) close() {
	fs, _ := newFramesClosingConnection()
	conn.sendFrames(fs)
	conn.c.Close()
}
