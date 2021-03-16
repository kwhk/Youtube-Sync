package routes

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/kwhk/sync/api/repository"
	"github.com/kwhk/sync/api/web/session"
	ws "github.com/kwhk/sync/api/web/websocket"
)

// SocketRouter for websocket endpoints
func SocketRouter(userRepo *repository.UserRepository, roomRepo *repository.RoomRepository, globalSessions *session.Manager) http.Handler {
	wsServer := ws.NewWebsocketServer(roomRepo, userRepo)
	go wsServer.Run()

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		sess := globalSessions.SessionStart(w, r)
		ct := sess.Get("connected")

		if ct == nil {
			sess.Set("connected", true)
		// Only connect to socket if user hasn't previously connected in session
		} else if ct == false {
			sess.Set("connected", true)
			ws.ServeWs(wsServer, w, r)
		}
	})

	return r
}