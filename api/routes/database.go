package routes

import (
	"github.com/go-chi/chi"
	"net/http"
	"encoding/json"
	"log"
	"github.com/kwhk/sync/api/repository"
)

// DatabaseRouter for database queries
func DatabaseRouter(userRepo *repository.UserRepository, roomRepo *repository.RoomRepository) http.Handler {
	r := chi.NewRouter()
	r.Get("/rooms", func(w http.ResponseWriter, r *http.Request) {
		json, err := json.Marshal(roomRepo.GetAllRooms())
		if err != nil {
			log.Println("Could not convert GetAllRooms() to json.")
		} else {
			w.Write(json)
		}
	})

	r.Get("/find-room-by-id", func(w http.ResponseWriter, r *http.Request) {
		roomID := r.URL.Query().Get("roomID")
		json, err := json.Marshal(roomRepo.FindRoomByID(roomID))
		if err != nil {
			log.Println("Could not convert FindRoomByID() to json.")
		} else {
			w.Write(json)
		}
	})
	
	r.Get("/find-room-by-name", func(w http.ResponseWriter, r *http.Request) {
		roomName := r.URL.Query().Get("roomName")
		json, err := json.Marshal(roomRepo.FindRoomByName(roomName))
		if err != nil {
			log.Println("Could not convert FindRoomByName() to json.")
		} else {
			w.Write(json)
		}
	})
	
	return r
}