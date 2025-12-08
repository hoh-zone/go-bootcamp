package main

import (
	"log"
	"os"
	"time"

	server "example.com/go-class/14"
)

func main() {
	logger := log.New(os.Stdout, "[secure] ", log.LstdFlags|log.Lmicroseconds)
	cfg := server.Config{
		Addr:         ":8081",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
		Logger:       logger,
		Token:        "demo-token",
		AllowOrigin:  "*",
	}

	srv := server.NewServer(cfg)
	logger.Printf("listening on %s", cfg.Addr)
	if err := srv.ListenAndServe(); err != nil && err.Error() != "http: Server closed" {
		logger.Fatalf("server error: %v", err)
	}
}
