package secure

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Config controls server settings.
type Config struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	Logger       *log.Logger
	// Shared secret for signing JWT.
	JWTSecret   string
	AllowOrigin string
}

// NewMux wires routes and middlewares.
func NewMux(cfg Config) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", HelloHandler)
	mux.HandleFunc("/echo", EchoHandler)
	mux.HandleFunc("/healthz", HealthHandler)
	mux.HandleFunc("/login", LoginHandler(cfg))

	return Chain(
		mux,
		SecurityHeaders(cfg.AllowOrigin),
		BearerAuthMiddleware(cfg.JWTSecret, []string{"/healthz", "/login"}),
		RecoverMiddleware(cfg.Logger),
		LoggingMiddleware(cfg.Logger),
	)
}

// NewServer builds the HTTP server with timeouts.
func NewServer(cfg Config) *http.Server {
	return &http.Server{
		Addr:         cfg.Addr,
		Handler:      NewMux(cfg),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
}

// HelloHandler responds with a greeting; requires GET.
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
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("hello " + name))
}

// EchoPayload is used for validation demo.
type EchoPayload struct {
	Message string `json:"message"`
	Email   string `json:"email,omitempty"`
}

// EchoHandler validates the payload and echoes it back.
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
	if err := validatePayload(payload); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

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

func validatePayload(p EchoPayload) error {
	if strings.TrimSpace(p.Message) == "" {
		return errors.New("message is required")
	}
	if p.Email != "" && !strings.Contains(p.Email, "@") {
		return errors.New("email must contain @")
	}
	return nil
}

// HealthHandler returns 200 OK.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

// LoginHandler issues a bearer token after verifying credentials.
func LoginHandler(cfg Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		defer r.Body.Close()
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
			return
		}
		if req.Username == "" || req.Password == "" {
			http.Error(w, "username and password required", http.StatusBadRequest)
			return
		}
		const expectedUser = "alice"
		const expectedPass = "123"
		if req.Username != expectedUser || req.Password != expectedPass {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		token := issueJWT(cfg.JWTSecret, req.Username)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}

// ErrServerClosed signals server shutdown.
var ErrServerClosed = errors.New("server closed")

func issueJWT(secret, sub string) string {
	claims := jwt.RegisteredClaims{
		Subject:   sub,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return ""
	}
	return signed
}
