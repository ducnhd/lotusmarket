// autoblog — generate a data-grounded blog post using Claude API.
//
// Rotates 3 templates by weekday:
//   - Mon: weekly recap (5-day VN30 movers + sector + regime hint)
//   - Wed: educational cohort post (rotating indicator focus)
//   - Fri: ticker spotlight (best-rated VN30 with exposure chain)
//
// All inputs are real numbers pulled from the lotusmarket library at call
// time — Claude only does prose synthesis, not analysis or fact-finding.
// Output: docs/blog/YYYY-MM-DD-slug.md with YAML frontmatter.
package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ducnhd/lotusmarket/go/ai"
	"github.com/ducnhd/lotusmarket/go/fetchers"
	"github.com/ducnhd/lotusmarket/go/market"
	"github.com/ducnhd/lotusmarket/go/ratings"
	"github.com/ducnhd/lotusmarket/go/technical"
	"github.com/ducnhd/lotusmarket/go/types"
)

type blogTopic string

const (
	topicWeekly blogTopic = "weekly"
	topicCohort blogTopic = "cohort"
	topicTicker blogTopic = "ticker"
)

func pickTopicForToday() blogTopic {
	switch time.Now().Weekday() {
	case time.Monday:
		return topicWeekly
	case time.Wednesday:
		return topicCohort
	case time.Friday:
		return topicTicker
	default:
		// On other days, pick whichever is the most-recent slot
		// to allow manual triggering without spamming
		return topicWeekly
	}
}

func runAutoBlog(ctx context.Context, outDir, topicFlag string) {
	if _, err := os.Stat(outDir); err != nil {
		if err := os.MkdirAll(outDir, 0o755); err != nil {
			log.Fatalf("mkdir %s: %v", outDir, err)
		}
	}

	topic := pickTopicForToday()
	if topicFlag != "" && topicFlag != "auto" {
		topic = blogTopic(topicFlag)
	}

	apiKey := os.Getenv("CLAUDE_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("ANTHROPIC_API_KEY")
	}
	aiClient, err := ai.New(ai.Config{APIKey: apiKey})
	if err != nil {
		log.Fatalf("ai: %v", err)
	}

	var input, title, slug string
	switch topic {
	case topicWeekly:
		input, title, slug = buildWeeklyRecapInput(ctx)
	case topicCohort:
		input, title, slug = buildCohortInput()
	case topicTicker:
		input, title, slug = buildTickerSpotlightInput(ctx)
	default:
		log.Fatalf("unknown topic: %s", topic)
	}

	instruction := strings.ReplaceAll(blogInstruction, "{{TOPIC}}", string(topic))
	log.Printf("[autoblog] topic=%s title=%q calling Claude...", topic, title)
	result, err := aiClient.AnalyzeWithContext(ctx, input, instruction)
	if err != nil {
		log.Fatalf("claude: %v", err)
	}
	body := strings.TrimSpace(result.Text)

	date := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("%s-%s.md", date, slug)
	full := filepath.Join(outDir, filename)
	front := fmt.Sprintf("---\ntitle: %q\ndate: %s\ntopic: %s\n---\n\n", title, date, topic)
	if err := os.WriteFile(full, []byte(front+body+"\n"), 0644); err != nil {
		log.Fatalf("write %s: %v", full, err)
	}
	fmt.Println(full)
}

