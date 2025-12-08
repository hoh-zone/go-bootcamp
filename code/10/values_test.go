package contextdemo

import (
	"context"
	"testing"
)

func TestRequestID(t *testing.T) {
	ctx := context.Background()
	if _, ok := RequestID(ctx); ok {
		t.Fatalf("expected no request id")
	}
	ctx = WithRequestID(ctx, "req-123")
	if id, ok := RequestID(ctx); !ok || id != "req-123" {
		t.Fatalf("got %q ok=%v, want req-123 true", id, ok)
	}
}
