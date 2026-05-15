// lmcli — lotusmarket command-line tool.
//
// One-binary entry to the library. Designed for automation: prints structured
// markdown / JSON to stdout, intended to be piped into a file, a Telegram bot,
// or committed to a GitHub Pages repo for free static hosting.
//
// Subcommands:
//
//	lmcli pulse                 → daily VN30 + global summary (markdown)
//	lmcli quote TICKER          → real-time quote (text)
//	lmcli rate TICKER           → 6-dim star ratings + verdict (markdown)
//	lmcli screen [--rsi=X-Y]    → screen VN30 by criteria
//	lmcli sectors               → sector flow leaderboard
//	lmcli global                → global indices snapshot
//	lmcli dividends TICKER      → corporate action calendar
//	lmcli report                → full daily report (combines all of above)
//
// All commands take no external state. No DB. No API key for the core flows
// (Yahoo + VPS + VCI are all key-free). FRED requires FRED_API_KEY env if you
// extend later.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/ducnhd/lotusmarket/go/fetchers"
	"github.com/ducnhd/lotusmarket/go/market"
	"github.com/ducnhd/lotusmarket/go/ratings"
	"github.com/ducnhd/lotusmarket/go/signals"
	"github.com/ducnhd/lotusmarket/go/technical"
	"github.com/ducnhd/lotusmarket/go/types"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(2)
	}
	cmd := os.Args[1]
	args := os.Args[2:]
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	switch cmd {
	case "pulse":
		runPulse(ctx)
	case "quote":
		runQuote(ctx, args)
	case "rate":
		runRate(ctx, args)
	case "screen":
		runScreen(ctx, args)
	case "sectors":
		runSectors(ctx)
	case "global":
		runGlobal(ctx)
	case "dividends":
		runDividends(ctx, args)
	case "report":
		runReport(ctx)
	case "-h", "--help", "help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", cmd)
		printUsage()
		os.Exit(2)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, `lmcli — lotusmarket CLI

Usage:
  lmcli pulse                Daily VN30 pulse (markdown)
  lmcli quote TICKER         Real-time quote
  lmcli rate TICKER          6-dim star ratings
  lmcli screen [flags]       Screen VN30 by criteria
  lmcli sectors              Sector flow leaderboard
  lmcli global               Global indices snapshot
  lmcli dividends TICKER     Upcoming corporate actions
  lmcli report [--out=path]  Full daily report

Examples:
  lmcli pulse > today.md
  lmcli rate ACB
  lmcli screen --rsi=30-50 --json
  lmcli report --out=reports/today.md`)
}

// VN30 default tracking universe — keep in sync with HOSE listing as of 2026.
var vn30 = []string{
	"ACB", "BCM", "BID", "BVH", "CTG", "FPT", "GAS", "GVR", "HDB", "HPG",
	"MBB", "MSN", "MWG", "PLX", "POW", "SAB", "SHB", "SSB", "SSI", "STB",
	"TCB", "TPB", "VCB", "VHM", "VIB", "VIC", "VJC", "VNM", "VPB", "VRE",
}

func runQuote(ctx context.Context, args []string) {
	if len(args) == 0 {
		log.Fatal("usage: lmcli quote TICKER")
	}
	ticker := strings.ToUpper(args[0])
	s, err := fetchers.StockWithFallback(ctx, ticker)
	if err != nil {
		log.Fatalf("fetch %s: %v", ticker, err)
	}
	fmt.Printf("%s: %.0f VND  Δ %+.2f%%  vol %d  foreign net %d\n",
		s.Ticker, s.Close, s.ChangePercent, s.Volume, s.ForeignNetVol)
}

