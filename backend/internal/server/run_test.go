package server

import (
	"context"
	"io"
	"log/slog"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wonderlic-calc/backend/internal/config"
)

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func freeAddr(t *testing.T) string {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	addr := l.Addr().String()
	require.NoError(t, l.Close())
	return addr
}

func testConfig(addr string) config.Config {
	return config.Config{
		Addr:           addr,
		AllowedOrigins: []string{"*"},
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
	}
}

func TestRun_CleanShutdown(t *testing.T) {
	addr := freeAddr(t)
	cfg := testConfig(addr)
	logger := discardLogger()

	ctx, cancel := context.WithCancel(context.Background())

	runErr := make(chan error, 1)
	go func() {
		runErr <- Run(ctx, cfg, logger)
	}()

	// Wait for server to be ready.
	require.Eventually(t, func() bool {
		resp, err := http.Get("http://" + addr + "/healthz")
		if err != nil {
			return false
		}
		resp.Body.Close()
		return resp.StatusCode == http.StatusOK
	}, 3*time.Second, 50*time.Millisecond, "server did not become ready")

	// Confirm /healthz returns 200.
	resp, err := http.Get("http://" + addr + "/healthz")
	require.NoError(t, err)
	resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Cancel ctx → graceful shutdown.
	cancel()

	select {
	case err := <-runErr:
		assert.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("Run did not return after ctx cancel")
	}
}

func TestRun_ServeError(t *testing.T) {
	// Keep the port occupied so ListenAndServe fails immediately.
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer l.Close()

	cfg := testConfig(l.Addr().String())
	logger := discardLogger()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	runErr := make(chan error, 1)
	go func() {
		runErr <- Run(ctx, cfg, logger)
	}()

	select {
	case err := <-runErr:
		assert.Error(t, err)
	case <-time.After(3 * time.Second):
		t.Fatal("Run did not return a serve error within timeout")
	}
}
