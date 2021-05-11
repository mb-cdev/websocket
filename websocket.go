package websocket

import (
	"log"
	"net"
)

func ListenAndServe(addr string, mux *WebSocketMux) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Default().Fatal("Error in ListenAndServe#1", err, "\n")
	}
	defer l.Close()

	for {
		c, errAccept := l.Accept()

		if errAccept != nil {
			log.Default().Fatal("Error in ListenAndServe#2", errAccept, "\n")
		}

		go handleWebSocketConnection(c, mux)
	}
}

/*
func handleConnection(c net.Conn, mux *WebSocketMux) {

	//defer c.Close()

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
	hnd := mux.Dispatch(req.URL)
	if hnd == nil {
		throw404HTTPError(c)
		return
	}

	defaultHandshakeResponse, err := newHandshakeResponse(req)
	if err != nil {
		log.Default().Println("Error in handleConnection#2", err)
		return
	}

	if defaultHandshakeResponse.StatusCode != http.StatusSwitchingProtocols {
		sendHTTPResponse(c, defaultHandshakeResponse)
		log.Default().Println("Error in handleConnection#3", err)
		return
	}

	errResp := sendHTTPResponse(c, defaultHandshakeResponse)

	if errResp != nil {
		log.Default().Fatalln(errResp)
	}

	in := make(chan string)
	out := make(chan string)

	go func(c net.Conn, in chan string, out chan string) {
		readBuff := make([]byte, 0)
		for {
			select {
			case d := <-out:
				fmt.Println([]byte(d))
			default:
				c.SetReadDeadline(time.Now().Add(time.Millisecond * 100))
				buff := make([]byte, 10)
				n, _ := c.Read(buff)

				if n > 0 {
					readBuff = append(readBuff, buff[:n]...)
				}

				if n == 0 && len(readBuff) > 0 {
					fs := newFramesFromBytes(readBuff)

					fmt.Println(len(fs.Bytes()))
					readBuff = make([]byte, 0)
				}

			}
		}
	}(c, in, out)

	go func(hnd Handler, in chan string, out chan string) {
		hnd.ServeConnection(in, out)
	}(hnd, in, out)
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
*/
