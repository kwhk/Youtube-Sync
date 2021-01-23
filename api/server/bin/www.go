/* Start up script to set up web server */

package bin

import (
	"fmt"
	"os"
	"os/signal"
	"net/http"
	"context"
	"time"
)

const (
	connPort = "8000"
	connHost = "localhost"
	connType = "tcp"
	connURL = connHost + ":" + connPort
)

// InitWebServer starts web server
func InitWebServer(router http.Handler) {
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