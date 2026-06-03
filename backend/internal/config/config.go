// Package config provides runtime configuration loaded from environment variables.
package config

import (
	"os"
	"strings"
	"time"
)

// Config holds all runtime configuration for the server.
type Config struct {
	// Addr is the TCP address the HTTP server listens on (e.g. ":8080").
	Addr string

	// AllowedOrigins is the list of CORS-permitted origins.
	AllowedOrigins []string

	// ReadTimeout is the maximum duration for reading an entire request.
	ReadTimeout time.Duration

	// WriteTimeout is the maximum duration before timing out writes of the response.
	WriteTimeout time.Duration

	// IdleTimeout is the maximum amount of time to wait for the next request on a keep-alive connection.
	IdleTimeout time.Duration
}

// Load reads configuration from the environment with sane defaults.
//
// Environment variables:
//   - ADDR: TCP listen address; default ":8080".
//   - CORS_ALLOWED_ORIGINS: comma-separated list of allowed origins; default "http://localhost:5173".
//
// Timeout values are constants and are not environment-driven:
//   - ReadTimeout  = 5s
//   - WriteTimeout = 10s
//   - IdleTimeout  = 120s
func Load() Config {
	return Config{
		Addr:           addr(),
		AllowedOrigins: allowedOrigins(),
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
	}
}

func addr() string {
	if v := os.Getenv("ADDR"); v != "" {
		return v
	}
	return ":8080"
}

func allowedOrigins() []string {
	raw, ok := os.LookupEnv("CORS_ALLOWED_ORIGINS")
	if !ok {
		return []string{"http://localhost:5173"}
	}
	var origins []string
	for _, part := range strings.Split(raw, ",") {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			origins = append(origins, trimmed)
		}
	}
	if len(origins) == 0 {
		return []string{"http://localhost:5173"}
	}
	return origins
}
