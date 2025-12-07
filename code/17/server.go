package chatserver

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	openai "github.com/sashabaranov/go-openai"
)

type Config struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	Logger       *log.Logger
	JWTSecret    string
	ModelID      string
	AllowOrigin  string
	APIKey       string
}

func NewServer(cfg Config) *http.Server {
	return &http.Server{
		Addr:         cfg.Addr,
		Handler:      NewMux(cfg),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
}

func NewMux(cfg Config) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/login", LoginHandler(cfg))
	mux.HandleFunc("/chat", ChatHandler(cfg))
	mux.HandleFunc("/healthz", HealthHandler)

	return Chain(
		mux,
		SecurityHeaders(cfg.AllowOrigin),
		BearerAuthMiddleware(cfg.JWTSecret, []string{"/login", "/healthz"}),
		RecoverMiddleware(cfg.Logger),
		LoggingMiddleware(cfg.Logger),
	)
}

// LoginHandler issues a JWT for fixed demo credentials.
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
		if req.Username != "alice" || req.Password != "123" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		token := issueJWT(cfg.JWTSecret, req.Username)
		if token == "" {
			http.Error(w, "failed to sign token", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}

// ChatRequest defines input for the SSE chat endpoint.
type ChatRequest struct {
	Message string `json:"message"`
	Model   string `json:"model,omitempty"`
}

// ChatHandler proxies to the Ark streaming API and returns SSE chunks.
func ChatHandler(cfg Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming unsupported", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		var req ChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
			return
		}
		if req.Message == "" {
			http.Error(w, "message is required", http.StatusBadRequest)
			return
		}

		apiKey := cfg.APIKey
		if envKey := os.Getenv("ARK_API_KEY"); apiKey == "" && envKey != "" {
			apiKey = envKey
		}
		if apiKey == "" {
			http.Error(w, "server not configured: missing ARK_API_KEY", http.StatusInternalServerError)
			return
		}

		model := cfg.ModelID
		if req.Model != "" {
			model = req.Model
		}

		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		stream, err := newArkClient(apiKey).CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleSystem, Content: "你是人工智能助手"},
				{Role: openai.ChatMessageRoleUser, Content: req.Message},
			},
		})
		if err != nil {
			http.Error(w, "chat stream error: "+err.Error(), http.StatusBadGateway)
			return
		}
		defer stream.Close()

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		for {
			select {
			case <-ctx.Done():
				writeSSE(w, "event: error\ndata: "+ctx.Err().Error()+"\n\n")
				flusher.Flush()
				return
			default:
			}
			resp, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				writeSSE(w, "event: done\ndata: [DONE]\n\n")
				flusher.Flush()
				return
			}
			if err != nil {
				writeSSE(w, "event: error\ndata: "+err.Error()+"\n\n")
				flusher.Flush()
				return
			}
			if len(resp.Choices) > 0 {
				chunk := resp.Choices[0].Delta.Content
				writeSSE(w, "data: "+chunk+"\n\n")
				flusher.Flush()
			}
		}
	}
}

// HealthHandler returns 200 OK.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func writeSSE(w http.ResponseWriter, payload string) {
	_, _ = w.Write([]byte(payload))
}

func issueJWT(secret, sub string) string {
	if secret == "" {
		secret = "demo-secret"
	}
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

func newArkClient(apiKey string) *openai.Client {
	cfg := openai.DefaultConfig(apiKey)
	cfg.BaseURL = "https://ark.cn-beijing.volces.com/api/v3"
	return openai.NewClientWithConfig(cfg)
}
