package routes

import (
	"github.com/go-chi/chi"
	"net/http"
)

// IndexRouter for all basic routes
func IndexRouter() http.Handler {
	r := chi.NewRouter()
	r.Mount("/", SocketRouter())
	return r
}