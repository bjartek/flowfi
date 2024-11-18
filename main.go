package main

import (
	"context"
	"flowfi_tg_bot/pkg/bot"
	"log"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func main() {
	flowFi := bot.NewBot()

	// Create a context that listens for interrupt signals
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	flowFi.Listen(ctx)

	eg, ctx := errgroup.WithContext(ctx)

	// Command listener
	eg.Go(func() error {
		return flowFi.Listen(ctx)
	})

	// Update processor
	eg.Go(func() error {
		return processUpdates(ctx, bot, subscriptions)
	})

	// Wait for all goroutines to complete
	log.Println("Bot is running... Press Ctrl+C to exit.")
	if err := eg.Wait(); err != nil {
		log.Printf("Shutting down with error: %v", err)
	} else {
		log.Println("Bot shut down gracefully.")
	}
}
