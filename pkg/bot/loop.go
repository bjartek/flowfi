package bot

import (
	"context"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (flowFi *FlowFi) Listen(ctx context.Context) error {
	updates := flowFi.Tgbot.GetUpdatesChan(flowFi.UpdateConfig)
	subscriptions := flowFi.Subscriptions
	bot := flowFi.Tgbot

	for {
		select {
		case update := <-updates:
			if update.Message != nil && update.Message.IsCommand() {
				chatID := update.Message.Chat.ID
				command := update.Message.Command()
				args := update.Message.CommandArguments()

				switch command {
				case "subscribe":
					if args != "" {
						subscriptions.AddSubscription(chatID, args)
						msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Subscribed to %s", args))
						bot.Send(msg)
					} else {
						msg := tgbotapi.NewMessage(chatID, "Please specify a pair to subscribe to.")
						bot.Send(msg)
					}
				case "unsubscribe":
					if args != "" {
						subscriptions.RemoveSubscription(chatID, args)
						msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Unsubscribed from %s", args))
						bot.Send(msg)
					} else {
						msg := tgbotapi.NewMessage(chatID, "Please specify a pair to unsubscribe from.")
						bot.Send(msg)
					}
				case "status":
					var subscribedPairs []string
					subscriptions.mu.RLock()
					for pair, data := range subscriptions.pairs {
						for _, id := range data.chatIDs {
							if id == chatID {
								subscribedPairs = append(subscribedPairs, pair)
								break
							}
						}
					}
					subscriptions.mu.RUnlock()

					if len(subscribedPairs) > 0 {
						msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("You are subscribed to the following pairs: %v", subscribedPairs))
						bot.Send(msg)
					} else {
						msg := tgbotapi.NewMessage(chatID, "You are not subscribed to any pairs.")
						bot.Send(msg)
					}
				case "help":
					helpText := "⚙️   /status to see what pairs are monitored\n" +
						"⚙️   /subscribe to subscribe to a new pair\n" +
						"⚙️   /unsubscribe to unsubscribe from a new pair"
					msg := tgbotapi.NewMessage(chatID, helpText)
					bot.Send(msg)
				default:
					msg := tgbotapi.NewMessage(chatID, "Unknown command.")
					bot.Send(msg)
				}
			}

		case <-ctx.Done():
			log.Println("Stopping command listener...")
			return ctx.Err()
		}
	}
}