func runRate(ctx context.Context, args []string) {
	if len(args) == 0 {
		log.Fatal("usage: lmcli rate TICKER")
	}
	ticker := strings.ToUpper(args[0])
	hist, err := fetchers.EntradeHistory(ctx, ticker, 250)
	if err != nil {
		log.Fatalf("history %s: %v", ticker, err)
	}
	closes := make([]float64, len(hist))
	vols := make([]int64, len(hist))
	for i, h := range hist {
		closes[i] = h.Close
		vols[i] = h.Volume
	}
	d := technical.Dashboard(closes)
	r := ratings.Compute(closes, vols, d.Score, d.RSI, d.MA20, d.MA50, d.MA200)

	fmt.Printf("# %s — Star Ratings\n\n", ticker)
	fmt.Printf("**Overall: %s (%d/100)** · RSI %.1f · Signal %s · Score %.0f\n\n",
		r.OverallVerdict, r.OverallGauge, d.RSI, d.Signal, d.Score)
	fmt.Printf("| Dimension | Stars (1-5) |\n|---|---|\n")
	fmt.Printf("| Price strength | %s |\n", stars(r.PriceStrength))
	fmt.Printf("| Trend strength | %s |\n", stars(r.TrendStrength))
	fmt.Printf("| Short-term position (RSI) | %s |\n", stars(r.ShortTermPos))
	fmt.Printf("| Money flow | %s |\n", stars(r.MoneyFlow))
	fmt.Printf("| Volatility | %s |\n", stars(r.VolatilityRating))
	fmt.Printf("| Base range | %s |\n", stars(r.BaseRange))
}

func stars(n int) string {
	if n < 1 {
		n = 1
	}
	if n > 5 {
		n = 5
	}
	return strings.Repeat("★", n) + strings.Repeat("☆", 5-n)
}

func runScreen(ctx context.Context, args []string) {
	fs := flag.NewFlagSet("screen", flag.ExitOnError)
	rsiRange := fs.String("rsi", "", "RSI range filter, e.g. 30-50 or 60-")
	signalFilter := fs.String("signal", "", "Signal filter: BUY, SELL, HOLD")
	jsonOut := fs.Bool("json", false, "JSON output")
	_ = fs.Parse(args)

	rsiLo, rsiHi := parseRange(*rsiRange)

	type result struct {
		Ticker string  `json:"ticker"`
		Close  float64 `json:"close"`
		RSI    float64 `json:"rsi"`
		Signal string  `json:"signal"`
		Score  float64 `json:"score"`
		ChgPct float64 `json:"change_pct"`
	}
	results := []result{}
	stocks, _ := fetchers.VPSMultiple(ctx, vn30)
	for _, s := range stocks {
		hist, err := fetchers.EntradeHistory(ctx, s.Ticker, 250)
		if err != nil || len(hist) < 50 {
			continue
		}
		closes := make([]float64, len(hist))
		for i, h := range hist {
			closes[i] = h.Close
		}
		d := technical.Dashboard(closes)
		if rsiLo > 0 || rsiHi > 0 {
			if rsiLo > 0 && d.RSI < rsiLo {
				continue
			}
			if rsiHi > 0 && d.RSI > rsiHi {
				continue
			}
		}
		if *signalFilter != "" && !strings.EqualFold(*signalFilter, string(d.Signal)) {
			continue
		}
		results = append(results, result{
			Ticker: s.Ticker, Close: s.Close, RSI: d.RSI,
			Signal: string(d.Signal), Score: d.Score, ChgPct: s.ChangePercent,
		})
	}
	sort.Slice(results, func(i, j int) bool { return results[i].Score > results[j].Score })

	if *jsonOut {
		_ = json.NewEncoder(os.Stdout).Encode(results)
		return
	}
	fmt.Printf("# Screen results (%d matches)\n\n", len(results))
	fmt.Println("| Ticker | Close | Δ% | RSI | Signal | Score |")
	fmt.Println("|---|---|---|---|---|---|")
	for _, r := range results {
		fmt.Printf("| %s | %.0f | %+.2f%% | %.1f | %s | %.0f |\n",
			r.Ticker, r.Close, r.ChgPct, r.RSI, r.Signal, r.Score)
	}
}

