package fetchers

import (
	"context"
	"log"

	"github.com/lotusmarket/lotusmarket-go/types"
)

func StockWithFallback(ctx context.Context, ticker string) (*types.StockData, error) {
	stock, err := VPS(ctx, ticker)
	if err != nil {
		log.Printf("[lotusmarket] VPS failed for %s, trying Entrade: %v", ticker, err)
		return EntradeLatest(ctx, ticker)
	}
	return stock, nil
}
