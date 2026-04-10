package fetchers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/lotusmarket/lotusmarket-go/types"
)

const vpsBaseURL = "https://bgapidatafeed.vps.com.vn/getliststockdata"

func VPS(ctx context.Context, ticker string) (*types.StockData, error) {
	stocks, err := VPSMultiple(ctx, []string{ticker})
	if err != nil {
		return nil, err
	}
	if len(stocks) == 0 {
		return nil, fmt.Errorf("no data for %s", ticker)
	}
	return &stocks[0], nil
}

func VPSMultiple(ctx context.Context, tickers []string) ([]types.StockData, error) {
	url := vpsBaseURL + "/" + strings.Join(tickers, ",")
	client := &http.Client{Timeout: 15 * time.Second}
	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(attempt) * time.Second)
		}
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("create request: %w", err)
		}
		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = err
			continue
		}
		return parseVPSResponse(body)
	}
	return nil, fmt.Errorf("VPS API failed after 3 attempts: %w", lastErr)
}

func parseVPSResponse(data []byte) ([]types.StockData, error) {
	var raw []map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	var stocks []types.StockData
	for _, item := range raw {
		ticker := getString(item, "sym")
		if ticker == "" {
			continue
		}
		refPrice := getFloat(item, "r")
		lastPrice := getFloat(item, "lastPrice")
		openPrice := getFloat(item, "openPrice")
		if lastPrice == 0 && refPrice > 0 {
			lastPrice = refPrice
		}
		var changeValue, changePc float64
		if refPrice > 0 {
			changeValue = lastPrice - refPrice
			changePc = (changeValue / refPrice) * 100
		}
		stocks = append(stocks, types.StockData{
			Ticker:         ticker,
			Date:           time.Now().Format("2006-01-02"),
			Open:           openPrice,
			High:           getFloat(item, "highPrice"),
			Low:            getFloat(item, "lowPrice"),
			Close:          lastPrice,
			Volume:         getInt(item, "lot"),
			ChangePercent:  changePc,
			ChangeValue:    changeValue,
			RefPrice:       refPrice,
			Ceiling:        getFloat(item, "c"),
			Floor:          getFloat(item, "f"),
			ForeignBuyVol:  getInt(item, "fBVol"),
			ForeignSellVol: getInt(item, "fSVolume"),
			ForeignNetVol:  getInt(item, "fBVol") - getInt(item, "fSVolume"),
			Bid:            parseFloatStr(getString(item, "g1")),
			Ask:            parseFloatStr(getString(item, "g2")),
			BidVol:         parsePipeVol(getString(item, "g1")),
			AskVol:         parsePipeVol(getString(item, "g2")),
		})
	}
	for i := range stocks {
		stocks[i].ToVND()
	}
	return stocks, nil
}