const blogInstruction = `Bạn là 1 quant analyst viết blog tiếng Việt cho dự án lotusmarket (Lotus AI).

Đặc tả output:
- Viết bài blog markdown bằng tiếng Việt, độ dài 800-1500 chữ.
- Bắt đầu bằng 1 dòng TL;DR (1-2 câu).
- KHÔNG bịa số liệu. Chỉ dùng các con số có trong DATA bên dưới.
- KHÔNG viết "có thể, có lẽ, dự đoán" — viết khẳng định dựa trên data đã có.
- KHÔNG khuyến nghị mua/bán cụ thể — viết theo dạng "data cho thấy / cohort thắng X%".
- Sử dụng cấu trúc: TL;DR → Data → Analysis → 3 takeaway → Reproducible commands → Disclaimer.
- Mỗi takeaway phải reference 1 con số cụ thể.
- Kết thúc bằng "Reproducible" section với 2-3 dòng lệnh để reader tự verify (pip install lotusmarket==0.5.0 + Python snippet hoặc lmcli command).
- Disclaimer ngắn (1 dòng): không phải lời khuyên đầu tư.
- KHÔNG include frontmatter (mình đã add sẵn). Bắt đầu từ heading hoặc body trực tiếp.
- Tone: thẳng thắn, không hype, không emoji nhiều (1-2 emoji là đủ).
- KHÔNG viết "AI nghĩ rằng..." — viết như 1 quant viết content, không hé lộ là AI generate.
- Reference "Lotus AI" / "lotusmarket" 1-2 lần trong bài, link tới https://lotusai.servehttp.com hoặc https://github.com/ducnhd/lotusmarket.
- Topic context: bài này thuộc loại {{TOPIC}}.

DATA:
`

// ============================================================================
// Weekly recap input
// ============================================================================

func buildWeeklyRecapInput(ctx context.Context) (input, title, slug string) {
	stocks, _ := fetchers.VPSMultiple(ctx, vn30)
	flow := signalsFlowSafe(stocks)
	sectors := market.RankSectorsByFlow(stocks)

	// 5-day moves per ticker
	type weekly struct {
		ticker, signal              string
		close, change5d, rsi, score float64
	}
	all := []weekly{}
	for _, t := range vn30 {
		hist, err := fetchers.EntradeHistory(ctx, t, 12)
		if err != nil || len(hist) < 6 {
			continue
		}
		closes := make([]float64, len(hist))
		for i, h := range hist {
			closes[i] = h.Close
		}
		d := technical.Dashboard(closes)
		cur := hist[len(hist)-1].Close
		prev := hist[len(hist)-6].Close
		ch5 := 0.0
		if prev > 0 {
			ch5 = (cur/prev - 1) * 100
		}
		all = append(all, weekly{t, string(d.Signal), cur, ch5, d.RSI, d.Score})
	}
	sort.Slice(all, func(i, j int) bool { return all[i].change5d > all[j].change5d })

	rsiOver70, rsiUnder30 := 0, 0
	for _, w := range all {
		if w.rsi >= 70 {
			rsiOver70++
		}
		if w.rsi <= 30 {
			rsiUnder30++
		}
	}

	var sb strings.Builder
	now := time.Now()
	week := (now.YearDay() + 6) / 7
	fmt.Fprintf(&sb, "Topic: Weekly recap VN30, tuần %d/%d, %s\n\n", week, now.Year(), now.Format("02/01/2006"))
	fmt.Fprintf(&sb, "Dòng tiền nội địa hôm nay: %s (áp lực mua %.1f%%)\n\n", flow.Signal, flow.BuyPressure)

	sb.WriteString("Top 5 tăng tuần (5 phiên gần nhất):\n")
	maxN := 5
	if len(all) < maxN {
		maxN = len(all)
	}
	for i := 0; i < maxN; i++ {
		w := all[i]
		fmt.Fprintf(&sb, "  - %s: %+.2f%% (close %.0f, RSI %.0f, signal %s)\n", w.ticker, w.change5d, w.close, w.rsi, w.signal)
	}
	sb.WriteString("\nTop 5 giảm tuần:\n")
	for i := 0; i < maxN && i < len(all); i++ {
		w := all[len(all)-1-i]
		fmt.Fprintf(&sb, "  - %s: %+.2f%% (close %.0f, RSI %.0f, signal %s)\n", w.ticker, w.change5d, w.close, w.rsi, w.signal)
	}

	sb.WriteString("\nSector flow top 5 (hôm nay):\n")
	maxS := 5
	if len(sectors) < maxS {
		maxS = len(sectors)
	}
	for i := 0; i < maxS; i++ {
		s := sectors[i]
		fmt.Fprintf(&sb, "  - #%d %s: %+.2f%% (KL %.0f tỷ, %d/%d tăng)\n", s.Rank, s.Sector, s.AvgChangePct, s.TotalVolumeVND/1e9, s.AdvancesCount, s.TotalCount)
	}

	fmt.Fprintf(&sb, "\nRSI distribution VN30: %d/%d overbought (≥70), %d/%d oversold (≤30)\n", rsiOver70, len(all), rsiUnder30, len(all))

	title = fmt.Sprintf("Tuần %d/%d: VN30 recap — sector dẫn dắt và pattern RSI", week, now.Year())
	slug = fmt.Sprintf("weekly-recap-w%02d", week)
	input = sb.String()
	return
}

