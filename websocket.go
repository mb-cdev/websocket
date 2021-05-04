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

func ListenAndServe(addr string, mux *WebSocketMux) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Default().Fatal("Error in ListenAndServe#1", err, "\n")
	}
	defer l.Close()

	for {
		c, errAccept := l.Accept()
		fmt.Println("Accepted")

		if errAccept != nil {
			log.Default().Fatal("Error in ListenAndServe#2", errAccept, "\n")
		}

		go handleConnection(c, mux)
	}
}

func handleConnection(c net.Conn, mux *WebSocketMux) {

	defer c.Close()

	reqBytes := make([]byte, 0)
	buff := make([]byte, 10)
	c.SetReadDeadline(time.Now().Add(time.Millisecond * 500))
	for {
		n, err := c.Read(buff)
		if err != nil || n == 0 {
			fmt.Println(err)
			break
		}
		reqBytes = append(reqBytes, buff...)
	}
	c.SetReadDeadline(time.Time{})

	req, err := http.ReadRequest(bufio.NewReader(bytes.NewBuffer(reqBytes)))
	if err != nil {
		throw404HTTPError(c)
		return
	}

	//find handler
	handler := mux.Dispatch(req.URL)
	//or throw 404
	if handler == nil {
		throw404HTTPError(c)
		return
	}

	defaultHandshakeResponse, err := NewHandshakeResponse(req)
	if err != nil {
		log.Default().Println("Error in handleConnection#2", err)
		return
	}

	if defaultHandshakeResponse.StatusCode != http.StatusSwitchingProtocols {
		sendHTTPResponse(c, defaultHandshakeResponse)
		log.Default().Println("Error in handleConnection#3", err)
		return
	}

	if ok := (*handler).ConfirmHandshake(req, defaultHandshakeResponse); !ok {
		sendHTTPResponse(c, defaultHandshakeResponse)
		log.Default().Println("Error in handleConnection#4", err)
		return
	}

	confirmed := (*handler).ConfirmHandshake(req, defaultHandshakeResponse)
	errResp := sendHTTPResponse(c, defaultHandshakeResponse)

	if errResp != nil {
		log.Default().Fatalln(errResp)
	}

	if !confirmed {
		return
	}

}

func sendHTTPResponse(c net.Conn, resp *http.Response) error {
	errWrite := resp.Write(c)
	return errWrite
}

func throw404HTTPError(c net.Conn) error {
	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n\r\n404 page not found", http.StatusNotFound, http.StatusText(http.StatusNotFound))

	buff := bufio.NewReader(strings.NewReader(statusLine))
	resp, err := http.ReadResponse(buff, nil)

	if err != nil {
		log.Default().Panic("Error in throwing 404", err, "\n")
	}
	errWrite := resp.Write(c)

	return errWrite
}
