// Package historical — cohort analysis on a panel of (ticker, date, features,
// forward returns) rows. Buckets rows by RSI band, MA trend, MACD signal,
// Wyckoff stage, regime, and joint (trend×RSI), then computes mean forward
// return + win rate per bucket.
//
// Use this to derive empirical edges from your own data: pass in a slice of
// FeatureRow and get back a Report. The report is data-driven, no AI involved.
//
// Designed to be DB-agnostic — caller queries their own table and converts to
// FeatureRow. lotusmarket does not impose a schema.
package historical

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// FeatureRow is one row of cohort input — typically one (ticker, date).
// fwd_return values are interpreted as percent points (1.0 = 1%).
type FeatureRow struct {
	Ticker  string
	Date    time.Time
	Close   float64
	RSI14   *float64
	MA20    *float64
	MA50    *float64
	MA200   *float64
	MACD    *float64
	MACDSig *float64
	Wyckoff *int // 1-4
	Regime  *string
	Fwd5    *float64 // percent points
	Fwd20   *float64
	Fwd60   *float64
}

// Stats summarises one cohort.
type Stats struct {
	Bucket  string
	N       int
	Mean5d  float64
	Mean20d float64
	Mean60d float64
	Win5d   float64 // fraction 0-1
	Win20d  float64
	Win60d  float64
	Edge60d float64 // mean60d − baseline60d
}

// Report is the full cohort result.
type Report struct {
	Window    string
	CleanRows int
	Baseline  Stats
	Groups    map[string][]Stats // group name → sorted cohorts
}

// clipFwd rejects implausible fwd returns from unadjusted corporate actions.
// Bounds are: [-80%, +cap%]. Use 30 for 5d, 60 for 20d, 150 for 60d.
func clipFwd(x, cap float64) (float64, bool) {
	if x > cap || x < -80 {
		return 0, false
	}
	return x, true
}

type cohort struct {
	fwd5  []float64
	fwd20 []float64
	fwd60 []float64
}

func (c *cohort) add(r FeatureRow) {
	if r.Fwd5 != nil {
		if v, ok := clipFwd(*r.Fwd5, 30); ok {
			c.fwd5 = append(c.fwd5, v)
		}
	}
	if r.Fwd20 != nil {
		if v, ok := clipFwd(*r.Fwd20, 60); ok {
			c.fwd20 = append(c.fwd20, v)
		}
	}
	if r.Fwd60 != nil {
		if v, ok := clipFwd(*r.Fwd60, 150); ok {
			c.fwd60 = append(c.fwd60, v)
		}
	}
}

func summary(xs []float64) (mean, win float64, n int) {
	n = len(xs)
	if n == 0 {
		return
	}
	wins, sum := 0, 0.0
	for _, x := range xs {
		sum += x
		if x > 0 {
			wins++
		}
	}
	mean = sum / float64(n)
	win = float64(wins) / float64(n)
	return
}

// Analyze runs cohort analysis. Pass the data window string for the report
// header (e.g., "2017-01-01 → present").
func Analyze(rows []FeatureRow, window string) Report {
	cohorts := map[string]map[string]*cohort{
		"rsi": {}, "trend": {}, "macd": {}, "wyckoff": {}, "regime": {}, "joint": {},
	}
	get := func(g, lbl string) *cohort {
		if c, ok := cohorts[g][lbl]; ok {
			return c
		}
		c := &cohort{}
		cohorts[g][lbl] = c
		return c
	}

	for _, r := range rows {
		if r.RSI14 != nil {
			get("rsi", RSIBucket(*r.RSI14)).add(r)
		}
		trend := MATrend(r.Close, r.MA50, r.MA200)
		if trend != "" {
			get("trend", trend).add(r)
		}
		if r.MACD != nil && r.MACDSig != nil {
			lbl := "macd_bear"
			if *r.MACD > *r.MACDSig {
				lbl = "macd_bull"
			}
			get("macd", lbl).add(r)
		}
		if r.Wyckoff != nil {
			get("wyckoff", WyckoffLabel(*r.Wyckoff)).add(r)
		}
		if r.Regime != nil && *r.Regime != "" {
			get("regime", *r.Regime).add(r)
		}
		if r.RSI14 != nil && trend != "" {
			get("joint", trend+" × "+RSIBucket(*r.RSI14)).add(r)
		}
	}

	// Baseline
	var all5, all20, all60 []float64
	for _, r := range rows {
		if r.Fwd5 != nil {
			if v, ok := clipFwd(*r.Fwd5, 30); ok {
				all5 = append(all5, v)
			}
		}
		if r.Fwd20 != nil {
			if v, ok := clipFwd(*r.Fwd20, 60); ok {
				all20 = append(all20, v)
			}
		}
		if r.Fwd60 != nil {
			if v, ok := clipFwd(*r.Fwd60, 150); ok {
				all60 = append(all60, v)
			}
		}
	}
	m5, w5, _ := summary(all5)
	m20, w20, _ := summary(all20)
	m60, w60, _ := summary(all60)
	baseline := Stats{
		Bucket: "BASELINE",
		N:      len(all60),
		Mean5d: m5, Win5d: w5,
		Mean20d: m20, Win20d: w20,
		Mean60d: m60, Win60d: w60,
	}

	report := Report{
		Window:    window,
		CleanRows: len(rows),
		Baseline:  baseline,
		Groups:    map[string][]Stats{},
	}
	for g, m := range cohorts {
		stats := make([]Stats, 0, len(m))
		labels := make([]string, 0, len(m))
		for k := range m {
			labels = append(labels, k)
		}
		sort.Strings(labels)
		for _, lbl := range labels {
			c := m[lbl]
			mm5, ww5, n5 := summary(c.fwd5)
			mm20, ww20, _ := summary(c.fwd20)
			mm60, ww60, _ := summary(c.fwd60)
			stats = append(stats, Stats{
				Bucket: lbl, N: n5,
				Mean5d: mm5, Win5d: ww5,
				Mean20d: mm20, Win20d: ww20,
				Mean60d: mm60, Win60d: ww60,
				Edge60d: mm60 - m60,
			})
		}
		report.Groups[g] = stats
	}
	return report
}

