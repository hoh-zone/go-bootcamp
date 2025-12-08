package secure

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newTestConfig() Config {
	return Config{
		JWTSecret:   "secret",
		AllowOrigin: "https://example.com",
	}
}

func TestHelloRequiresAPIKey(t *testing.T) {
	mux := NewMux(newTestConfig())
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status=%d want 401", rec.Code)
	}
}

func TestHelloSuccess(t *testing.T) {
	cfg := newTestConfig()
	mux := NewMux(cfg)
	req := httptest.NewRequest(http.MethodGet, "/hello?name=foo", nil)
	token := issueJWT(cfg.JWTSecret, "alice")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d want 200", rec.Code)
	}
	if rec.Body.String() != "hello foo" {
		t.Fatalf("body=%q want %q", rec.Body.String(), "hello foo")
	}
}

func TestEchoValidation(t *testing.T) {
	cfg := newTestConfig()
	mux := NewMux(cfg)
	body, _ := json.Marshal(EchoPayload{Message: ""})
	req := httptest.NewRequest(http.MethodPost, "/echo", bytes.NewReader(body))
	token := issueJWT(cfg.JWTSecret, "alice")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status=%d want 400", rec.Code)
	}
}

func TestEchoContextCancel(t *testing.T) {
	cfg := newTestConfig()
	mux := NewMux(cfg)
	body, _ := json.Marshal(EchoPayload{Message: "hi"})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	req := httptest.NewRequest(http.MethodPost, "/echo", bytes.NewReader(body)).WithContext(ctx)
	token := issueJWT(cfg.JWTSecret, "alice")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusRequestTimeout {
		t.Fatalf("status=%d want timeout", rec.Code)
	}
}

func TestSecurityHeaders(t *testing.T) {
	mux := NewMux(newTestConfig())
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	req.Header.Set("Origin", "https://example.com")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Fatalf("missing security header")
	}
	if rec.Header().Get("Access-Control-Allow-Origin") != "https://example.com" {
		t.Fatalf("cors header not set")
	}
}

func TestHealthzNoAuth(t *testing.T) {
	mux := NewMux(newTestConfig())
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d want 200", rec.Code)
	}
}

func TestLogin(t *testing.T) {
	cfg := newTestConfig()
	mux := NewMux(cfg)
	body, _ := json.Marshal(map[string]string{"username": "alice", "password": "123"})
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d want 200", rec.Code)
	}
	var resp map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp["token"] != cfg.Token {
		t.Fatalf("token=%q want %q", resp["token"], cfg.Token)
	}
}

func TestLoginBadCredentials(t *testing.T) {
	cfg := newTestConfig()
	mux := NewMux(cfg)
	body, _ := json.Marshal(map[string]string{"username": "bad", "password": "bad"})
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status=%d want 401", rec.Code)
	}
}

func TestNewServerConfig(t *testing.T) {
	cfg := Config{
		Addr:         ":0",
		ReadTimeout:  time.Second,
		WriteTimeout: 2 * time.Second,
		IdleTimeout:  3 * time.Second,
	}
	srv := NewServer(cfg)
	if srv.ReadTimeout != cfg.ReadTimeout || srv.WriteTimeout != cfg.WriteTimeout || srv.IdleTimeout != cfg.IdleTimeout {
		t.Fatalf("server timeouts not applied")
	}
}
