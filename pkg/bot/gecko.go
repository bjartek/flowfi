package bot

import (
	"context"
	"fmt"
	"time"

	"github.com/samber/lo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func (flowFi *FlowFi) GetTrades(ctx context.Context, pool string, lastRead uint64) ([]Attributes, uint64) {
	url := fmt.Sprintf("%s/%s/trades", flowFi.BaseUrl, pool)
	l := flowFi.Logger.With(zap.String("pool", pool), zap.String("url", url), zap.Any("lastRead", lastRead))

	l.Info("getting trades")
	// TODO; configure this
	httpLogger := l.WithOptions(zap.IncreaseLevel(zapcore.InfoLevel))

	trades, err := HttpGet[Trades](ctx, url, httpLogger)
	if err != nil {
		flowFi.Logger.Warn("Failed getting trades", zap.Error(err))
	}

	attr := []Attributes{}
	lastProgressed := uint64(0)
	if len(trades.Data) > 0 {
		lastProgressed = trades.Data[0].Attributes.BlockNumber
	}
	l = l.With(zap.Any("lastProgressed", lastProgressed))

	for _, d := range trades.Data {
		// we are not interested in sells

		if d.Attributes.Kind == "sell" {
			l.Debug("skipping sell")
			continue
		}
		// if we have read this before we return the items progressed reversed
		if d.Attributes.BlockNumber <= lastRead {
			l.Debug("block is before last read")
			return lo.Reverse(attr), lastProgressed
		}
		l.Debug("appending")
		attr = append(attr, d.Attributes)
	}
	l.Info("We might have missed trades since we added them all")
	return lo.Reverse(attr), lastProgressed
}

type Trades struct {
	Data []Data `json:"data"`
}
type Attributes struct {
	BlockTimestamp           time.Time `json:"block_timestamp"`
	FromTokenAmount          string    `json:"from_token_amount"`
	TxFromAddress            string    `json:"tx_from_address"`
	ToTokenAmount            string    `json:"to_token_amount"`
	PriceFromInCurrencyToken string    `json:"price_from_in_currency_token"`
	PriceToInCurrencyToken   string    `json:"price_to_in_currency_token"`
	PriceFromInUsd           string    `json:"price_from_in_usd"`
	PriceToInUsd             string    `json:"price_to_in_usd"`
	TxHash                   string    `json:"tx_hash"`
	Kind                     string    `json:"kind"`
	VolumeInUsd              string    `json:"volume_in_usd"`
	FromTokenAddress         string    `json:"from_token_address"`
	ToTokenAddress           string    `json:"to_token_address"`
	BlockNumber              uint64    `json:"block_number"`
}
type Data struct {
	ID         string     `json:"id"`
	Type       string     `json:"type"`
	Attributes Attributes `json:"attributes"`
}
