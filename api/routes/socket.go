package routes

import (
	"net/http"

	"github.com/go-chi/chi"
	repository "github.com/kwhk/sync/api/repository/redis"
	"github.com/kwhk/sync/api/web/session"
	ws "github.com/kwhk/sync/api/web/websocket"
)

// SocketRouter for websocket endpoints
func SocketRouter(userRepo repository.UserRepository, roomRepo repository.RoomRepository, globalSessions *session.Manager) http.Handler {
	wsServer := ws.NewWebsocketServer(roomRepo, userRepo)
	go wsServer.Run()

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(wsServer, w, r)
	})

	return r
}