func parseRange(s string) (lo, hi float64) {
	if s == "" {
		return 0, 0
	}
	parts := strings.SplitN(s, "-", 2)
	if len(parts) >= 1 {
		fmt.Sscanf(parts[0], "%f", &lo)
	}
	if len(parts) == 2 {
		fmt.Sscanf(parts[1], "%f", &hi)
	}
	return
}

func runSectors(ctx context.Context) {
	stocks, _ := fetchers.VPSMultiple(ctx, vn30)
	sectors := market.RankSectorsByFlow(stocks)
	fmt.Println("# Sector flow")
	fmt.Println()
	fmt.Println("| # | Sector | Δ% trung bình | KL (tỷ) | Mã tăng/tổng |")
	fmt.Println("|---|---|---|---|---|")
	for _, s := range sectors {
		fmt.Printf("| %d | %s | %+.2f%% | %.0f | %d/%d |\n",
			s.Rank, s.Sector, s.AvgChangePct, s.TotalVolumeVND/1e9, s.AdvancesCount, s.TotalCount)
	}
}

func runGlobal(ctx context.Context) {
	symbols := make([]string, 0, len(fetchers.GlobalIndexRegistry))
	nameByIdx := map[string]string{}
	regionByIdx := map[string]string{}
	for _, e := range fetchers.GlobalIndexRegistry {
		symbols = append(symbols, e.Symbol)
		nameByIdx[e.Symbol] = e.NameVi
		regionByIdx[e.Symbol] = e.Region
	}
	quotes := fetchers.YahooMultiple(ctx, symbols)
	fmt.Println("# Global markets")
	fmt.Println()
	fmt.Println("| Khu vực | Chỉ số | Giá | Δ% | Thời gian |")
	fmt.Println("|---|---|---|---|---|")
	for _, q := range quotes {
		fmt.Printf("| %s | %s | %.2f | %+.2f%% | %s |\n",
			regionByIdx[q.Symbol], nameByIdx[q.Symbol], q.Price, q.ChangePct,
			q.MarketTime.Format("15:04 02/01"))
	}
}

func runDividends(ctx context.Context, args []string) {
	if len(args) == 0 {
		log.Fatal("usage: lmcli dividends TICKER")
	}
	ticker := strings.ToUpper(args[0])
	events, err := fetchers.VCIDividends(ctx, ticker)
	if err != nil {
		log.Fatalf("vci %s: %v", ticker, err)
	}
	fmt.Printf("# %s — Lịch sự kiện cổ phiếu\n\n", ticker)
	if len(events) == 0 {
		fmt.Println("_(Không có sự kiện nào)_")
		return
	}
	fmt.Println("| Ngày GDKHQ | Loại | Tỷ lệ | Tiền/cp | Ngày trả |")
	fmt.Println("|---|---|---|---|---|")
	for _, e := range events {
		typ := "Cổ tức tiền"
		rate := fmt.Sprintf("%.1f%%", e.RateCash)
		amt := fmt.Sprintf("%.0f VND", e.AmountVND())
		if e.EventType == 2 {
			typ = "Phát hành (CT/quyền/ESOP)"
			rate = fmt.Sprintf("%.1f%%", e.RateSplit)
			amt = "—"
		}
		fmt.Printf("| %s | %s | %s | %s | %s |\n",
			e.ExDate, typ, rate, amt, e.PaymentDate)
	}
}

