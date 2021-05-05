package websocket

import (
	"log"
	"net/http"
	"net/url"
	"reflect"
)

type Handler interface {
	ConfirmHandshake(clientHandshakeRequest *http.Request, serverHandshakeResponse *http.Response) bool
	ServeConnection(inData chan string, outData chan string)
}

type WebSocketMux struct {
	handlers map[string]Handler
}

func NewWebSocketMux() *WebSocketMux {
	return &WebSocketMux{handlers: make(map[string]Handler)}
}

func (w *WebSocketMux) Handle(pattern string, handle Handler) {
	if _, ok := w.handlers[pattern]; ok {
		log.Default().Panic("Route already assigned")
	}
	w.handlers[pattern] = handle
}

func (w *WebSocketMux) Dispatch(url *url.URL) Handler {
	if h, ok := w.handlers[url.Path]; ok {

		typ := reflect.TypeOf(h).Elem()
		t := reflect.New(typ)

		return t.Interface().(Handler)
	}

	return nil
}
