package routes

import (
	"github.com/go-chi/chi"
	"net/http"
	ws "github.com/kwhk/sync/api/web/websocket"
	"github.com/kwhk/sync/api/repository"
	wp "github.com/kwhk/sync/api/utils/workerPool"
)

// SocketRouter for websocket endpoints
func SocketRouter(userRepo *repository.UserRepository, roomRepo *repository.RoomRepository) http.Handler {
	wsServer := ws.NewWebsocketServer(roomRepo, userRepo)
	msgWorkerPool := wp.NewPool()
	go wsServer.Run()

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(wsServer, w, r)
	})

	return r
}