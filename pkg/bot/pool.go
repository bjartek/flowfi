package bot

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Pool struct {
	Data []TokenData `json:"data,omitempty"`
}
type TokenAttributes struct {
	CoingeckoCoinID any     `json:"coingecko_coin_id,omitempty"`
	DiscordURL      any     `json:"discord_url,omitempty"`
	TwitterHandle   any     `json:"twitter_handle,omitempty"`
	Address         string  `json:"address,omitempty"`
	Name            string  `json:"name,omitempty"`
	Symbol          string  `json:"symbol,omitempty"`
	ImageURL        string  `json:"image_url,omitempty"`
	Description     string  `json:"description,omitempty"`
	TelegramHandle  string  `json:"telegram_handle,omitempty"`
	Websites        []any   `json:"websites,omitempty"`
	Decimals        int     `json:"decimals,omitempty"`
	GtScore         float64 `json:"gt_score,omitempty"`
}
type TokenData struct {
	ID         string          `json:"id,omitempty"`
	Type       string          `json:"type,omitempty"`
	Attributes TokenAttributes `json:"attributes,omitempty"`
}

func (flowFi *FlowFi) GetPoolInformation(ctx context.Context, logger *zap.Logger, pool string) *TokenAttributes {
	httpLogger := logger.WithOptions(zap.IncreaseLevel(zapcore.InfoLevel))
	url := fmt.Sprintf("%s/%s/info", flowFi.BaseUrl, pool)
	poolData, err := HttpGet[Pool](ctx, url, httpLogger)
	if err != nil {
		flowFi.Logger.Warn("Failed getting pool", zap.Error(err))
	}

	for _, p := range poolData.Data {
		if p.Attributes.Symbol != "WFLOW" {
			return &p.Attributes
		}
	}
	return nil
}
