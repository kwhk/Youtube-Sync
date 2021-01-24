package main

import (
	"fmt"
	"os"
	"os/signal"
	"net/http"
	"context"
	"time"

	"github.com/kwhk/sync/api/server/init"
	"github.com/kwhk/sync/api/server/routes"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	connPort = "8000"
	connHost = "localhost"
	connType = "tcp"
	connURL = connHost + ":" + connPort
)

// WebServer starts web server
func initWebServer(router http.Handler) {
	fmt.Println("Starting server on " + connURL)
	srv := &http.Server {
		Addr: ":" + connPort,
		Handler: router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		srv.ListenAndServe()
	}()

	// Wait for an interrupt
	c := make(chan os.Signal, 1)
	// send os interrupt signals only to channel c
	signal.Notify(c, os.Interrupt)
	//
	<-c

	// Attempt a graceful shutdon
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Mount("/api", routes.IndexRouter())
	initWebServer(r)
}