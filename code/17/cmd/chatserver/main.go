package main

import (
	"log"
	"os"
	"time"

	server "example.com/go-class/17"
)

func main() {
	logger := log.New(os.Stdout, "[chat] ", log.LstdFlags|log.Lmicroseconds)
	modelID := os.Getenv("ARK_MODEL_ID")
	if modelID == "" {
		modelID = "deepseek-v3-250324"
	}
	cfg := server.Config{
		Addr:         ":8082",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 0, // streaming
		IdleTimeout:  60 * time.Second,
		Logger:       logger,
		JWTSecret:    "demo-secret",
		ModelID:      modelID,
		AllowOrigin:  "*",
		APIKey:       os.Getenv("ARK_API_KEY"),
	}
	logger.Printf("listening on %s", cfg.Addr)
	srv := server.NewServer(cfg)
	if err := srv.ListenAndServe(); err != nil && err.Error() != "http: Server closed" {
		logger.Fatalf("server error: %v", err)
	}
}
