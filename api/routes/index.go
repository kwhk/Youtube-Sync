package routes

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/kwhk/sync/api/repository"
)

// IndexRouter for all basic routes
func IndexRouter(userRepo *repository.UserRepository, roomRepo *repository.RoomRepository) http.Handler {
	r := chi.NewRouter()	
	r.Mount("/ws", SocketRouter(userRepo, roomRepo))
	r.Mount("/db", DatabaseRouter(userRepo, roomRepo))
	return r
}