package routes

import (
	"net/http"

	"github.com/go-chi/chi"
	repository "github.com/kwhk/sync/api/repository/redis"
	"github.com/kwhk/sync/api/web/session"
)


// IndexRouter for all basic routes
func IndexRouter(userRepo repository.UserRepository, roomRepo repository.RoomRepository, globalSessions *session.Manager) http.Handler {
	r := chi.NewRouter()	
	r.Mount("/ws", SocketRouter(userRepo, roomRepo, globalSessions))
	r.Mount("/db", DatabaseRouter(userRepo, roomRepo))
	r.Mount("/test", TestRouter(globalSessions))
	return r
}