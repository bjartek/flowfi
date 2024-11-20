package bot

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

// takeScreenshot captures the screenshot of the provided URL and stores it in memory.
func (flowFi *FlowFi) Screenshot(ctx context.Context, pool string) ([]byte, error) {
	url := fmt.Sprintf(flowFi.ScreenshotUrlPattern, pool)
	fmt.Println(url)

	ctx, cancel := chromedp.NewContext(ctx,
		chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	// Create a buffer to store the screenshot
	var buf []byte
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(10*time.Second),
		chromedp.FullScreenshot(&buf, 90), // Capture full page screenshot
	)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
