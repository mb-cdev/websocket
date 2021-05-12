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

		go func(c net.Conn, mux *WebSocketMux) {
			err := handleWebSocketConnection(c, mux)
			if err != nil {
				log.Default().Println(err)
				c.Close()
			}
		}(c, mux)
	}
}
