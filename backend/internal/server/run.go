package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/wonderlic-calc/backend/internal/config"
)

// Run builds the router and serves HTTP until ctx is cancelled, then shuts down
// gracefully (15s timeout). Returns nil on clean shutdown, or the serve/shutdown error.
func Run(ctx context.Context, cfg config.Config, logger *slog.Logger) error {
	srv := &http.Server{
		Addr:         cfg.Addr,
		Handler:      NewRouter(cfg, logger),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	errCh := make(chan error, 1)

	go func() {
		logger.Info("server starting", "addr", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	logger.Info("server stopped")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	return srv.Shutdown(shutdownCtx)
}
