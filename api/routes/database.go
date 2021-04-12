package routes

import (
	"github.com/go-chi/chi"
	"net/http"
	"encoding/json"
	"log"
	repository "github.com/kwhk/sync/api/repository/redis"
)

// DatabaseRouter for database queries
func DatabaseRouter(userRepo repository.UserRepository, roomRepo repository.RoomRepository) http.Handler {
	r := chi.NewRouter()
	r.Get("/find-room-by-id", func(w http.ResponseWriter, r *http.Request) {
		roomID := r.URL.Query().Get("roomID")
		if room, ok := roomRepo.FindRoomByID(roomID); ok {
			json, err := json.Marshal(room)
			if err != nil {
				log.Println("Could not convert FindRoomByID() to json.")
			} else {
				w.Write(json)
			}
		}

	})
	
	return r
}