package nlu

import (
	"regexp"
	"strings"

	"github.com/ducnhd/lotusmarket/go/types"
)

type Result struct {
	Intent     string
	Confidence float64
	Tickers    []string
	Timeframe  string
	Query      string
}

type Parser struct {
	intentKeywords map[string][]string
	tickerSet      map[string]bool
}

func New(customTickers []string) *Parser {
	tickerSet := make(map[string]bool)
	tickers := customTickers
	if tickers == nil {
		tickers = types.VN30
	}
	for _, t := range tickers {
		tickerSet[t] = true
	}
	return &Parser{
		intentKeywords: map[string][]string{
			"GET_STOCK_INFO":     {"giá", "bao nhiêu", "hiện tại", "price", "current", "info", "thông tin"},
			"ANALYZE_TREND":      {"xu hướng", "tăng", "giảm", "chart", "trend", "biểu đồ", "phân tích"},
			"GET_SENTIMENT":      {"tin tức", "ý kiến", "news", "sentiment", "dư luận", "thị trường"},
			"GET_RECOMMENDATION": {"nên", "mua", "bán", "buy", "sell", "signal", "khuyến nghị", "đầu tư"},
			"COMPARE_STOCKS":     {"so sánh", "so với", "compare", "vs", "và"},
		},
		tickerSet: tickerSet,
	}
}

var vnStopWords = map[string]bool{
	"CAC": true, "CON": true, "CUA": true, "CHO": true, "DAU": true,
	"DAN": true, "GIA": true, "HAI": true, "HOA": true, "HOI": true,
	"LAM": true, "MOI": true, "NAM": true, "NEN": true, "SAU": true,
	"TAI": true, "THE": true, "TIN": true, "VAN": true, "VOI": true,
}

var tickerPattern = regexp.MustCompile(`\b([A-Z]{3})\b`)

func (p *Parser) Parse(query string) Result {
	lower := strings.ToLower(query)
	result := Result{Query: query, Tickers: p.extractTickers(query)}
	bestIntent := "GET_STOCK_INFO"
	bestScore := 0.0
	for intent, keywords := range p.intentKeywords {
		score := 0.0
		for _, kw := range keywords {
			if strings.Contains(lower, kw) {
				score += 1.0
			}
		}
		if score > bestScore {
			bestScore = score
			bestIntent = intent
		}
	}
	result.Intent = bestIntent
	result.Confidence = bestScore / 3.0
	if result.Confidence > 1.0 {
		result.Confidence = 1.0
	}
	result.Timeframe = extractTimeframe(lower)
	return result
}

func (p *Parser) extractTickers(query string) []string {
	matches := tickerPattern.FindAllString(strings.ToUpper(query), -1)
	var tickers []string
	seen := make(map[string]bool)
	for _, m := range matches {
		if vnStopWords[m] {
			continue
		}
		if p.tickerSet[m] && !seen[m] {
			tickers = append(tickers, m)
			seen[m] = true
		}
	}
	return tickers
}

func extractTimeframe(lower string) string {
	if strings.Contains(lower, "tuần") || strings.Contains(lower, "week") {
		return "1w"
	}
	if strings.Contains(lower, "tháng") || strings.Contains(lower, "month") {
		return "1m"
	}
	if strings.Contains(lower, "năm") || strings.Contains(lower, "year") {
		return "1y"
	}
	return "1d"
}
