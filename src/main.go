package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jnsoft/htmxgo/src/logger"
	"github.com/jnsoft/htmxgo/src/server"
)

func main() {
	log.SetFlags(0) // Disable default flags
	log.SetOutput(&logger.IsoLogWriter{})

	mux := server.NewServer()
	addr := ":8080"
	srv := &http.Server{
		Addr:    addr,
		Handler: logger.LoggingMiddleware(mux),
	}

	go func() {
		log.Printf("Starting server on %s\n", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %s\n", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %s\n", err)
	}

	log.Println("Server exited gracefully")

}
