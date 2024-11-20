package bot

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

// Process updates for each unique pair
func (flowFi *FlowFi) SendUpdates(ctx context.Context) error {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	subscriptions := flowFi.Subscriptions
	for {
		select {
		case <-ticker.C:
			flowFi.Logger.Debug("tick")
			// Fetch all unique pairs
			pairs := subscriptions.GetPairs()

			// screenshot
			for _, pair := range pairs {

				l := flowFi.Logger.With(zap.String("pair", pair))

				l.Debug("process pair")
				data := subscriptions.GetSubscriptionData(pair)

				screenshot, err := flowFi.Screenshot(ctx, pair)
				if err != nil {
					return err
				}

				tbytes := tgbotapi.FileBytes{
					Name:  "trend.png",
					Bytes: screenshot,
				}

				//				image := data.TokenAttributes.ImageURL
				// Create a new photo message with the file
				photo := tgbotapi.NewPhoto(0, tbytes)

				trades, lastProgressed := flowFi.GetTrades(ctx, pair, data.BlockNumber)
				l = l.With(zap.Any("lastProgressed", lastProgressed), zap.Int("trades", len(trades)))

				l.Debug("got trades")

				// if we are just starting to listen to this we do not process existing trades
				if data.BlockNumber == 0 {
					subscriptions.SetLastProgressed(pair, lastProgressed)
					return nil
				}

				for _, trade := range trades {

					msg, err := flowFi.FormatTelegram(pair, trade, *data.TokenAttributes, data.Emoticon)
					if err != nil {
						l.Warn("failed formating", zap.Error(err))
					}
					photo.Caption = msg
					photo.ParseMode = "MarkdownV2"
					for _, chatID := range data.ChatIDs {
						//		l2 := l.With(zap.Int64("chatId", chatID))
						photo.ChatID = chatID
						flowFi.Tgbot.Send(photo)
					}
				}
				subscriptions.SetLastProgressed(pair, lastProgressed)
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
