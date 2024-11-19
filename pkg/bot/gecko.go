package bot

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
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

func (a Attributes) String(pool string) string {
	token := "Avocado"
	return fmt.Sprintf(`
%s Buy\!
 🥑🥑🥑🥑🥑🥑🥑🥑🥑🥑🥑🥑🥑🥑🥑🥑🥑🥑🥑🥑🥑🥑🥑🥑🥑 🥑🥑🥑🥑🥑🥑🥑

🔀 Spent $%s \(%s Flow\)
🔀 Got %s %s
👤 [Buyer](https://evm.flowscan.io/address/%s) / [TX](https://evm.flowscan.io/tx/%s)
[Buy](https://swap.kittypunch.xyz/?tokens=%s-%s) \| [Gecko](https://www.geckoterminal.com/flow-evm/pools/%s) | [Dexscreener](https://dexscreener.com/flowevm/%s)

    `,
		token,
		formatAmount(a.VolumeInUsd), formatAmount(a.FromTokenAmount),
		formatAmount(a.ToTokenAmount), token,
		a.TxFromAddress, a.TxHash,
		a.FromTokenAddress, a.ToTokenAddress,
		pool, pool,
	)
}

func formatAddress(input string) string {
	// Check if the string length is at least 10 (6 + 4 minimum)
	if len(input) <= 10 {
		return input // No shortening needed
	}

	// Extract the first 6 and last 4 characters
	prefix := input[:6]
	suffix := input[len(input)-4:]

	// Combine them with "..." in between
	return fmt.Sprintf("%s\\.\\.\\.%s", prefix, suffix)
}

func formatAmount(input string) string {
	// Parse the string as a float
	value, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return "Invalid number"
	}

	// Handle exact zero
	if value == 0 {
		return "0"
	}

	// Dynamically determine precision based on magnitude
	var precision int
	if math.Abs(value) < 0.01 {
		// For very small numbers, retain up to 6 decimal places
		precision = 6
	} else if math.Abs(value) < 1 {
		// For numbers < 1 but >= 0.01, keep up to 4 decimal places
		precision = 4
	} else {
		// For numbers >= 1, retain up to 2 decimal places
		precision = 1
	}

	// Format the value with the determined precision
	result := strconv.FormatFloat(value, 'f', precision, 64)

	// Trim trailing zeros and unnecessary decimal point
	result = strings.TrimRight(result, "0")
	result = strings.TrimRight(result, ".")

	result = strings.ReplaceAll(result, ".", "\\.")

	return result
}

func EscapeMarkdownV2(text string) string {
	specialChars := []string{
		"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!",
	}

	// Escape each special character by prefixing it with a backslash
	for _, char := range specialChars {
		text = strings.ReplaceAll(text, char, "\\"+char)
	}

	return text
}
