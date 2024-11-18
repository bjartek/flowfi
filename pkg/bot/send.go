package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func (flowFi *FlowFi) Send(l *zap.Logger, msg tgbotapi.MessageConfig) {
	msg.DisableWebPagePreview = true
	msg.ParseMode = "markdownv2"
	_, err := flowFi.Tgbot.Send(msg)
	if err != nil {
		l.Warn("failed sending", zap.Error(err))
	}
}
