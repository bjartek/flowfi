package bot

import (
	"fmt"
	"html/template"
	"math"
	"strconv"
	"strings"
)

const telegramTemplate = `
{{ .Emoticon }} 
{{ .PoolAttributes.Name }} Buy\!
🔀 Spent ${{ formatAmount .Attributes.VolumeInUsd }} \({{ formatAmount .Attributes.FromTokenAmount }} Flow\)
🔀 Got {{ formatAmount .Attributes.ToTokenAmount }} {{ .PoolAttributes.Symbol }}
👤 [Buyer](https://evm.flowscan.io/address/{{ .Attributes.TxFromAddress }}) / [TX](https://evm.flowscan.io/tx/{{ .Attributes.TxHash }})
💰 FDV ${{ formatAmount .TokenAttributes.FdvUsd }} 

🛒 [Buy](https://swap.kittypunch.xyz/?tokens={{ .Attributes.FromTokenAddress }}-{{ .Attributes.ToTokenAddress }}) 
📊 [Gecko](https://www.geckoterminal.com/flow-evm/pools/{{ .Pool }}) \| [Dexscreener](https://dexscreener.com/flowevm/{{ .Pool }})
`

func (flowFi *FlowFi) FormatTelegram(pool string, a Attributes, ta PoolAttributes, emoticon string, token *TokenAttributes) (string, error) {
	// Register functions for the template
	funcMap := template.FuncMap{
		"formatAmount":  FormatAmount,
		"formatAddress": formatAddress,
	}

	// Parse the template with the FuncMap
	tmpl, err := template.New("telegram").Funcs(funcMap).Parse(telegramTemplate)
	if err != nil {
		return "", err
	}

	value, err := strconv.ParseFloat(a.FromTokenAmount, 64)
	if err != nil {
		return "", err
	}

	// Use the existing structs directly
	data := struct {
		TokenAttributes TokenAttributes
		Pool            string
		Emoticon        string
		Attributes      Attributes
		PoolAttributes  PoolAttributes
	}{
		Emoticon:        repeatByCustomFactor(emoticon, value, flowFi.Config.EmoticonStep),
		Pool:            pool,
		Attributes:      a,
		PoolAttributes:  ta,
		TokenAttributes: *token,
	}
	// Render the template
	var output strings.Builder
	if err := tmpl.Execute(&output, data); err != nil {
		return "", err
	}

	return output.String(), nil
}

func FormatAmount(input string) string {
	// Parse the string as a float
	value, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return "Invalid number"
	}

	// Calculate the absolute value once
	absValue := math.Abs(value)

	// Dynamically determine precision based on magnitude
	var result string
	switch {
	case absValue >= 1_000_000:
		// For millions, use "m" and format with 1 decimal place
		result = fmt.Sprintf("%.3fm", value/1_000_000)
	case absValue >= 1_000:
		// For thousands, use "k" and format with 1 decimal place
		result = fmt.Sprintf("%.3fk", value/1_000)
	case absValue < 0.01:
		// For very small numbers, retain up to 6 decimal places
		result = strconv.FormatFloat(value, 'f', 6, 64)
	case absValue < 1:
		// For numbers < 1 but >= 0.01, keep up to 4 decimal places
		result = strconv.FormatFloat(value, 'f', 4, 64)
	default:
		// For numbers >= 1 and < 1000, retain up to 2 decimal places
		result = strconv.FormatFloat(value, 'f', 1, 64)
	}

	// Replace '.' with '\.'
	result = strings.ReplaceAll(result, ".", "\\.")

	return result
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

func repeatByCustomFactor(input string, value float64, factor int) string {
	var result strings.Builder
	bucketSize := getBucketForValue(value)
	result.WriteString(strings.Repeat(input, bucketSize*factor))
	return result.String()
}

func getBucketForValue(number float64) int {
	value := uint64(number)

	// If the value is 0, it should fall into bucket 1 (special case)
	if value == 0 {
		return 1
	}

	// Initialize bucket index to 1
	bucket := 1

	// Perform bit shifting to determine which power of 2 the number fits into
	for value > 0 {
		// Right shift by 1 bit
		value >>= 1
		bucket++
	}

	return bucket
}
