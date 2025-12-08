package main

import (
	"log"
	"os"
	"time"

	server "example.com/go-class/13"
)

func main() {
	logger := log.New(os.Stdout, "[http] ", log.LstdFlags|log.Lmicroseconds)
	cfg := server.Config{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
		Logger:       logger,
	}

	srv := server.NewServer(cfg)
	logger.Printf("listening on %s", cfg.Addr)
	if err := srv.ListenAndServe(); err != nil && err.Error() != "http: Server closed" {
		logger.Fatalf("server error: %v", err)
	}
}