// Markdown renders a Report as a markdown leaderboard.
func (r Report) Markdown() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "# Historical Cohort Analysis\n\n")
	fmt.Fprintf(&sb, "Generated: %s  \nWindow: %s  \nClean rows: %d (excludes outliers)\n\n",
		time.Now().Format("2006-01-02 15:04"), r.Window, r.CleanRows)
	fmt.Fprintf(&sb, "## Baseline\n\n| Horizon | Mean fwd return | Win rate | N |\n|---|---|---|---|\n")
	fmt.Fprintf(&sb, "| 5d | %+.2f%% | %.0f%% | %d |\n", r.Baseline.Mean5d, r.Baseline.Win5d*100, r.Baseline.N)
	fmt.Fprintf(&sb, "| 20d | %+.2f%% | %.0f%% | %d |\n", r.Baseline.Mean20d, r.Baseline.Win20d*100, r.Baseline.N)
	fmt.Fprintf(&sb, "| 60d | %+.2f%% | %.0f%% | %d |\n\n", r.Baseline.Mean60d, r.Baseline.Win60d*100, r.Baseline.N)

	order := []string{"rsi", "trend", "macd", "wyckoff", "regime", "joint"}
	titles := map[string]string{
		"rsi": "## RSI bucket", "trend": "## MA trend",
		"macd": "## MACD signal", "wyckoff": "## Wyckoff stage",
		"regime": "## Market regime", "joint": "## Joint: trend × RSI",
	}
	for _, g := range order {
		stats, ok := r.Groups[g]
		if !ok || len(stats) == 0 {
			continue
		}
		fmt.Fprintf(&sb, "%s\n\n| Bucket | N | Mean 5d | Win 5d | Mean 20d | Win 20d | **Mean 60d** | Win 60d | Edge 60d |\n", titles[g])
		sb.WriteString("|---|---|---|---|---|---|---|---|---|\n")
		for _, s := range stats {
			fmt.Fprintf(&sb, "| %s | %d | %+.2f%% | %.0f%% | %+.2f%% | %.0f%% | **%+.2f%%** | %.0f%% | %+.2f%% |\n",
				s.Bucket, s.N, s.Mean5d, s.Win5d*100, s.Mean20d, s.Win20d*100,
				s.Mean60d, s.Win60d*100, s.Edge60d)
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

// Helpers — exported so callers can apply same bucketing logic in their own
// pipelines without re-implementing.

func RSIBucket(rsi float64) string {
	switch {
	case rsi < 30:
		return "RSI<30 (oversold)"
	case rsi < 50:
		return "RSI 30-50"
	case rsi < 70:
		return "RSI 50-70"
	default:
		return "RSI>70 (overbought)"
	}
}

func MATrend(close float64, ma50, ma200 *float64) string {
	if ma200 == nil || ma50 == nil {
		return ""
	}
	above := close >= *ma200
	ma50Above := *ma50 >= *ma200
	if above && ma50Above {
		return "uptrend"
	}
	if !above && !ma50Above {
		return "downtrend"
	}
	return "mixed"
}

func WyckoffLabel(s int) string {
	switch s {
	case 1:
		return "1-accumulation"
	case 2:
		return "2-markup"
	case 3:
		return "3-distribution"
	case 4:
		return "4-decline"
	}
	return fmt.Sprintf("%d", s)
}
