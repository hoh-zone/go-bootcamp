package server

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHelloHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/hello?name=foo", nil)
	rec := httptest.NewRecorder()
	HelloHandler(rec, req)
	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status=%d want 200", res.StatusCode)
	}
	body := rec.Body.String()
	if body != "hello foo" {
		t.Fatalf("body=%q want %q", body, "hello foo")
	}
}

func TestHelloHandlerMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/hello", nil)
	rec := httptest.NewRecorder()
	HelloHandler(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status=%d want 405", rec.Code)
	}
}

func TestEchoHandler(t *testing.T) {
	payload := EchoPayload{Message: "hi"}
	buf, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/echo", bytes.NewReader(buf))
	rec := httptest.NewRecorder()
	EchoHandler(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d want 200", rec.Code)
	}
	var got EchoPayload
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Message != payload.Message {
		t.Fatalf("message=%q want %q", got.Message, payload.Message)
	}
}

func TestEchoHandlerBadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/echo", bytes.NewReader([]byte(`{}`)))
	rec := httptest.NewRecorder()
	EchoHandler(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status=%d want 400", rec.Code)
	}
}

func TestEchoHandlerContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	req := httptest.NewRequest(http.MethodPost, "/echo", bytes.NewReader([]byte(`{"message":"x"}`))).WithContext(ctx)
	rec := httptest.NewRecorder()
	EchoHandler(rec, req)
	if rec.Code != http.StatusRequestTimeout {
		t.Fatalf("status=%d want 408-ish", rec.Code)
	}
}

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	HealthHandler(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d want 200", rec.Code)
	}
	if rec.Body.String() != "ok" {
		t.Fatalf("body=%q want ok", rec.Body.String())
	}
}

func TestRecoverMiddleware(t *testing.T) {
	panicHandler := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic("boom")
	})
	h := Chain(panicHandler, RecoverMiddleware(log.Default()))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status=%d want 500", rec.Code)
	}
}

func TestNewServerConfig(t *testing.T) {
	cfg := Config{
		Addr:         ":0",
		ReadTimeout:  time.Second,
		WriteTimeout: 2 * time.Second,
		IdleTimeout:  3 * time.Second,
		Logger:       log.Default(),
	}
	srv := NewServer(cfg)
	if srv.ReadTimeout != cfg.ReadTimeout || srv.WriteTimeout != cfg.WriteTimeout || srv.IdleTimeout != cfg.IdleTimeout {
		t.Fatalf("server timeouts not applied")
	}
}
