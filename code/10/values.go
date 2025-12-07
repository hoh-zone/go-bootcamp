package contextdemo

import "context"

type requestIDKey struct{}

// WithRequestID embeds a request id into the context.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, id)
}

// RequestID fetches the request id if present.
func RequestID(ctx context.Context) (string, bool) {
	v := ctx.Value(requestIDKey{})
	if id, ok := v.(string); ok {
		return id, true
	}
	return "", false
}
