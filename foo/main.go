package main

import (
	"fmt"

	. "flowfi_tg_bot/pkg/bot"
)

func main() {
	// Test cases
	fmt.Println(formatAmount("1500.0")) // Expected: "1500"
	fmt.Println(formatAmount("1000.5")) // Expected: "1000"
	fmt.Println(formatAmount("15.5"))   // Expected: "15.5"
	fmt.Println(formatAmount("0.0056")) // Expected: "0.0056"
	fmt.Println(formatAmount("0.1234")) // Expected: "0.1"
	fmt.Println(formatAmount("100000")) // Expected: "100000"
}
