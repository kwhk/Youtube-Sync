package main

import (
	"time"

	start "github.com/kwhk/sync/api/server/bin"
	"github.com/kwhk/sync/api/server/routes"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Mount("/api", routes.IndexRouter())
	start.InitWebServer(r)
}