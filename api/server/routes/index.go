package routes

import (
	"github.com/go-chi/chi"
	"net/http"
	"github.com/kwhk/sync/api/server/utils"
)

// HelloResponse type
type HelloResponse struct {
	Message string `json:"message"`
}

// HelloWorld returns Hello World response
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	response := HelloResponse {
		Message: "Hello there!",
	}

	utils.JSONResponse(w, response, http.StatusOK)
}

// IndexRouter for all basic routes
func IndexRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", HelloWorld)
	r.Mount("/", SocketRouter())
	return r
}