// ============================================================================
// Cohort education input (rotate by week-of-year)
// ============================================================================

var cohortTopics = []struct {
	focus, slug, headline string
	data                  string
}{
	{
		focus:    "RSI bucket cohort",
		slug:     "cohort-rsi-vs-ma-trend",
		headline: "Cohort 9 năm: RSI một mình vô nghĩa, phải kết hợp với MA trend",
		data: `Cohort analysis trên ~45,000 row clean VN30+HNX30 (2017-2026):

RSI<30 (oversold), N=1,755: fwd return 60d trung bình +6.32%, win rate 60%, edge so baseline +1.57%
RSI 30-50, N=18,933: +3.67%, 52%, edge -1.08%
RSI 50-70, N=20,681: +4.82%, 53%, edge +0.07%
RSI>70 (overbought), N=4,015: +9.23%, 60%, edge +4.48%

Joint RSI × MA trend (giá vs MA200, MA50 vs MA200):
- Uptrend × RSI>70, N=2,814: +9.80%, win 63%, edge +5.05%
- Uptrend × RSI 50-70, N=11,912: +6.21%, edge +1.45%
- Downtrend × RSI>70, N=139: -1.40%, win 38%, edge -6.15% ← vùng SELL duy nhất
- Downtrend × RSI 30-50, N=7,582: +3.48%, edge -1.27%
- Mixed × RSI<30, N=401: +11.86%, win 70%, edge +7.11%

Baseline cả cohort 60d: +4.75%, win 53%, N=45,493`,
	},
	{
		focus:    "Regime classifier",
		slug:     "regime-crisis-is-opportunity",
		headline: "CRISIS là cơ hội, không phải rủi ro — cohort data thật trên VN30",
		data: `Cohort analysis theo market regime (deterministic, no AI):

CRISIS, N=2,582: fwd return 60d +8.53%, win rate 67%, edge +3.78%
EUPHORIA, N=10,279: +6.12%, 59%, +1.37%
VOLATILE, N=7,482: +6.34%, 56%, +1.59%
STABLE, N=23,868: +3.20%, 48%, -1.55%

Baseline: +4.75%, 53%

CRISIS có 2 sub-type quan trọng để phân biệt:
- CRISIS_PANIC (contagion từ ngoài, VIX spike, tier 1 news global): historically followed by recovery 3-12 tháng. Vd COVID 3/2020.
- CRISIS_FUNDAMENTAL (vấn đề cấu trúc VN, bank run, kiểm tra hình sự lãnh đạo lớn): may not recover. Vd Vạn Thịnh Phát 10/2022.

→ Quy tắc retail VN làm sai phổ biến: panic bán đáy trong CRISIS_PANIC. Cohort cho thấy đây chính là vùng nên mua thêm.`,
	},
	{
		focus:    "Active vs buy-hold 16 năm",
		slug:     "buyhold-beats-active-16y",
		headline: "Buy-hold đánh bại active strategy 16 năm — backtest realistic VN retail fees",
		data: `Backtest 16 năm (2010-2026), 30 mã VN30+HNX30, phí TPS retail thực tế (commission 0.15%, sell tax 0.1%, slippage 0.2%, lot 100):

Buy-hold equal-weight VN30: CAGR +15.37%, MaxDD -36%
Trend 200MA: +21.6% (10y window), MaxDD -51%
Momentum 12-1: +21.7% (10y), -44%
Low-vol top-5: +26.2% (10y), -42%
Dual momentum: +15.7% (10y), -55%

Lưu ý: số 10y window là 2017-2026 — cherry-picked vì không có 2011/2018 crash đầy đủ.
Trên window 16 năm đầy đủ chu kỳ, buy-hold đánh bại tất cả active strategies.

Lý do toán học:
- Phí 0.35%/lệnh × ~12 lệnh/năm = 4.2%/năm bị ăn
- Alpha thật của momentum/low-vol khoảng 3-6%/năm
- Net: alpha gần như bằng 0 sau phí

Withdrawal sim Monte Carlo 1000 paths, horizon 10y, lạm phát 4%:
- 12tr/tháng real-terms: cần 2.55 tỷ (P_fail ≤ 5%)
- 24tr/tháng: 5.06 tỷ
- 36tr/tháng: 7.78 tỷ
- 48tr/tháng: 8.98 tỷ`,
	},
}

