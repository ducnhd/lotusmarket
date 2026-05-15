// Weekly best-of summary — runs Sunday 09:00 VN, recaps the past 5 trading
// days. Differentiated content from the daily noise: which stocks moved most
// over the week, sector rotation, regime hints. Posts to the same Telegram
// channel.
package main

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/ducnhd/lotusmarket/go/fetchers"
	"github.com/ducnhd/lotusmarket/go/market"
	"github.com/ducnhd/lotusmarket/go/technical"
)

type weeklyMover struct {
	ticker string
	change float64 // 5-day return %
	close  float64
}

func computeWeeklyReturn(ctx context.Context, ticker string) (weeklyMover, bool) {
	hist, err := fetchers.EntradeHistory(ctx, ticker, 12)
	if err != nil || len(hist) < 6 {
		return weeklyMover{}, false
	}
	cur := hist[len(hist)-1].Close
	prev := hist[len(hist)-6].Close
	if prev <= 0 {
		return weeklyMover{}, false
	}
	return weeklyMover{
		ticker: ticker,
		change: (cur/prev - 1) * 100,
		close:  cur,
	}, true
}

func runWeekly(ctx context.Context, format string) {
	movers := []weeklyMover{}
	for _, t := range vn30 {
		if m, ok := computeWeeklyReturn(ctx, t); ok {
			movers = append(movers, m)
		}
	}
	sort.Slice(movers, func(i, j int) bool { return movers[i].change > movers[j].change })

	// Sector flow snapshot for the past week
	stocks, _ := fetchers.VPSMultiple(ctx, vn30)
	sectors := market.RankSectorsByFlow(stocks)

	// RSI distribution (regime hint) using current data
	rsiOver70 := 0
	rsiUnder30 := 0
	for _, t := range vn30[:15] { // sample 15 to keep fast
		hist, err := fetchers.EntradeHistory(ctx, t, 30)
		if err != nil || len(hist) < 15 {
			continue
		}
		closes := make([]float64, len(hist))
		for i, h := range hist {
			closes[i] = h.Close
		}
		rsi := technical.RSI(closes, 14)
		if rsi == 0 {
			continue
		}
		if rsi >= 70 {
			rsiOver70++
		}
		if rsi <= 30 {
			rsiUnder30++
		}
	}

	now := time.Now()
	weekNum := (now.YearDay() + 6) / 7
	header := fmt.Sprintf("🌸 Tuần %d/%d — VN stock recap", weekNum, now.Year())

	if format == "telegram" {
		// HTML for Telegram
		var sb strings.Builder
		fmt.Fprintf(&sb, "<b>%s</b>\n\n", header)
		sb.WriteString("📈 <b>Top tăng tuần (VN30, 5 phiên):</b>\n")
		for i := 0; i < 5 && i < len(movers); i++ {
			m := movers[i]
			fmt.Fprintf(&sb, "• <code>%s</code> %+.1f%% (%.0f)\n", m.ticker, m.change, m.close)
		}
		sb.WriteString("\n📉 <b>Top giảm tuần:</b>\n")
		for i := 0; i < 5 && i < len(movers); i++ {
			m := movers[len(movers)-1-i]
			fmt.Fprintf(&sb, "• <code>%s</code> %+.1f%% (%.0f)\n", m.ticker, m.change, m.close)
		}
		sb.WriteString("\n🏭 <b>Sector flow (top 3):</b>\n")
		maxS := 3
		if len(sectors) < maxS {
			maxS = len(sectors)
		}
		for i := 0; i < maxS; i++ {
			s := sectors[i]
			arrow := "🟢"
			if s.AvgChangePct < 0 {
				arrow = "🔴"
			}
			fmt.Fprintf(&sb, "%s %s · %+.2f%%\n", arrow, s.Sector, s.AvgChangePct)
		}
		fmt.Fprintf(&sb, "\n⚖️ <b>Regime hint:</b> %d/15 VN30 overbought (RSI≥70) · %d/15 oversold (RSI≤30)\n",
			rsiOver70, rsiUnder30)
		sb.WriteString(`<i>Mua đáy/bán đỉnh hay momentum? Tuần này nghiêng về </i>`)
		switch {
		case rsiOver70 >= 6:
			sb.WriteString("<i>thị trường nóng — cần cẩn trọng.</i>")
		case rsiUnder30 >= 6:
			sb.WriteString("<i>oversold — có thể là cơ hội nếu regime chính ổn.</i>")
		default:
			sb.WriteString("<i>cân bằng, không có signal mạnh.</i>")
		}
		sb.WriteString("\n\n")
		sb.WriteString(`📊 <a href="https://ducnhd.github.io/lotusmarket/">Dashboard</a> · `)
		sb.WriteString(`<a href="https://github.com/ducnhd/lotusmarket">GitHub</a>`)
		fmt.Print(sb.String())
		return
	}

	// Default: markdown
	var sb strings.Builder
	fmt.Fprintf(&sb, "# %s\n\nGenerated: %s\n\n", header, now.Format("2006-01-02 15:04"))
	sb.WriteString("## 📈 Top tăng tuần (VN30, 5 phiên)\n\n| Ticker | Giá | Δ% |\n|---|---|---|\n")
	for i := 0; i < 5 && i < len(movers); i++ {
		m := movers[i]
		fmt.Fprintf(&sb, "| %s | %.0f | %+.2f%% |\n", m.ticker, m.close, m.change)
	}
	sb.WriteString("\n## 📉 Top giảm tuần\n\n| Ticker | Giá | Δ% |\n|---|---|---|\n")
	for i := 0; i < 5 && i < len(movers); i++ {
		m := movers[len(movers)-1-i]
		fmt.Fprintf(&sb, "| %s | %.0f | %+.2f%% |\n", m.ticker, m.close, m.change)
	}
	sb.WriteString("\n## 🏭 Sector flow snapshot\n\n| # | Sector | Δ% | KL tỷ |\n|---|---|---|---|\n")
	maxS := 5
	if len(sectors) < maxS {
		maxS = len(sectors)
	}
	for i := 0; i < maxS; i++ {
		s := sectors[i]
		fmt.Fprintf(&sb, "| %d | %s | %+.2f%% | %.0f |\n", s.Rank, s.Sector, s.AvgChangePct, s.TotalVolumeVND/1e9)
	}
	fmt.Fprintf(&sb, "\n## ⚖️ Regime hint\n\n- %d/15 VN30 overbought (RSI≥70)\n- %d/15 VN30 oversold (RSI≤30)\n", rsiOver70, rsiUnder30)
	fmt.Print(sb.String())
}
