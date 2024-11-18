package bot

import (
	"context"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Process updates for each unique pair
func (flowFi *FlowFi) SendUpdates(ctx context.Context) error {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	subscriptions := flowFi.Subscriptions
	for {
		select {
		case <-ticker.C:
			// Fetch all unique pairs
			pairs := subscriptions.GetPairs()

			for _, pair := range pairs {

				data := subscriptions.GetSubscriptionData(pair)

				trades, lastProgressed := flowFi.GetTrades(ctx, pair, data.blockNumber)

				for _, chatID := range data.chatIDs {

					msg := tgbotapi.NewMessage(chatID, update)
					_, err := bot.Send(msg)
					if err != nil {
						log.Printf("Failed to send message to chat %d: %v", chatID, err)
					}
				}
			}

		case <-ctx.Done():
			log.Println("Stopping update processor...")
			return ctx.Err()
		}
	}
}
