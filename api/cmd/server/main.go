package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/kwhk/sync/api/config"
	repository "github.com/kwhk/sync/api/repository/redis"
	"github.com/kwhk/sync/api/routes"
	"github.com/kwhk/sync/api/web/session"
	_ "github.com/kwhk/sync/api/web/session/providers/memory"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	connPort = "8000"
	connType = "tcp"
)

// WebServer starts web server
func initWebServer(router http.Handler) {
	fmt.Println("Starting server on port " + connPort)
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

	// Attempt a graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func main() {
	config.CreateRedisClient()
	
	userRepo := repository.UserRepository{Redis: config.Redis}
	roomRepo := repository.RoomRepository{Redis: config.Redis}
	playerRepo := repository.PlayerRepository{Redis: config.Redis}
	globalSessions, _ := session.NewManager("memory", "gosessionid", 3600)
	go globalSessions.GC()

	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Mount("/", routes.IndexRouter(userRepo, roomRepo, playerRepo, globalSessions))

	initWebServer(r)
}
