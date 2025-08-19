package graceful

import (
	"context"
	"fmt"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

type Process interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

func Graceful(processes map[string]Process) {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	startgroup, ctx := errgroup.WithContext(ctx)

	for name, process := range processes {
		process := process
		startgroup.Go(func() error {
			slog.Info(fmt.Sprintf("starting %v", name))
			return process.Start(ctx)
		})
	}

	if err := startgroup.Wait(); err != nil {
		slog.Error("error starting", slog.String("error", err.Error()))
	}

	<-ctx.Done()

	slog.Info("received termination signal, shutting down...")

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			slog.Info("waiting for graceful shutdown...")
		}
	}()

	stopCtx, stopCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer stopCancel()

	stopgroup, stopCtx := errgroup.WithContext(stopCtx)

	for name, process := range processes {
		process := process
		stopgroup.Go(func() error {
			slog.Info(fmt.Sprintf("stopping %v", name))
			return process.Stop(stopCtx)
		})
	}

	if err := stopgroup.Wait(); err != nil {
		slog.Error("error stopping", slog.String("error", err.Error()))
	}

	slog.Info("shutdown complete")
}
