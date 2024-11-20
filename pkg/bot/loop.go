package bot

import (
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func (flowFi *FlowFi) Listen(ctx context.Context) error {
	updates := flowFi.Tgbot.GetUpdatesChan(flowFi.UpdateConfig)
	subscriptions := flowFi.Subscriptions

	for {
		select {
		case update := <-updates:
			if update.Message != nil && update.Message.IsCommand() {
				chatID := update.Message.Chat.ID
				command := update.Message.Command()
				args := update.Message.CommandArguments()

				l := flowFi.Logger.With(zap.String("cmd", command), zap.Int64("chatID", chatID), zap.String("args", args))
				switch command {
				case "subscribe":
					if args != "" {

						argList := strings.Split(args, " ")

						emoticon := "💰"
						if len(argList) == 2 {
							emoticon = argList[1]
						}
						ti := flowFi.GetPoolInformation(ctx, l, argList[0])
						subscriptions.AddSubscription(chatID, args, ti, emoticon)

						l.Info("Subscribed")
						msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Subscribed to %s", args))
						flowFi.Send(l, msg)
					} else {
						msg := tgbotapi.NewMessage(chatID, "Please specify a pair to subscribe to")
						flowFi.Send(l, msg)
					}
				case "unsubscribe":
					if args != "" {

						l.Info("Unsubscribed")
						subscriptions.RemoveSubscription(chatID, args)
						msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Unsubscribed from %s", args))

						flowFi.Send(l, msg)
					} else {
						msg := tgbotapi.NewMessage(chatID, "Please specify a pair to unsubscribe from")
						flowFi.Send(l, msg)
					}
				case "status":
					var subscribedPairs []string
					subscriptions.mu.RLock()
					for pair, data := range subscriptions.pairs {
						for _, id := range data.ChatIDs {
							if id == chatID {
								subscribedPairs = append(subscribedPairs, pair)
								break
							}
						}
					}
					subscriptions.mu.RUnlock()

					if len(subscribedPairs) > 0 {
						msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("You are subscribed to the following pairs: %v", subscribedPairs))
						flowFi.Send(l, msg)
					} else {
						msg := tgbotapi.NewMessage(chatID, "You are not subscribed to any pairs")
						flowFi.Send(l, msg)
					}
				case "help":
					helpText := "⚙️   /status to see what pairs are monitored\n" +
						"⚙️   /subscribe to subscribe to a new pair\n" +
						"⚙️   /unsubscribe to unsubscribe from a new pair"
					msg := tgbotapi.NewMessage(chatID, helpText)
					flowFi.Send(l, msg)
				default:
					msg := tgbotapi.NewMessage(chatID, "Unknown command.")
					flowFi.Send(l, msg)
				}
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
