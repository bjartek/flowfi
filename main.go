package main

import (
	"context"
	"flowfi_tg_bot/pkg/bot"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func main() {
	flowFi := bot.NewBot()

	// Create a context that listens for interrupt signals
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	eg, ctx := errgroup.WithContext(ctx)

	// Command listener
	eg.Go(func() error {
		return flowFi.Listen(ctx)
	})

	// Update processor
	eg.Go(func() error {
		return flowFi.SendUpdates(ctx)
	})

	l := flowFi.Logger
	// Wait for all goroutines to complete
	l.Info("Bot is running... Press Ctrl+C to exit.")
	if err := eg.Wait(); err != nil {
		if err != context.Canceled {
			l.Warn("Shutting down with error", zap.Error(err))
		}
	} else {
		l.Info("Bot shut down gracefully.")
	}
	// Save subscriptions on shutdown
	if err := flowFi.StoreProgress(); err != nil {
		flowFi.Logger.Warn("Error saving subscriptions", zap.Error(err))
	} else {
		flowFi.Logger.Info("Subscriptions saved successfully.")
	}
}
