package bot

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func (flowFi *FlowFi) GetToken(ctx context.Context, token string) *TokenAttributes {
	//*  'https://api.geckoterminal.com/api/v2/networks/flow-evm/tokens/0x995258cea49c25595cd94407fad9e99b81406a84' \

	url := fmt.Sprintf("%s/tokens/%s", flowFi.BaseUrl, token)
	l := flowFi.Logger.With(zap.String("token", token), zap.String("url", url))

	httpLogger := l.WithOptions(zap.IncreaseLevel(zapcore.InfoLevel))
	data, err := HttpGet[Token](ctx, url, httpLogger)
	if err != nil {
		flowFi.Logger.Warn("Failed getting token", zap.Error(err))
		return nil
	}

	return &data.Data.Attributes
}

type Token struct {
	Data TokenData `json:"data,omitempty"`
}

type TokenAttributes struct {
	TotalSupply       string `json:"total_supply,omitempty"`
	PriceUsd          string `json:"price_usd,omitempty"`
	FdvUsd            string `json:"fdv_usd,omitempty"`
	TotalReserveInUsd string `json:"total_reserve_in_usd,omitempty"`
}

type TokenData struct {
	ID         string          `json:"id,omitempty"`
	Type       string          `json:"type,omitempty"`
	Attributes TokenAttributes `json:"attributes,omitempty"`
}