func runPulse(ctx context.Context) {
	stocks, _ := fetchers.VPSMultiple(ctx, vn30)
	if len(stocks) == 0 {
		log.Fatal("no stock data")
	}

	// Domestic flow
	flow := signals.ComputeDomesticFlow(stocks)
	// Sector ranking
	sectors := market.RankSectorsByFlow(stocks)
	// Volume surges
	surges := []signals.VolumeSurge{}
	for _, s := range stocks {
		hist, err := fetchers.EntradeHistory(ctx, s.Ticker, 30)
		if err != nil || len(hist) < 20 {
			continue
		}
		var sum, sumSq float64
		for _, h := range hist[len(hist)-20:] {
			v := float64(h.Volume)
			sum += v
			sumSq += v * v
		}
		avg := sum / 20
		variance := sumSq/20 - avg*avg
		stddev := 0.0
		if variance > 0 {
			stddev = math.Sqrt(variance)
		}
		surge := signals.ClassifyVolumeSurge(s.Volume, avg, stddev, s.ChangePercent)
		surge.Ticker = s.Ticker
		if surge.Ratio >= 2.0 {
			surges = append(surges, surge)
		}
	}

	fmt.Printf("# 🌸 Nhịp thị trường — %s\n\n", time.Now().Format("02/01/2006 15:04"))
	fmt.Printf("## Áp lực mua/bán nội địa\n\n")
	fmt.Printf("**%s** — áp lực mua: %.1f%%\n\n", flow.Signal, flow.BuyPressure)

	fmt.Println("## Dòng tiền sector (top 5)")
	fmt.Println()
	fmt.Println("| # | Ngành | Δ% | KL tỷ | Tăng/Tổng |")
	fmt.Println("|---|---|---|---|---|")
	max := 5
	if len(sectors) < max {
		max = len(sectors)
	}
	for i := 0; i < max; i++ {
		s := sectors[i]
		fmt.Printf("| %d | %s | %+.2f%% | %.0f | %d/%d |\n",
			s.Rank, s.Sector, s.AvgChangePct, s.TotalVolumeVND/1e9, s.AdvancesCount, s.TotalCount)
	}
	fmt.Println()

	if len(surges) > 0 {
		fmt.Println("## Đột biến khối lượng")
		fmt.Println()
		fmt.Println("| Ticker | Tín hiệu | Tỷ lệ |")
		fmt.Println("|---|---|---|")
		for _, s := range surges {
			fmt.Printf("| %s | %s | %.1fx |\n", s.Ticker, s.Label, s.Ratio)
		}
		fmt.Println()
	}

	// Top movers
	sorted := make([]types.StockData, len(stocks))
	copy(sorted, stocks)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].ChangePercent > sorted[j].ChangePercent })
	fmt.Println("## Top tăng/giảm VN30")
	fmt.Println()
	fmt.Println("| Top tăng | Δ% | Top giảm | Δ% |")
	fmt.Println("|---|---|---|---|")
	for i := 0; i < 5 && i < len(sorted); i++ {
		up := sorted[i]
		down := sorted[len(sorted)-1-i]
		fmt.Printf("| %s | %+.2f%% | %s | %+.2f%% |\n",
			up.Ticker, up.ChangePercent, down.Ticker, down.ChangePercent)
	}
	fmt.Println()
	fmt.Println("---")
	fmt.Println("_Generated by [lotusmarket](https://github.com/ducnhd/lotusmarket) — free & open source_")
}

func runReport(ctx context.Context) {
	fs := flag.NewFlagSet("report", flag.ExitOnError)
	outPath := fs.String("out", "", "output file (default: stdout)")
	_ = fs.Parse(os.Args[2:])

	var sb strings.Builder
	captureStdout := os.Stdout
	defer func() { os.Stdout = captureStdout }()

	// Capture each subcommand output into sb
	for _, fn := range []func(){
		func() { runPulse(ctx) },
		func() { fmt.Println(); runGlobal(ctx) },
		func() { fmt.Println(); runSectors(ctx) },
	} {
		r, w, _ := os.Pipe()
		os.Stdout = w
		fn()
		_ = w.Close()
		buf := make([]byte, 1<<20)
		n, _ := r.Read(buf)
		sb.Write(buf[:n])
	}

	os.Stdout = captureStdout
	report := sb.String()
	if *outPath != "" {
		if err := os.WriteFile(*outPath, []byte(report), 0644); err != nil {
			log.Fatalf("write: %v", err)
		}
		fmt.Fprintf(os.Stderr, "wrote %s (%d bytes)\n", *outPath, len(report))
		return
	}
	fmt.Print(report)
}
