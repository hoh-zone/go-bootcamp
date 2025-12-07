package chatserver

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// LoggingMiddleware logs method, path and duration.
func LoggingMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			if logger != nil {
				logger.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
			}
		})
	}
}

// RecoverMiddleware prevents panics from crashing the server.
func RecoverMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					if logger != nil {
						logger.Printf("panic: %v", rec)
					}
					http.Error(w, "internal error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// BearerAuthMiddleware validates JWT Bearer tokens unless path is allowlisted.
func BearerAuthMiddleware(secret string, allowlist []string) func(http.Handler) http.Handler {
	skip := make(map[string]struct{}, len(allowlist))
	for _, p := range allowlist {
		skip[p] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if secret == "" {
				next.ServeHTTP(w, r)
				return
			}
			if _, ok := skip[r.URL.Path]; ok {
				next.ServeHTTP(w, r)
				return
			}
			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			raw := strings.TrimPrefix(auth, "Bearer ")
			claims, err := parseJWT(secret, raw)
			if err != nil || claims == nil || claims.Subject == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func parseJWT(secret, token string) (*jwt.RegisteredClaims, error) {
	parsed, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := parsed.Claims.(*jwt.RegisteredClaims); ok && parsed.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid claims")
}

// SecurityHeaders adds common defense headers and simple CORS.
func SecurityHeaders(allowOrigin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("Content-Security-Policy", "default-src 'self'")
			if allowOrigin != "" {
				origin := r.Header.Get("Origin")
				if origin != "" && (allowOrigin == "*" || strings.EqualFold(origin, allowOrigin)) {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Vary", "Origin")
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

// Chain applies middlewares from left to right.
func Chain(h http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}
