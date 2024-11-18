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

			for _, pair := range pairs {

				image := "https://fcljjsnuzjacwqgiqiib.supabase.co/storage/v1/object/public/token_images/images/f7c8d943-74ff-4cd4-acee-3e5dfcad9636.jpeg"
				// Create a new photo message with the file
				photo := tgbotapi.NewPhoto(0, tgbotapi.FileURL(image))

				l := flowFi.Logger.With(zap.String("pair", pair))

				l.Debug("process pair")
				data := subscriptions.GetSubscriptionData(pair)

				trades, lastProgressed := flowFi.GetTrades(ctx, pair, data.BlockNumber)
				l = l.With(zap.Any("lastProgressed", lastProgressed), zap.Int("trades", len(trades)))

				l.Debug("got trades")

				// if we are just starting to listen to this we do not process existing trades
				if data.BlockNumber == 0 {
					subscriptions.SetLastProgressed(pair, lastProgressed)
					return nil
				}

				for _, trade := range trades {

					// TODO: maybe send trade id along as well?
					msg := trade.String(pair)
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
