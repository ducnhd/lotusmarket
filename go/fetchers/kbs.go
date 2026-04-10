package fetchers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ducnhd/lotusmarket/go/types"
)

const kbsBaseURL = "https://kbbuddywts.kbsec.com.vn/sas/kbsv-stock-data-store/stock"

func KBS(ctx context.Context, ticker string) (*types.KBSQuote, error) {
	to := time.Now().Format("2006-01-02")
	from := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	url := fmt.Sprintf("%s/%s/historical-quotes?from=%s&to=%s", kbsBaseURL, ticker, from, to)
	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("KBS request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("KBS HTTP %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("KBS read body: %w", err)
	}
	var quotes []types.KBSQuote
	if err := json.Unmarshal(body, &quotes); err != nil {
		return nil, fmt.Errorf("KBS parse: %w", err)
	}
	if len(quotes) == 0 {
		return nil, fmt.Errorf("no KBS data for %s", ticker)
	}
	return &quotes[0], nil
}
