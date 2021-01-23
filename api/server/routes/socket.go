package routes

import (
	"github.com/go-chi/chi"
	"net/http"
	"github.com/kwhk/sync/api/server/bin/websocket"
)

// SocketRouter for websocket endpoints
func SocketRouter() http.Handler {
	room := websocket.NewRoom()
	go room.Start()

	r := chi.NewRouter()
	r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(room, w, r)
	})

	return r
}