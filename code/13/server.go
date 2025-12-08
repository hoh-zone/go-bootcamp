package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Config controls server settings.
type Config struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	Logger       *log.Logger
}

// NewMux builds the mux with routes and middlewares.
func NewMux(logger *log.Logger) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", HelloHandler)
	mux.HandleFunc("/echo", EchoHandler)
	mux.HandleFunc("/healthz", HealthHandler)

	// chain middlewares: recover then logging
	return Chain(mux, RecoverMiddleware(logger), LoggingMiddleware(logger))
}

// NewServer creates *http.Server configured with timeouts.
func NewServer(cfg Config) *http.Server {
	return &http.Server{
		Addr:         cfg.Addr,
		Handler:      NewMux(cfg.Logger),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
}

// HelloHandler responds with a greeting; defaults name to "gopher".
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "gopher"
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "hello %s", name)
}

// EchoPayload is the expected request/response shape for /echo.
type EchoPayload struct {
	Message string `json:"message"`
}

// EchoHandler echos posted JSON payload.
func EchoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	var payload EchoPayload
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&payload); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	if payload.Message == "" {
		http.Error(w, "message is required", http.StatusBadRequest)
		return
	}

	// Respect cancellation: if ctx is done before writing, abort.
	select {
	case <-r.Context().Done():
		http.Error(w, r.Context().Err().Error(), http.StatusRequestTimeout)
		return
	default:
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(payload)
}

// HealthHandler returns 200 OK for probes.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

// ErrServerClosed signals server shutdown.
var ErrServerClosed = errors.New("server closed")
