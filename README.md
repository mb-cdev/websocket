# websocket
Module for handling websocket connections

# Usage

```golang
// sample websocket handler 
type TestWebsocketHandler struct {
}

func (t *TestWebsocketHandler) ConfirmHandshake(clientHandshakeRequest *http.Request, serverHandshakeResponse *http.Response) bool {
	fmt.Println("OK!")
	return true
}

func (t *TestWebsocketHandler) ServeConnection(inData chan string, outData chan string, disconnect chan bool) {

	fmt.Println("Listening...")
	for {
		select {
		case <-disconnect:
			fmt.Println("Return from fc...")
			return
		case d := <-inData:
			fmt.Println("In Data: ", d)
			outData <- fmt.Sprintln("Received:", d)
			if d == "exit" {
				disconnect <- true
			}
		}
	}
}

func main(){
    // Create server and listen...
    mux := websocket.NewWebSocketMux()
    mux.Handle("/test", &TestWebsocketHandler{})
  
    websocket.ListenAndServe("0.0.0.0:9999", mux)
}
