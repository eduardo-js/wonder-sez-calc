// Command server starts the wonderlic-calc backend HTTP server.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/wonderlic-calc/backend/internal/config"
	"github.com/wonderlic-calc/backend/internal/server"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	if err := server.Run(ctx, config.Load(), logger); err != nil {
		logger.Error("server terminated", "error", err)
		os.Exit(1)
	}
}