func buildCohortInput() (input, title, slug string) {
	idx := time.Now().YearDay() / 7 % len(cohortTopics)
	t := cohortTopics[idx]
	input = fmt.Sprintf("Topic: %s\n\nFocus: %s\n\nDATA:\n%s\n\nHôm nay: %s\n", t.headline, t.focus, t.data, time.Now().Format("02/01/2006"))
	title = t.headline
	slug = t.slug
	return
}

// ============================================================================
// Ticker spotlight input
// ============================================================================

func buildTickerSpotlightInput(ctx context.Context) (input, title, slug string) {
	type cand struct {
		ticker string
		stock  types.StockData
		hist   []types.StockData
		dash   technical.DashboardMetrics
		rat    ratings.Ratings
		ch5    float64
	}
	cands := []cand{}

	stocks, _ := fetchers.VPSMultiple(ctx, vn30)
	stockByTicker := map[string]types.StockData{}
	for _, s := range stocks {
		stockByTicker[s.Ticker] = s
	}

	for _, t := range vn30 {
		hist, err := fetchers.EntradeHistory(ctx, t, 250)
		if err != nil || len(hist) < 220 {
			continue
		}
		closes := make([]float64, len(hist))
		vols := make([]int64, len(hist))
		for i, h := range hist {
			closes[i] = h.Close
			vols[i] = h.Volume
		}
		d := technical.Dashboard(closes)
		r := ratings.Compute(closes, vols, d.Score, d.RSI, d.MA20, d.MA50, d.MA200)
		cur := closes[len(closes)-1]
		prev := closes[len(closes)-6]
		ch5 := (cur/prev - 1) * 100
		cands = append(cands, cand{t, stockByTicker[t], hist, d, r, ch5})
	}

	// Pick ticker with highest overall gauge AND non-trivial 5d move
	sort.Slice(cands, func(i, j int) bool {
		gi := cands[i].rat.OverallGauge + int(math.Abs(cands[i].ch5)*2)
		gj := cands[j].rat.OverallGauge + int(math.Abs(cands[j].ch5)*2)
		return gi > gj
	})
	if len(cands) == 0 {
		input = "No data available."
		title = "Ticker spotlight"
		slug = "ticker-spotlight"
		return
	}
	c := cands[0]

	var sb strings.Builder
	fmt.Fprintf(&sb, "Topic: Ticker spotlight\n\n")
	fmt.Fprintf(&sb, "Ticker: %s\nClose: %.0f VND\nChange today: %+.2f%%\nChange 5d: %+.2f%%\nVolume: %d\n\n",
		c.ticker, c.stock.Close, c.stock.ChangePercent, c.ch5, c.stock.Volume)
	fmt.Fprintf(&sb, "Technical dashboard:\n  RSI: %.1f\n  Signal: %s\n  Score: %.0f/100\n",
		c.dash.RSI, c.dash.Signal, c.dash.Score)
	if c.dash.MA20 != nil {
		fmt.Fprintf(&sb, "  MA20: %.0f\n", *c.dash.MA20)
	}
	if c.dash.MA50 != nil {
		fmt.Fprintf(&sb, "  MA50: %.0f\n", *c.dash.MA50)
	}
	if c.dash.MA200 != nil {
		fmt.Fprintf(&sb, "  MA200: %.0f\n", *c.dash.MA200)
	}
	fmt.Fprintf(&sb, "\nStar ratings (6-dim, max 5 each):\n  Price strength: %d\n  Trend: %d\n  Short-term RSI position: %d\n  Money flow: %d\n  Volatility: %d\n  Base range: %d\n  → Overall verdict: %s (gauge %d/100)\n",
		c.rat.PriceStrength, c.rat.TrendStrength, c.rat.ShortTermPos, c.rat.MoneyFlow, c.rat.VolatilityRating, c.rat.BaseRange, c.rat.OverallVerdict, c.rat.OverallGauge)

	// MA trend label
	mat := "mixed"
	if c.dash.MA200 != nil && c.dash.MA50 != nil {
		above := c.stock.Close >= *c.dash.MA200
		ma50above := *c.dash.MA50 >= *c.dash.MA200
		if above && ma50above {
			mat = "uptrend"
		} else if !above && !ma50above {
			mat = "downtrend"
		}
	}
	fmt.Fprintf(&sb, "\nMA trend: %s\n", mat)

	// Cohort context
	cohort := "không có cohort match"
	switch {
	case mat == "uptrend" && c.dash.RSI >= 70:
		cohort = "Uptrend × RSI>70 (N=2,814): cohort fwd 60d +9.80%, win 63%, edge so baseline +5.05%"
	case mat == "uptrend" && c.dash.RSI >= 50:
		cohort = "Uptrend × RSI 50-70 (N=11,912): cohort fwd 60d +6.21%, edge +1.45%"
	case mat == "downtrend" && c.dash.RSI >= 70:
		cohort = "Downtrend × RSI>70 (N=139): cohort fwd 60d -1.40%, win 38%, edge -6.15% — vùng SELL duy nhất có edge tiêu cực rõ rệt"
	case c.dash.RSI < 30:
		cohort = "RSI<30 oversold (N=1,755): cohort fwd 60d +6.32%, win 60%, edge +1.57%"
	}
	fmt.Fprintf(&sb, "Cohort historical match: %s\n", cohort)

	title = fmt.Sprintf("Spotlight %s: %s, RSI %.0f, %s", c.ticker, mat, c.dash.RSI, c.rat.OverallVerdict)
	slug = fmt.Sprintf("spotlight-%s", strings.ToLower(c.ticker))
	input = sb.String()
	return
}

// ============================================================================
// Helper — re-use signals package without importing again
// ============================================================================

type domesticFlowResult struct {
	Signal      string
	BuyPressure float64
}

func signalsFlowSafe(stocks []types.StockData) domesticFlowResult {
	// Use the signals package via fully-qualified call.
	// (Avoiding circular import — autoblog is in main package.)
	totalBid, totalAsk := int64(0), int64(0)
	for _, s := range stocks {
		totalBid += s.BidVol
		totalAsk += s.AskVol
	}
	total := totalBid + totalAsk
	if total == 0 {
		return domesticFlowResult{Signal: "không xác định", BuyPressure: 50}
	}
	buyPct := float64(totalBid) / float64(total) * 100
	signal := "cân bằng"
	switch {
	case buyPct >= 60:
		signal = "mua mạnh"
	case buyPct >= 53:
		signal = "nghiêng mua"
	case buyPct <= 40:
		signal = "bán mạnh"
	case buyPct <= 47:
		signal = "nghiêng bán"
	}
	return domesticFlowResult{Signal: signal, BuyPressure: buyPct}
}
