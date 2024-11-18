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
)

func (flowFi *FlowFi) GetTrades(ctx context.Context, pool string, lastRead uint64) ([]Attributes, uint64) {
	url := fmt.Sprintf("%s/%s/trades", flowFi.BaseUrl, pool)
	l := flowFi.Logger.With(zap.String("pool", pool), zap.String("url", url))

	l.Info("getting trades")
	trades, err := HttpGet[Trades](ctx, url, l)
	if err != nil {
		flowFi.Logger.Warn("Failed getting trades", zap.Error(err))
	}

	attr := []Attributes{}
	lastProgressed := uint64(0)
	for _, d := range trades.Data {
		// if we have read this before we return the items progressed reversed
		if d.Attributes.BlockNumber <= lastRead {
			return lo.Reverse(attr), d.Attributes.BlockNumber
		}
		lastProgressed = d.Attributes.BlockNumber
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

func (a Attributes) String() string {
	return fmt.Sprintf(`
🚀 *Buy\!*  
💵_%s_ 💰 *%s*  
🕒 %s
📈 %s$
[tx](https://evm.flowscan.io/tx/%s) 
    `,
		formatAmount(a.ToTokenAmount), formatAmount(a.PriceToInCurrencyToken),
		EscapeMarkdownV2(a.BlockTimestamp.String()),
		formatAmount(a.VolumeInUsd), a.TxHash)
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
		precision = 2
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
