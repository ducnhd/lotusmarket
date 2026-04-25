// Package earnings provides pure-logic extraction of profit/revenue numbers
// from Vietnamese news headlines, period detection, and BCTC deadline generation.
// No database coupling — all functions are stateless.
package earnings

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// Deadline is a BCTC (financial report) regulatory submission deadline.
type Deadline struct {
	Date        string // YYYY-MM-DD
	Description string // Vietnamese description
}

// skipTickers are uppercase acronyms that must not be treated as stock tickers.
var skipTickers = map[string]bool{
	"CEO": true, "CFO": true, "IPO": true, "GDP": true,
	"VND": true, "USD": true, "ETF": true, "EPS": true,
	"FDI": true, "FED": true, "THE": true, "AND": true,
}

// isSkipTicker returns true when s is a blocked acronym.
func isSkipTicker(s string) bool {
	return skipTickers[s]
}

// Compiled regexps (package-level, compiled once).
var (
	numberPat = `([\d.,]+)`
	qualOpt   = `(?:hơn|gần|khoảng|đạt|vượt)?\s*`

	profitPattern = regexp.MustCompile(
		`(?i)(?:lãi(?:\s+(?:ròng|sau\s+thuế|trước\s+thuế))?|lợi\s+nhuận|LNST)` +
			`[^.]{0,40}?` + qualOpt + numberPat + `\s*tỷ`,
	)

	revenuePattern = regexp.MustCompile(
		`(?i)doanh\s+thu[^.]{0,40}?` + qualOpt + numberPat + `\s*tỷ`,
	)

	targetPattern = regexp.MustCompile(
		`(?i)(?:kế\s+hoạch|mục\s+tiêu|chỉ\s+tiêu)[^.]{0,80}?` +
			`(?:lãi(?:\s+(?:ròng|sau\s+thuế|trước\s+thuế))?|lợi\s+nhuận)` +
			`[^.]{0,30}?` + qualOpt + numberPat + `\s*tỷ`,
	)

	tickerPattern = regexp.MustCompile(`\b([A-Z]{2,5})\b`)
	yearPattern   = regexp.MustCompile(`\b(20\d{2})\b`)

	periodPatterns = map[string]*regexp.Regexp{
		"Q1": regexp.MustCompile(`(?i)(?:quý\s*(?:1|I|một)|Q1)(?:\s*[/\-]\s*(\d{4}))?`),
		"Q2": regexp.MustCompile(`(?i)(?:quý\s*(?:2|II|hai)|Q2)(?:\s*[/\-]\s*(\d{4}))?`),
		"Q3": regexp.MustCompile(`(?i)(?:quý\s*(?:3|III|ba)|Q3)(?:\s*[/\-]\s*(\d{4}))?`),
		"Q4": regexp.MustCompile(`(?i)(?:quý\s*(?:4|IV|bốn)|Q4)(?:\s*[/\-]\s*(\d{4}))?`),
		"9M": regexp.MustCompile(`(?i)(?:9\s*tháng)(?:\s+(?:đầu\s+năm|năm)\s*(\d{4}))?`),
		"6M": regexp.MustCompile(`(?i)(?:6\s*tháng|nửa\s+đầu\s+năm|bán\s+niên)(?:\s*(?:năm)?\s*(\d{4}))?`),
		"FY": regexp.MustCompile(`(?i)(?:cả\s+năm|cả\s+năm\s+(\d{4})|năm\s+(\d{4}))`),
	}

	periodOrder = []string{"Q1", "Q2", "Q3", "Q4", "9M", "6M", "FY"}
)

// ParseVietnameseNumber parses a Vietnamese-formatted number string.
// "35.000" → 35000 (dots are thousands separators in Vietnamese).
func ParseVietnameseNumber(s string) int64 {
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.TrimSpace(s)
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return val
}

// DetectPeriod identifies the reporting period from a Vietnamese headline.
// Returns strings like "Q3/2025", "9M/2026", "FY/2025", "Q4" (no year if not found).
func DetectPeriod(title string) string {
	for _, key := range periodOrder {
		pat := periodPatterns[key]
		m := pat.FindStringSubmatch(title)
		if m == nil {
			continue
		}
		year := ""
		for i := 1; i < len(m); i++ {
			if m[i] != "" {
				year = m[i]
				break
			}
		}
		if year == "" {
			ym := yearPattern.FindStringSubmatch(title)
			if ym != nil {
				year = ym[1]
			}
		}
		if year != "" {
			return key + "/" + year
		}
		return key
	}
	return ""
}

// ExtractTickers finds 3-letter uppercase stock tickers in a Vietnamese headline.
// Filters out common non-ticker acronyms (CEO, GDP, etc.).
func ExtractTickers(title string) []string {
	var tickers []string
	seen := map[string]bool{}
	matches := tickerPattern.FindAllString(title, -1)
	for _, m := range matches {
		if len(m) == 3 && !isSkipTicker(m) && !seen[m] {
			valid := true
			for _, r := range m {
				if !unicode.IsUpper(r) || r > 'Z' {
					valid = false
					break
				}
			}
			if valid {
				tickers = append(tickers, m)
				seen[m] = true
			}
		}
	}
	return tickers
}

// ExtractProfit extracts net profit (tỷ VND) and the reporting period from a headline.
// Returns (profit, period). profit=0 if not found.
func ExtractProfit(title string) (int64, string) {
	m := profitPattern.FindStringSubmatch(title)
	if m == nil {
		return 0, ""
	}
	profit := ParseVietnameseNumber(m[len(m)-1])
	if profit == 0 {
		return 0, ""
	}
	return profit, DetectPeriod(title)
}

// ExtractRevenue extracts revenue (tỷ VND) from a Vietnamese headline.
// Returns 0 if not found.
func ExtractRevenue(title string) int64 {
	m := revenuePattern.FindStringSubmatch(title)
	if m == nil {
		return 0
	}
	return ParseVietnameseNumber(m[1])
}

// ExtractAnnualTarget extracts an annual profit target (tỷ VND) from ĐHCĐ-style headlines.
// Returns 0 if not found.
func ExtractAnnualTarget(title string) int64 {
	m := targetPattern.FindStringSubmatch(title)
	if m == nil {
		return 0
	}
	return ParseVietnameseNumber(m[1])
}

// BCTCDeadlines returns the 6 BCTC regulatory submission deadlines for the given year.
// Deadlines:
//   - Q1: April 20
//   - Q2: July 20
//   - H1 audited: August 14
//   - Q3: October 20
//   - Q4 (preliminary): January 20 next year
//   - FY audited: March 31 next year
func BCTCDeadlines(year int) []Deadline {
	return []Deadline{
		{fmt.Sprintf("%d-04-20", year), fmt.Sprintf("Hạn nộp BCTC Q1/%d", year)},
		{fmt.Sprintf("%d-07-20", year), fmt.Sprintf("Hạn nộp BCTC Q2/%d", year)},
		{fmt.Sprintf("%d-08-14", year), fmt.Sprintf("Hạn nộp BCTC bán niên kiểm toán %d", year)},
		{fmt.Sprintf("%d-10-20", year), fmt.Sprintf("Hạn nộp BCTC Q3/%d", year)},
		{fmt.Sprintf("%d-01-20", year+1), fmt.Sprintf("Hạn nộp BCTC Q4/%d", year)},
		{fmt.Sprintf("%d-03-31", year+1), fmt.Sprintf("Hạn nộp BCTC năm %d kiểm toán", year)},
	}
}
