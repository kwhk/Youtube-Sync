package websocket 

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { return true },
}

func upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
			return nil, err
	}

	return conn, nil
}

// ServeWs defines our WebSocket endpoint
func ServeWs(wsServer *WsServer, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+V\n", err)
	}

	client := initClient(conn, wsServer)

	// Allow connection of memory referenced by the caller by doing
	// all work in new goroutines.
	go client.writePump()
	go client.readPump()

	wsServer.register <- client
}

