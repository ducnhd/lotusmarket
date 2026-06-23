// autoblog — generate diverse, story-driven blog posts using Claude API.
//
// Topic catalog (10 types, randomly selected, no repeat within 7 days):
//
//	tech-classic      Phân tích kỹ thuật cổ điển (RSI/MACD/Bollinger/...) chứng minh đúng/sai
//	myth-buster       Backtest niềm tin sai lầm phổ biến trong cộng đồng retail VN
//	random-ticker     Đánh giá ngẫu nhiên 1 mã VN30 với star ratings + cohort context
//	career-investing  Nghề nghiệp + đầu tư (full-time/part-time, FIRE math, side income)
//	supply-chain      Phân tích input/output 1 công ty + sensitivity với commodity/FX/peer
//	psychology        Behavioral finance với góc nhìn VN (loss aversion, recency, herd)
//	data-insight      Surprising findings từ cohort 47K data points
//	comparison        So sánh 2 mã cùng ngành (VCB vs CTG, HPG vs HSG, VNM vs MSN, ...)
//	regime-now        Current market regime breakdown + insights
//	news-impact       Macro/global events tuần qua + ảnh hưởng VN
//
// Selection algorithm: scan docs/blog/ for last 7 days of "<date>-<topic>.md",
// exclude those topics, pick random from remaining. Ensures variety.
//
// Each topic has its own data builder + prompt twist. Data anchors are real
// (cohort numbers from validated backtests + live lmcli where applicable);
// Claude does prose synthesis only. No fabrication.
package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ducnhd/lotusmarket/go/ai"
	"github.com/ducnhd/lotusmarket/go/fetchers"
	"github.com/ducnhd/lotusmarket/go/ratings"
	"github.com/ducnhd/lotusmarket/go/technical"
	"github.com/ducnhd/lotusmarket/go/types"
)

// ============================================================================
// Topic registry
// ============================================================================

type topicSpec struct {
	Key     string
	Title   string
	Slug    string
	Twist   string // extra prompt instruction unique to this topic
	Builder func(ctx context.Context) (data, title, slug string)
}

var topicRegistry = []topicSpec{
	{Key: "tech-classic", Builder: buildTechClassic},
	{Key: "myth-buster", Builder: buildMythBuster},
	{Key: "random-ticker", Builder: buildRandomTicker},
	{Key: "career-investing", Builder: buildCareerInvesting},
	{Key: "supply-chain", Builder: buildSupplyChain},
	{Key: "psychology", Builder: buildPsychology},
	{Key: "data-insight", Builder: buildDataInsight},
	{Key: "comparison", Builder: buildComparison},
	{Key: "regime-now", Builder: buildRegimeNow},
	{Key: "news-impact", Builder: buildNewsImpact},
}

// ============================================================================
// Entry point + selector
// ============================================================================

func runAutoBlog(ctx context.Context, outDir, topicFlag string) {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		log.Fatalf("mkdir %s: %v", outDir, err)
	}

	// Load variant recency so builders can avoid recently-published content.
	recentVariantKeys = recentContentKeys(outDir, variantCooldownDays)

	apiKey := os.Getenv("CLAUDE_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("ANTHROPIC_API_KEY")
	}
	aiClient, err := ai.New(ai.Config{APIKey: apiKey})
	if err != nil {
		log.Fatalf("ai: %v", err)
	}

	t := selectTopic(outDir, topicFlag)
	log.Printf("[autoblog] selected topic: %s", t.Key)

	data, title, slug := t.Builder(ctx)
	if data == "" {
		log.Fatalf("[autoblog] empty data for topic %s — abort", t.Key)
	}

	instruction := blogPromptBase + "\n\nTopic-specific guidance:\n" + t.Twist + "\n\nDATA:\n"
	result, err := aiClient.AnalyzeWithContext(ctx, data, instruction)
	if err != nil {
		log.Fatalf("claude: %v", err)
	}
	body := strings.TrimSpace(result.Text)

	date := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("%s-%s.md", date, slug)
	full := filepath.Join(outDir, filename)
	front := fmt.Sprintf("---\ntitle: %q\ndate: %s\ntopic: %s\n---\n\n", title, date, t.Key)
	if err := os.WriteFile(full, []byte(front+body+"\n"), 0644); err != nil {
		log.Fatalf("write %s: %v", full, err)
	}
	fmt.Println(full)
}

// Recency windows. Two granularities work together:
//   - topicCooldownDays: how long before a topic *type* may recur. Controls the
//     single-variant timely topics (regime-now, news-impact) — their only lever.
//   - variantCooldownDays: how long before a specific *variant* (slug/content)
//     may recur. Carries content diversity even when a topic type returns.
const (
	topicCooldownDays   = 14
	variantCooldownDays = 60
)

// recentVariantKeys maps a normalized content key -> days ago (smaller = more
// recent). Populated once at the start of runAutoBlog and consulted by builders
// via pickFresh so variant selection avoids recently-published content.
var recentVariantKeys map[string]int

// slugTopicPrefixes maps a filename slug prefix to its topic key. Order matters:
// the first matching prefix wins. Built from the actual slug shapes the builders
// emit (e.g. comparison -> "compare-", data-insight -> "insight-"), which the
// old prefix-vs-key matching got wrong for half the topics.
var slugTopicPrefixes = []struct{ prefix, topic string }{
	{"tech-classic", "tech-classic"},
	{"myth-", "myth-buster"},
	{"random-", "random-ticker"},
	{"career-", "career-investing"},
	{"supply-chain", "supply-chain"},
	{"psychology", "psychology"},
	{"insight-", "data-insight"},
	{"compare-", "comparison"},
	{"regime-now", "regime-now"},
	{"news-impact", "news-impact"},
}

// topicOf returns the topic key for a filename slug (the part after YYYY-MM-DD-).
func topicOf(rest string) string {
	for _, m := range slugTopicPrefixes {
		if strings.HasPrefix(rest, m.prefix) {
			return m.topic
		}
	}
	return rest
}

// contentKey normalizes a slug-rest into a stable content identity. The timely
// topics carry a -MM-DD suffix; we strip it so every daily snapshot collapses to
// one key for recency purposes. Catalog topics already carry their variant in
// the slug, so they pass through unchanged.
func contentKey(rest string) string {
	for _, p := range []string{"regime-now", "news-impact"} {
		if strings.HasPrefix(rest, p+"-") {
			return p
		}
	}
	return rest
}

// pickFresh returns the index of the candidate whose slug was used least
// recently. Never-used slugs win (random among them); if every candidate was
// used inside the variant window, the one used longest ago wins.
func pickFresh(slugs []string) int {
	fresh := []int{}
	for i, s := range slugs {
		if _, used := recentVariantKeys[contentKey(s)]; !used {
			fresh = append(fresh, i)
		}
	}
	if len(fresh) > 0 {
		return fresh[rand.Intn(len(fresh))]
	}
	oldest, bestAge := []int{}, -1
	for i, s := range slugs {
		age := recentVariantKeys[contentKey(s)]
		if age > bestAge {
			bestAge, oldest = age, []int{i}
		} else if age == bestAge {
			oldest = append(oldest, i)
		}
	}
	if len(oldest) == 0 {
		return rand.Intn(len(slugs))
	}
	return oldest[rand.Intn(len(oldest))]
}

// selectTopic — pick a topic not used in the last topicCooldownDays. If
// topicFlag is set, honor it directly.
func selectTopic(outDir, topicFlag string) topicSpec {
	if topicFlag != "" && topicFlag != "auto" {
		for _, t := range topicRegistry {
			if t.Key == topicFlag {
				return t
			}
		}
		log.Printf("[autoblog] unknown topic %q, falling back to auto", topicFlag)
	}
	recent := recentTopicKeys(outDir, topicCooldownDays)
	eligible := []topicSpec{}
	for _, t := range topicRegistry {
		if !contains(recent, t.Key) {
			eligible = append(eligible, t)
		}
	}
	if len(eligible) == 0 {
		eligible = topicRegistry
	}
	rand.Seed(time.Now().UnixNano())
	return eligible[rand.Intn(len(eligible))]
}

func recentTopicKeys(dir string, days int) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	cutoff := time.Now().AddDate(0, 0, -days)
	keys := []string{}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		// Filename pattern: YYYY-MM-DD-<rest>.md
		name := strings.TrimSuffix(e.Name(), ".md")
		if len(name) < 11 {
			continue
		}
		d, err := time.Parse("2006-01-02", name[:10])
		if err != nil || d.Before(cutoff) {
			continue
		}
		keys = append(keys, topicOf(name[11:]))
	}
	return keys
}

// recentContentKeys scans outDir for posts within the last `days` and returns a
// map of content key -> days ago (smaller = more recent).
func recentContentKeys(dir string, days int) map[string]int {
	out := map[string]int{}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return out
	}
	now := time.Now()
	cutoff := now.AddDate(0, 0, -days)
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		name := strings.TrimSuffix(e.Name(), ".md")
		if len(name) < 11 {
			continue
		}
		d, err := time.Parse("2006-01-02", name[:10])
		if err != nil || d.Before(cutoff) {
			continue
		}
		key := contentKey(name[11:])
		age := int(now.Sub(d).Hours() / 24)
		if prev, ok := out[key]; !ok || age < prev {
			out[key] = age
		}
	}
	return out
}

func contains(s []string, v string) bool {
	for _, x := range s {
		if x == v {
			return true
		}
	}
	return false
}

// ============================================================================
// Base prompt — applies to ALL topics
// ============================================================================

const blogPromptBase = `Bạn là 1 quant analyst viết blog tiếng Việt cho dự án Lotus AI / lotusmarket. Bạn KHÔNG được tiết lộ mình là AI.

OUTPUT spec:
- Markdown tiếng Việt 800-1500 chữ. Hấp dẫn, có story angle, KHÔNG khô khan.
- KHÔNG include YAML frontmatter (đã có sẵn). Bắt đầu từ heading hoặc body.
- Mở đầu: 1 dòng TL;DR + 1 câu hook gây tò mò (KHÔNG bắt đầu bằng "TL;DR" literal — viết tự nhiên).
- Câu chuyện > số liệu khô. Số liệu cụ thể là anchor, story là khung.
- KHÔNG bịa số liệu. CHỈ dùng số trong DATA section.
- KHÔNG viết "có thể, có lẽ, dự đoán" — viết khẳng định từ data.
- KHÔNG khuyến nghị mua/bán cụ thể. Phrasing: "data cho thấy cohort thắng X%", "lịch sử cho thấy...".
- Cuối bài: 3 takeaway ngắn + 1 đoạn "Verify reproducible" (lệnh pip install + 1-3 dòng code Python hoặc lmcli command).
- 1 dòng disclaimer cuối: không phải lời khuyên đầu tư.
- Link 1-2 lần: https://lotusai.servehttp.com hoặc https://github.com/ducnhd/lotusmarket.
- Tone: thẳng, không hype, sarcasm nhẹ với myth/hype trong cộng đồng retail VN.
- Emoji: tối đa 2-3, dùng có chủ đích, không spam.

QUAN TRỌNG:
- Nếu DATA không đủ để khẳng định, ghi rõ "data hạn chế, finding sơ bộ" thay vì bịa.
- Đừng dùng "chúng ta", "chúng tôi" — viết theo ngôi đơn (mình / tôi / dạng impersonal).
`

// ============================================================================
// Topic builders
// ============================================================================

// ---------- T1: tech-classic ----------

var techClassicCatalog = []struct {
	indicator, headline, data string
}{
	{
		indicator: "RSI > 70",
		headline:  "RSI > 70 có thực sự là tín hiệu BÁN? Backtest 9 năm VN30 nói ngược",
		data: `Cohort analysis VN30+HNX30, 45,000 row clean (2017-2026):

RSI buckets, fwd return 60d:
- RSI < 30 (oversold), N=1,755: +6.32%, win 60%, edge +1.57%
- RSI 30-50, N=18,933: +3.67%, win 52%, edge -1.08%
- RSI 50-70, N=20,681: +4.82%, win 53%, edge +0.07%
- RSI > 70 (overbought), N=4,015: +9.23%, win 60%, edge +4.48%

Joint với MA trend:
- Uptrend (close>MA200, MA50>MA200) × RSI>70, N=2,814: +9.80%, win 63%, edge +5.05%
- Downtrend × RSI>70, N=139: -1.40%, win 38%, edge -6.15% ← vùng SELL duy nhất

Baseline 60d cohort: +4.75%, win 53%.

Hypothesis cũ trong cộng đồng VN: "RSI > 70 = chốt lời ngay". Backtest cho thấy SAI trong uptrend (cohort thắng mạnh), ĐÚNG chỉ khi downtrend.`,
	},
	{
		indicator: "Golden Cross (MA50 vượt MA200)",
		headline:  "Golden Cross — tín hiệu thần kỳ hay over-rated? Data VN30 nói gì",
		data: `Lưu ý: data dưới được trích từ cohort phân loại MA trend.

Định nghĩa "uptrend full" = (close > MA200) AND (MA50 > MA200). Đây là vùng SAU khi Golden Cross đã hoàn thành.

Cohort fwd return 60d trong Uptrend:
- Uptrend × RSI 30-50, N=7,142: +4.00%, edge -0.76%
- Uptrend × RSI 50-70, N=11,912: +6.21%, edge +1.45%
- Uptrend × RSI > 70, N=2,814: +9.80%, edge +5.05%
- Uptrend × RSI < 30, N=89: +6.88%, edge +2.13%

Baseline 60d: +4.75%

Insight:
- Sau Golden Cross, momentum continuation (RSI > 70) đem lại edge lớn nhất.
- Không phải "vào ngay sau cross". Cohort cho thấy edge tăng dần theo độ "đỏ" của RSI.
- Downtrend full (close<MA200, MA50<MA200) × RSI>70 chỉ N=139 nhưng cho edge -6.15% — đây là bear trap, KHÔNG nên fade.`,
	},
	{
		indicator: "MACD bullish crossover",
		headline:  "MACD bullish cross trên VN30: edge thực có bằng tin đồn?",
		data: `Cohort theo MACD signal:
- MACD bullish (line > signal), N=23,743: fwd 60d +5.40%, win 54%, edge +0.65%
- MACD bearish (line ≤ signal), N=21,916: fwd 60d +4.05%, win 52%, edge -0.70%

Baseline 60d: +4.75%

Insight:
- MACD bullish cho edge +0.65% so với baseline — DƯƠNG nhưng KHÔNG đột phá.
- Edge của MACD nhỏ hơn rất nhiều so với RSI bucket (RSI>70 edge +4.48%) hoặc MA trend (uptrend edge +1.20%).
- Trong khi đó, mỗi lần đảo MACD trade tốn ~0.35% phí + slippage TPS retail.
- Net edge sau phí cho retail: gần bằng 0. MACD alone không đủ để vượt phí.

So sánh: cohort uptrend + RSI > 70 (kết hợp 2 chỉ báo) cho edge +5.05% — gấp 7-8 lần MACD alone.`,
	},
	{
		indicator: "Wyckoff stage classification",
		headline:  "Wyckoff 4 giai đoạn — VN30 thắng giai đoạn nào? (spoiler: không phải markup)",
		data: `Cohort theo Wyckoff stage (1=accumulation, 2=markup, 3=distribution, 4=decline):

- Stage 1 (accumulation), N=8,084: fwd 60d +2.83%, win 48%, edge -1.92%
- Stage 2 (markup, uptrend), N=19,259: +5.46%, win 53%, edge +0.71%
- Stage 3 (distribution), N=5,111: +6.63%, win 60%, edge +1.88% ← thắng nhất
- Stage 4 (decline), N=8,506: +4.24%, win 55%, edge -0.51%

Baseline: +4.75%, win 53%

Insight ngược trực giác:
- Wyckoff stage 3 (distribution — đỉnh đang phân phối) lại có forward return cao nhất.
- Tâm lý retail VN: thấy giá cao → cho là "đu đỉnh" → bán. Cohort cho thấy đây là vùng có edge dương rõ rệt.
- Stage 1 (accumulation) — vùng "mua đáy gom hàng" mà mọi guru rao giảng — thực ra có edge ÂM (-1.92%) và win rate dưới baseline.
- Trader trung bình làm ngược: gom stage 1 (sai), bán stage 3 (sai).`,
	},
	{
		indicator: "Volume surge >2x",
		headline:  "Volume tăng đột biến — có phải tín hiệu \"smart money\"? Data nói gì",
		data: `Cohort phân tích volume surge (volume hôm nay >= 2× MA20 volume):

Trên data backtest 10 năm:
- Volume surge với price up >0.5%: thường gắn label "accumulation" (mua chủ động).
- Volume surge với price down >0.5%: thường gắn label "distribution" (bán chủ động).
- Volume surge với price flat: "high activity / churning".

Edge thực tế đo được:
- Volume surge >2σ so với baseline: edge fwd return chỉ +0.7% — yếu hơn cả MACD bullish.
- "Smart money" narrative phổ biến hơn realistic edge gấp nhiều lần.

So sánh edge các signal:
- Uptrend × RSI > 70: +5.05% ← cao nhất
- CRISIS regime: +3.78%
- Wyckoff stage 3: +1.88%
- MACD bullish: +0.65%
- Volume surge >2σ: +0.7%
- Volume surge alone (không kết hợp): noise.

Insight: Volume surge KHÔNG phải edge độc lập. Phải kết hợp với MA trend + RSI mới có meaning.`,
	},
}

func buildTechClassic(ctx context.Context) (data, title, slug string) {
	slugs := make([]string, len(techClassicCatalog))
	for i, c := range techClassicCatalog {
		slugs[i] = "tech-classic-" + safeSlug(c.indicator)
	}
	t := techClassicCatalog[pickFresh(slugs)]
	data = t.data
	title = t.headline
	slug = "tech-classic-" + safeSlug(t.indicator)
	return
}

// ---------- T2: myth-buster ----------

var mythCatalog = []struct {
	myth, headline, data string
}{
	{
		myth:     "Mua đáy bán đỉnh",
		headline: "\"Mua đáy bán đỉnh\" — câu thần chú phổ biến nhất, và sai nhất với VN30",
		data: `Pattern observation:
- Retail VN thường mua khi giá rớt mạnh (RSI<30, "đáy") và bán khi giá tăng mạnh (RSI>70, "đỉnh").
- Backtest cohort 9 năm VN30+HNX30 (N=45,493) cho thấy:
  - RSI<30 (oversold) fwd 60d: +6.32%, win 60% ← OK nhưng N nhỏ (1,755)
  - RSI>70 (overbought) fwd 60d: +9.23%, win 60% ← thắng MẠNH hơn

- Joint với MA trend:
  - Uptrend × RSI>70: +9.80%, win 63% — vùng "tưởng là đỉnh nhưng letting winner run"
  - Downtrend × RSI>70: -1.40%, win 38% — vùng SELL duy nhất

Baseline 60d: +4.75%.

Vấn đề: retail VN thực hành "mua đáy bán đỉnh" mà không phân biệt MA trend. Result: bán đúng lúc nên giữ, mua đúng lúc nên tránh.

30-day plan audit cá nhân: BUY/HOLD recommendations đúng 14/15 (93%), SELL/CUT_LOSS đúng 0/3.`,
	},
	{
		myth:     "Khủng hoảng thị trường = rút tiền ngay",
		headline: "Crisis là rủi ro hay cơ hội? Backtest regime VN cho 1 câu trả lời",
		data: `Cohort theo market regime (deterministic, không AI):

- CRISIS, N=2,582: fwd 60d +8.53%, win 67%, edge +3.78% ← HIGHEST
- EUPHORIA, N=10,279: +6.12%, 59%, +1.37%
- VOLATILE, N=7,482: +6.34%, 56%, +1.59%
- STABLE, N=23,868: +3.20%, 48%, -1.55%

Baseline: +4.75%, 53%

Phân loại CRISIS theo nguyên nhân:
- CRISIS_PANIC (contagion từ ngoài, VIX spike, news tier 1 global): historically hồi phục 3-12 tháng. Vd COVID 3/2020 VN-Index về cũ trong 8 tháng.
- CRISIS_FUNDAMENTAL (vấn đề cấu trúc nội tại VN, bank run, hình sự lãnh đạo lớn): may not recover. Vd VTP 10/2022, BĐS sập, mã liên quan vẫn dưới giá pre-crisis.

Retail VN thường ngược: panic bán đáy trong CRISIS_PANIC. Cohort cho thấy đây là vùng nên giữ hoặc mua thêm — fwd 60d +8.53% > baseline +4.75%.`,
	},
	{
		myth:     "Cổ phiếu blue chip ổn định = an toàn",
		headline: "\"Blue chip an toàn\" — backtest 16 năm VN30 nói khác",
		data: `Backtest 16 năm (2010-2026) trên 30 mã VN30+HNX30. Realistic phí TPS retail.

VNM (Vinamilk, được coi là blue chip an toàn nhất VN):
- 2010 close: ~24,000 VND → 2026 close: ~58,000 VND
- Tổng 16 năm: +145% gross
- CAGR: +5.6%/năm

So với:
- HPG: +1,270% (CAGR +17.4%)
- FPT: +1,210% (+17.1%)
- CTG: +530% (+11.9%)
- ACB: +374% (+10.0%)

Buy-hold equal-weight VN30 16 năm: CAGR +15.37%.

Insight: "An toàn" trong tâm lý retail = "ít biến động" → thường đồng nghĩa "ít tăng". 16 năm sở hữu VNM có lẽ là một trong những choice tệ nhất so với mặt bằng VN30.

Define lại "an toàn":
- KHÔNG phải = biến động thấp ngắn hạn
- LÀ = sống sót qua chu kỳ + CAGR dương qua 10+ năm
- VN30 equal-weight đạt cả 2.`,
	},
	{
		myth:     "Sell in May and go away",
		headline: "\"Sell in May\" trên VN30 — câu thần chú từ Mỹ có áp dụng được?",
		data: `Pattern test: tổng hợp return tháng 5 trên VN-Index 2010-2025.

Note: Đây là analysis quan sát, không phải backtest đầy đủ.

Câu nói "Sell in May and go away" gốc từ thị trường Mỹ — quan sát thấy May-Oct historically yếu hơn Nov-Apr.

Cohort fwd return 60d theo regime (VN data):
- STABLE regime (chiếm 52% cohort): +3.20%, win 48%, edge -1.55%
- VOLATILE regime: +6.34%, win 56%

Trên VN data:
- Tháng 5-6 thường rơi vào VOLATILE/STABLE.
- Tháng 11-12 thường rơi vào EUPHORIA hoặc tail của uptrend.

Edge từ "bán tháng 5 mua tháng 11":
- Bỏ lỡ EUPHORIA mid-year nếu có
- Chịu phí 2 lần switch/năm (~0.7%)
- Tax 0.1% khi sell

Net: edge có thể không tồn tại trên VN khi sau phí.

Caveat: cần backtest đầy đủ với data 16 năm để confirm. Hiện tại finding sơ bộ: "Sell in May" trên VN data có lẽ là noise, không phải edge thật.`,
	},
	{
		myth:     "Cổ tức cao = đầu tư tốt",
		headline: "Cổ tức 10% — sang trọng hay bẫy giảm giá? Phân tích data ex-date VN30",
		data: `Phân tích thực tế: VN tickers thường trả cổ tức 5-15% par value (=500-1500 VND/cổ phiếu nếu par = 10,000).

Cơ chế:
- Ex-date: giá tự động giảm đúng = số tiền cổ tức
- Vd: cổ phiếu 30,000 VND, cổ tức 1,500 → ex-date giá còn ~28,500
- Net change: 0 (giá - cổ tức cash = giá pre-cổ tức)

Backtest 60 ngày sau ex-date:
- Tickers có yield cao (>8% par/năm) thường là mature businesses tăng trưởng chậm.
- Forward return 60d post ex-date ~baseline market: +4-5%.
- KHÔNG có "free money" từ cổ tức.

So sánh CAGR 16 năm:
- VNM (cổ tức ~5-8%/năm, được coi là "stable dividend stock"): +5.6%/năm
- HPG (cổ tức 0-5%/năm, growth-focused): +17.4%/năm
- VN30 equal-weight buy-hold: +15.37%/năm

Insight: yield cao thường là proxy cho stagnant growth. "Cổ tức cao = an toàn" là trade-off thật sự, không phải free upside.`,
	},
}

func buildMythBuster(ctx context.Context) (data, title, slug string) {
	slugs := make([]string, len(mythCatalog))
	for i, c := range mythCatalog {
		slugs[i] = "myth-" + safeSlug(c.myth)
	}
	t := mythCatalog[pickFresh(slugs)]
	data = "Myth dưới phân tích: \"" + t.myth + "\"\n\n" + t.data
	title = t.headline
	slug = "myth-" + safeSlug(t.myth)
	return
}

// ---------- T3: random-ticker ----------

func buildRandomTicker(ctx context.Context) (data, title, slug string) {
	rand.Seed(time.Now().UnixNano())
	slugs := make([]string, len(vn30))
	for i, tk := range vn30 {
		slugs[i] = "random-" + strings.ToLower(tk)
	}
	ticker := vn30[pickFresh(slugs)]

	hist, err := fetchers.EntradeHistory(ctx, ticker, 250)
	if err != nil || len(hist) < 50 {
		// Fallback: use a known ticker with synthetic context
		return synthTickerData(ticker), fmt.Sprintf("Spotlight %s — đánh giá ngẫu nhiên một mã VN30", ticker), "random-" + strings.ToLower(ticker)
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
	var ch5, ch20 float64
	if len(closes) > 6 {
		ch5 = (cur/closes[len(closes)-6] - 1) * 100
	}
	if len(closes) > 21 {
		ch20 = (cur/closes[len(closes)-21] - 1) * 100
	}

	mat := "mixed"
	if d.MA200 != nil && d.MA50 != nil {
		if cur >= *d.MA200 && *d.MA50 >= *d.MA200 {
			mat = "uptrend"
		} else if cur < *d.MA200 && *d.MA50 < *d.MA200 {
			mat = "downtrend"
		}
	}

	cohort := matchCohort(mat, d.RSI)

	var sb strings.Builder
	fmt.Fprintf(&sb, "Ticker (randomly selected): %s\n", ticker)
	fmt.Fprintf(&sb, "Close: %.0f VND\nChange 5d: %+.2f%%\nChange 20d: %+.2f%%\nRSI: %.1f\nSignal: %s\nTechnical score: %.0f/100\nMA trend: %s\n\n",
		cur, ch5, ch20, d.RSI, d.Signal, d.Score, mat)
	fmt.Fprintf(&sb, "Star ratings (6-dim, 1-5 each):\n  Price strength: %d\n  Trend strength: %d\n  Short-term RSI position: %d\n  Money flow: %d\n  Volatility: %d\n  Base range: %d\n  → Overall: %s, gauge %d/100\n\n",
		r.PriceStrength, r.TrendStrength, r.ShortTermPos, r.MoneyFlow, r.VolatilityRating, r.BaseRange, r.OverallVerdict, r.OverallGauge)
	fmt.Fprintf(&sb, "Cohort historical match:\n%s\n", cohort)

	data = sb.String()
	title = fmt.Sprintf("%s hôm nay: %s, RSI %.0f, %s — random spotlight", ticker, mat, d.RSI, r.OverallVerdict)
	slug = "random-" + strings.ToLower(ticker)
	return
}

func synthTickerData(ticker string) string {
	return fmt.Sprintf(`Ticker (randomly selected): %s

Lưu ý: Live data không khả dụng tại thời điểm tạo bài (có thể VN market đóng cửa).

Hãy viết một bài giới thiệu ngắn về %s — đặc điểm ngành, vị thế VN30, lý do mã này thường được retail VN thảo luận. Đặt câu hỏi: nếu phải đánh giá %s khách quan, retail nên nhìn vào các chỉ số nào?

Các chỉ số đáng nhìn (gợi ý structure, không phải số liệu cụ thể):
- MA trend (close vs MA200, MA50 vs MA200) — phân biệt uptrend/downtrend/mixed
- RSI 14 — vị thế ngắn hạn
- Volume vs MA20 volume — money flow
- 6-dim star ratings (price strength, trend, short-term, money flow, volatility, base range)

Vì là bài random, hãy đề xuất reader tự chạy "lmcli rate %s" để có data tươi.`, ticker, ticker, ticker, ticker)
}

func matchCohort(mat string, rsi float64) string {
	switch {
	case mat == "uptrend" && rsi >= 70:
		return "Cohort Uptrend × RSI>70 (N=2,814): fwd 60d trung bình +9.80%, win rate 63%, edge so baseline +5.05%. Đây là cohort thắng mạnh nhất trong toàn bộ phân loại."
	case mat == "uptrend" && rsi >= 50:
		return "Cohort Uptrend × RSI 50-70 (N=11,912): fwd 60d +6.21%, edge +1.45%. Sweet spot momentum-friendly."
	case mat == "uptrend":
		return "Cohort Uptrend × RSI<50 (N≈7,231): fwd 60d ~+4.00%, edge ~-0.75%. Pullback in uptrend — chờ thêm signal."
	case mat == "downtrend" && rsi >= 70:
		return "Cohort Downtrend × RSI>70 (N=139): fwd 60d -1.40%, win 38%, edge -6.15%. Đây là VÙNG SELL DUY NHẤT có edge tiêu cực rõ rệt."
	case mat == "downtrend" && rsi < 30:
		return "Cohort Downtrend × RSI<30 (N=1,130): fwd 60d +4.68%, edge -0.07%. Oversold trong downtrend — borderline."
	case rsi < 30:
		return "Cohort RSI<30 oversold (N=1,755): fwd 60d +6.32%, win 60%, edge +1.57%."
	default:
		return "Cohort không match pattern mạnh — đây là vùng neutral, baseline cohort edge ~0%."
	}
}

// ---------- T4: career-investing ----------

func buildCareerInvesting(ctx context.Context) (data, title, slug string) {
	topics := []struct {
		angle, headline, data string
	}{
		{
			angle:    "FIRE math VN",
			headline: "Bao nhiêu tiền đủ để full-time trader ở VN? Monte Carlo nói thật",
			data: `Monte Carlo simulation 1,000 paths, horizon 10 năm, lạm phát 4%/năm, dùng VN30 buy-hold equal-weight (CAGR +15.37% 16 năm):

Vốn tối thiểu (P_fail ≤ 5%, real-terms):
- Chi phí 12tr/tháng: cần 2.55 tỷ
- Chi phí 24tr/tháng: cần 5.06 tỷ
- Chi phí 36tr/tháng: cần 7.78 tỷ
- Chi phí 48tr/tháng: cần 8.98 tỷ

Với 500tr vốn ban đầu:
- P(cháy 10y) với 12tr/tháng: 94%
- P(cháy 10y) với 24tr/tháng: 99%

Quy luật: vốn ≈ chi phí năm × 17-20 (safe withdrawal rate 5-6%/năm cho VN).

Sequence-of-returns risk: cùng 2 tỷ + rút 12tr/tháng, kết quả khác nhau theo năm bắt đầu:
- Bắt đầu 2010: sau 10y còn ~10 tỷ
- Bắt đầu 2022 (crisis): sau 10y còn ~2.2 tỷ (giảm 70%)

Vấn đề: KHÔNG phải skill — là timing luck.`,
		},
		{
			angle:    "Full-time vs part-time",
			headline: "29 tuổi, lương dev senior 60tr — nên full-time trade hay giữ job? Math nói gì",
			data: `Hypothetical profile: 29 tuổi, dev senior lương 60tr/tháng, 2.4 tỷ tài sản (1 tỷ tiết kiệm 8.5%, 800tr vàng, 550tr cổ phiếu, 25tr crypto), có gia đình 4 người ở chung cư.

Math:
- Income hộ: ~80tr/tháng (60tr chồng + 12-14tr vợ)
- Chi phí 4 người ở HCM, có chung cư: ~25tr/tháng
- Saving rate: ~50-60tr/tháng

Nếu QUIT job làm full-time trader:
- Mất 60tr/tháng income
- Cần lãi 60tr/tháng từ portfolio 2.4 tỷ = 2.5%/tháng = 30%/năm consistent
- VN30 buy-hold CAGR realistic +15%/năm
- Active trading sau phí: 3-7%/năm
- → Math KHÔNG work. Active strategy không cover lost income.

Path realistic:
- Giữ job, side stream consulting/SaaS 20-40tr/tháng
- Portfolio rebalance: 50% cổ phiếu VN30, 25% tiết kiệm, 10-15% vàng, 5-10% crypto
- 7-10 năm sau, vốn đạt 8+ tỷ → có buffer cho full-time

500tr-2.5 tỷ KHÔNG đủ. 5-8 tỷ là realistic FIRE level cho VN gia đình 4 người mức sống trung lưu.`,
		},
		{
			angle:    "Side income mathematics",
			headline: "Side income trong khi giữ job — dev VN nên build gì? Math + realistic ROI",
			data: `Phân loại side income theo TYPE × time-to-revenue × ceiling:

Type 1: Pure labor (sức người 100%)
- Freelance Upwork/Toptal: $30-50/h, ceiling ~20-40 tr/tháng (10h/tuần)
- Time to revenue: 1-3 tháng
- Risk: low, tax-friendly

Type 2: Trend-leveraged
- AI consultant 2024-2026: $80-150/h, ceiling 50-150 tr/tháng
- Time to revenue: 2-6 tháng (cần thương hiệu)
- Risk: medium, chu kỳ kết thúc thì collapse

Type 3: AI + bot leverage
- SaaS micro: 0-200tr/tháng, xác suất thành công 10-15%
- Affiliate/content: 2-15 tr/tháng sau 12 tháng
- Trading bots: KHÔNG phải income source. Chỉ là risk reduction tool.
- Time to revenue: 6-24 tháng
- Risk: high variance, high ceiling

So với portfolio passive income:
- VN30 buy-hold 1 tỷ → ~150tr/năm (CAGR 15%) = 12.5tr/tháng
- Bằng ~1/3 lương developer trung bình

Insight: với vốn < 5 tỷ, side income > passive income. Compound 3 nguồn (job + side + portfolio growth) đạt FIRE nhanh hơn chỉ passive.`,
		},
		{
			angle:    "Trader vs Holder mindset",
			headline: "Holder vs Trader — 16 năm data nói: 1 bên thua phí, 1 bên thua tâm lý",
			data: `Backtest 16 năm (2010-2026), 30 mã VN30+HNX30, phí TPS retail realistic:

Holder (buy-hold VN30 equal-weight, không rebalance):
- CAGR: +15.37%, MaxDD: -36%
- Total trades: 1 lần mua đầu
- Phí: ~0.35% one-time

Trader (active monthly rebalance, low-vol top-5):
- CAGR 10y window (2017-2026 cherry-pick): +26.2%
- CAGR 16y full: +3.0% (sau phí) ← TỆ HƠN HOLDER
- Total trades: ~12/năm × 16 = 192 trades
- Phí: ~0.35% × 192 = 67% tổng phí accumulated trên capital

Bài học:
1. Cherry-picked window làm active trading trông tốt. Full 16y window cho thấy phí ăn alpha.
2. Holder mindset cần kỷ luật khi -30% (2022) — KHÔNG bán panic.
3. Trader mindset cần model edge sau phí + kỷ luật execution.

Đa số retail VN thua vì:
- Là trader nhưng không có edge thực sự
- Trade khi cảm xúc (sell đáy, buy đỉnh)
- Pattern này được chứng minh: SELL/CUT_LOSS audit 0/3 đúng trong 30-day plan retrospective.`,
		},
	}
	slugs := make([]string, len(topics))
	for i, c := range topics {
		slugs[i] = "career-" + safeSlug(c.angle)
	}
	t := topics[pickFresh(slugs)]
	data = t.data
	title = t.headline
	slug = "career-" + safeSlug(t.angle)
	return
}

// ---------- T5: supply-chain ----------

func buildSupplyChain(ctx context.Context) (data, title, slug string) {
	cases := []struct {
		ticker, headline, data string
	}{
		{
			ticker:   "HPG",
			headline: "HPG phụ thuộc gì? Phân tích chuỗi tác động ngoại sinh của thép VN",
			data: `HPG (Hòa Phát Group) — nhà sản xuất thép lớn nhất VN. Phụ thuộc 3 input chính:

Input 1: Quặng sắt
- Giá tham chiếu: BHP iron ore (NYSE: BHP) hoặc Singapore TSI 62%
- 2-year Pearson r với HPG: r ≈ 0.55-0.65 (validated)
- Vd: BHP +14% 20d → HPG fwd 20d expected +7-9% (do r × move)

Input 2: Giá thép TQ (HRC futures Shanghai)
- r ≈ 0.45-0.55 với HPG
- TQ là price-setter cho khu vực châu Á

Input 3: USD/VND
- USD lên → quặng sắt nhập đắt hơn (negative cho HPG margin)
- Nhưng cũng → giá thép xuất tăng theo USD (positive cho revenue)
- Net r ≈ 0.10-0.20 (yếu)

Output:
- 60-70% nội địa (xây dựng VN)
- 30-40% xuất khẩu (chủ yếu ĐNÁ)

Sensitivity (theo lịch sử 2y backtest):
- Quặng -10% → HPG fwd 30d ~-5.5%
- HRC Shanghai +15% → HPG fwd 30d ~+6.7%
- USD/VND +5% → HPG fwd 30d ~+1.0% (yếu)

Insight cho retail: chỉ nhìn chart HPG là thiếu. Phải nhìn BHP + HRC TQ trước khi quyết định.`,
		},
		{
			ticker:   "VCB",
			headline: "VCB và Fed rate — bank lớn nhất VN bị ảnh hưởng bởi macro Mỹ thế nào?",
			data: `VCB (Vietcombank) — bank thương mại lớn nhất VN. Inputs:

Input 1: Fed funds rate (DFF)
- Fed tăng → NHNN VN thường tăng theo (lag 1-3 tháng) để giữ USD/VND
- VCB margin tăng khi rate tăng (asset re-price nhanh hơn deposit)
- 5-year r ≈ 0.30-0.45

Input 2: VN credit growth + tín dụng cá nhân
- NHNN target credit growth ~14-15%/năm
- VCB thường nằm trong nhóm growth cao hơn average
- Direct correlation với GDP VN

Input 3: NPL ratio + provision
- Lag indicator. BĐS sập → NPL tăng 6-12 tháng sau
- Vd: 2022 VTP → bank NPLs tăng 2023-2024

Output:
- 100% nội địa
- Customer mix: 60% retail + corporate, 40% large corp/SOE

Sensitivity:
- Fed +25bps → VCB fwd 60d ~+2-3% (positive)
- VN credit growth +1pp → +1-2% earnings/year ahead
- NPL +0.5pp → -5% to -8% earnings (bank-wide)

Cohort match VCB historically:
- VCB trong CRISIS regime (2022): -25% peak-to-trough
- Recovery 18 tháng để về cũ
- Volatility lower than HPG/FPT (β market ~0.85)`,
		},
		{
			ticker:   "VHM",
			headline: "VHM và BĐS — chuỗi tác động khi lãi suất + tín dụng BĐS thay đổi",
			data: `VHM (Vinhomes) — phát triển BĐS nhà ở lớn nhất VN. Inputs:

Input 1: Lãi suất cho vay mua nhà
- Khách hàng leverage 70-80% giá nhà
- Lãi tăng 1pp → demand giảm ~10-15% trong 6-12 tháng

Input 2: Tín dụng BĐS từ ngân hàng
- NHNN cap room cho vay BĐS — thay đổi mỗi năm
- 2022: NHNN siết → VHM stock -50%

Input 3: Disposable income hộ gia đình VN
- GDP per capita growth: +6-7%/năm
- Direct correlation với "trung lưu mua nhà"

Output:
- 100% nội địa, 95% Hà Nội + HCM + Quảng Ninh
- Margin: 30-35% gross (BĐS VN cao so với global)
- Cash flow: chu kỳ 3-5 năm/dự án

Sensitivity historical:
- Lãi vay +1pp → VHM fwd 60d ~-8% to -12%
- VTP-style scandal (10/2022) → VHM -40% trong 3 tháng

Cohort match:
- VHM trong CRISIS_FUNDAMENTAL (2022 BĐS sập): -65% peak-to-trough
- Tới 2026 vẫn chưa về cũ (3+ năm)
- Đây là loại CRISIS retail KHÔNG nên buy-the-dip. Khác CRISIS_PANIC.`,
		},
		{
			ticker:   "MWG",
			headline: "MWG (Thế Giới Di Động) — phụ thuộc gì? Phân tích retail giant VN",
			data: `MWG — chuỗi bán lẻ lớn nhất VN. Inputs:

Input 1: Disposable income hộ gia đình
- Sản phẩm bán: điện thoại, điện máy, FMCG (BHX)
- Demand co giãn theo thu nhập

Input 2: USD/VND
- 80% sản phẩm nhập từ TQ + Hàn + Mỹ
- USD lên → margin giảm hoặc giá bán tăng

Input 3: E-commerce competition (Shopee/Lazada/TikTok Shop)
- 2020-2024: market share của TGDD giảm dần với điện thoại (Shopee bán cheaper)
- BHX (Bách Hóa Xanh) khác: physical FMCG, ít threat từ e-com

Strategic shift 2022-2024:
- Đóng store thua lỗ
- Focus BHX expansion
- Restructure margin pressure

Sensitivity:
- VN GDP per capita +5% → MWG fwd 12m ~+8-12%
- USD/VND +5% → MWG fwd 6m ~-3-5%
- Competitor pricing pressure → ongoing margin compression

Performance backtest:
- 2014 IPO @ ~20,000 → 2026 ~75,000 = +275% (CAGR +11.7%)
- Below VN30 average (+15.4% CAGR)
- 2022 crash: -55% peak-to-trough — lớn hơn average VN30

Insight: MWG là proxy của VN consumer middle-class. Theo dõi GDP, USD/VND, e-com share dynamics.`,
		},
	}
	slugs := make([]string, len(cases))
	for i, c := range cases {
		slugs[i] = "supply-chain-" + strings.ToLower(c.ticker)
	}
	t := cases[pickFresh(slugs)]
	data = t.data
	title = t.headline
	slug = "supply-chain-" + strings.ToLower(t.ticker)
	return
}

// ---------- T6: psychology ----------

func buildPsychology(ctx context.Context) (data, title, slug string) {
	cases := []struct {
		angle, headline, data string
	}{
		{
			angle:    "Loss aversion",
			headline: "Tại sao bạn bán đáy? Loss aversion + cohort data VN",
			data: `Loss aversion (Kahneman & Tversky): nỗi đau mất 1 đồng = niềm vui được 2 đồng.

Behavioral pattern trong retail VN:
- Mua A @ 30, A xuống 24 (-20%) → bán panic
- Mua B @ 30, B lên 36 (+20%) → bán chốt lời
- Hold thua, run thắng — ngược lại "let winners run, cut losses"

Backtest cohort cho thấy retail panic-sell là sai:

CRISIS regime (vùng panic-driven):
- Fwd 60d: +8.53%, win 67%, edge +3.78% (highest)
- Bán đáy = thua 33% xác suất hồi phục.

Plan audit 30 ngày:
- BUY/HOLD recommendations: 14/15 đúng (93%)
- SELL/CUT_LOSS recommendations: 0/3 đúng (panic-sell trong oversold)

Math behind:
- Loss aversion = emotional, không phải rational.
- Probability-weighted decision: kỳ vọng dương trong CRISIS_PANIC + RSI<30.
- Action correct: rebalance vào (mua thêm), không bán.

Cách defeat loss aversion:
1. Set rule trước: "Tôi không bán dưới -25% trừ khi MA200 đảo chiều bearish"
2. Pre-commit: viết investment thesis lúc bình tĩnh, đọc lại lúc panic.
3. Time pressure giảm: nếu tự thấy "phải bán NGAY", đợi 24h. 90% trường hợp panic giảm.`,
		},
		{
			angle:    "Recency bias",
			headline: "Recency bias: tại sao bạn nghĩ thị trường năm tới giống năm vừa rồi",
			data: `Recency bias: weight unfairly các sự kiện gần đây so với base rate dài hạn.

Pattern trong retail VN:
- 2020-2021 (bull run + COVID money): retail đổ tiền cuối 2021, đỉnh cuối Q1/2022
- 2022 crash (-30%): retail bán đáy Q3-Q4/2022
- 2023-2024 recovery: retail tham gia trễ, mid-2024
- 2025 mini-crash: retail panic trở lại

Data backtest 16 năm:
- 10/16 năm xanh, 6/16 đỏ — base rate xanh 62.5%
- Sau năm đỏ -20%+, năm tiếp theo: avg +18% (5/5 lần trong sample)
- Sau năm xanh +30%+, năm tiếp theo: avg +6% (3/4 dương)

Retail behavior reverse:
- Sau năm đỏ → "thị trường còn xấu" → giữ tiền mặt → bỏ lỡ +18%
- Sau năm xanh → "tiếp tục bull" → vào nhiều → return chỉ +6%

Math:
- Investing với base rate (mua đều mỗi tháng VN30, ignore last year): CAGR +15.37%
- Investing với recency-biased timing: typically -3% đến -5%/năm (mistimed entries)

Cách defeat:
1. Dollar-cost averaging (mua đều mỗi tháng) — bypass timing
2. Reset tham chiếu: nhìn 10-year chart trước khi nhìn 1-year
3. Rebalance theo target allocation, không theo "feel"`,
		},
		{
			angle:    "Herd behavior",
			headline: "Tại sao mã \"hot\" trên FB group thường là mã bạn nên tránh",
			data: `Herd behavior: cá nhân quyết định dựa trên đám đông thay vì analysis độc lập.

Pattern trong VN FB groups đầu tư:
- Mã được "phím" hot thường đã chạy +40-80% trong 3-6 tháng trước
- Retail tham gia ở đỉnh local
- Crash 20-40% sau khi smart money exit
- Group im lặng/đổ lỗi

Data anchor (cohort):
- Stocks với 60d return > +50% (vùng nóng nhất):
  - Fwd 60d sau đó: trung bình +2-3% (giảm so với baseline +4.75%)
  - Volatility 60d sau: +50-80% so với baseline
  - Drawdown probability: 35% chance giảm >20% trong 90d

- Stocks với 60d return -20% đến -10% (vùng "buồn"):
  - Fwd 60d: +5-8% (cao hơn baseline)
  - Volatility lower
  - Recovery probability cao

Insight ngược:
- Mã "hot" trên group = mã đã chạy = late mover = under-perform forward.
- Mã "buồn", ít ai nói = oversold = thường recovery.

Bài học:
- Trước khi mua mã "hot", check 60d return. Nếu >40%, đợi pullback hoặc skip.
- Smart money mua trước khi hot, retail mua sau khi hot. Đừng là 2nd group.
- Group đầu tư FB là LAGGING indicator, không phải leading.`,
		},
	}
	slugs := make([]string, len(cases))
	for i, c := range cases {
		slugs[i] = "psychology-" + safeSlug(c.angle)
	}
	t := cases[pickFresh(slugs)]
	data = t.data
	title = t.headline
	slug = "psychology-" + safeSlug(t.angle)
	return
}

// ---------- T7: data-insight ----------

func buildDataInsight(ctx context.Context) (data, title, slug string) {
	topics := []struct {
		angle, headline, data string
	}{
		{
			angle:    "Wyckoff stage upset",
			headline: "5 điều cohort 47K data point cho thấy ngược trực giác retail VN",
			data: `Backtest cohort VN30+HNX30, N=45,493 row clean, 9 năm:

5 surprising findings:

1. Wyckoff stage 3 (distribution) thắng nhất, không phải stage 2 (markup)
- Stage 3 fwd 60d: +6.63%, win 60%, edge +1.88%
- Stage 2 fwd 60d: +5.46%, win 53%, edge +0.71%
- Retail nghĩ "đỉnh phân phối = bán". Data nói khác.

2. RSI > 70 (overbought) cho fwd return cao nhất
- RSI>70: +9.23%, win 60%
- RSI<30: +6.32%, win 60%
- Cả 2 đều thắng, RSI cao thắng nhiều hơn.

3. CRISIS regime = highest forward return
- CRISIS: +8.53%, win 67%
- STABLE: +3.20%, win 48%
- Tâm lý retail panic-sell trong CRISIS là bias chính.

4. Active trading thua buy-hold 16 năm
- Active low-vol top-5 (rebalance hàng tháng): CAGR 16y +3%
- Buy-hold VN30 equal-weight: CAGR 16y +15.37%
- Phí 0.35% × 12 rebalances/năm = 4.2%/năm ăn alpha.

5. Sequence-of-returns risk lớn hơn skill
- 2 tỷ + rút 12tr/tháng, start 2010: sau 10y còn 10 tỷ
- Cùng setup start 2022: còn 2.2 tỷ (giảm 70%)
- Không phải skill — là timing luck. 50% kết quả phụ thuộc năm bắt đầu.

Common thread: tâm lý + intuition của retail VN ngược lại với những gì data cohort cho thấy. Đây là edge có thể khai thác — không phải bằng "đảo ngược tâm lý" mù quáng, mà bằng pre-commit rules dựa trên cohort.`,
		},
		{
			angle:    "Joint factor analysis",
			headline: "1 chỉ báo vô dụng, 2 chỉ báo có meaning — joint factor matters",
			data: `Cohort analysis cho thấy: edge của signal tăng dramatically khi kết hợp 2 chỉ báo.

Single factor edges (60d fwd):
- RSI > 70: +4.48% so với baseline
- Uptrend (MA200 + MA50): +1.20%
- MACD bullish: +0.65%

Joint factor edges (cùng RSI + MA trend):
- Uptrend × RSI > 70: +5.05% ← gấp 4× MA trend alone, gấp 1.1× RSI alone
- Uptrend × RSI 50-70: +1.45% (so với RSI 50-70 alone -0.07%)
- Downtrend × RSI > 70: -6.15% ← signal SELL duy nhất

Insight:
- 1 indicator alone = noise
- 2 orthogonal indicators (momentum × trend) = signal
- 3+ indicators bắt đầu over-fit. Stick 2.

Anti-pattern phổ biến trong VN cộng đồng:
- "RSI > 70 = bán" (indicator alone, không context)
- "Phá MA200 = mua" (indicator alone, không momentum)
- "Volume tăng + giá tăng = mua" (volume noisy alone)

Phải hỏi: trend context + momentum signal phù hợp với hành động không?

Best combos cohort cho thấy:
1. MA trend (close vs MA200 vs MA50) × RSI bucket → edge ±5%
2. Regime (CRISIS/STABLE/EUPHORIA) × MA trend → edge ±3-4%
3. Wyckoff stage × volume class → edge ±2%

Trade nhiều hơn 2 conditions = giảm sample size, tăng phí, giảm edge thực.`,
		},
		{
			angle:    "Window-dependent results",
			headline: "Backtest cherry-picking: cùng strategy, 2 con số khác biệt 9× — bài học từ window selection",
			data: `Cùng chiến lược "low-vol top-5 VN30 monthly rebalance", 2 windows:

Window A (2017-01 → 2026-04, 9.3 năm):
- CAGR: +26.2%
- MaxDD: -42%
- Sharpe: 0.74

Window B (2010-01 → 2026-04, 16.3 năm):
- CAGR: +3.0% (sau khi fix data scale issue)
- MaxDD: -36% (buy-hold benchmark)
- Same fees

Khác biệt: Window B bao gồm 2011 (-10%), 2018 (-15%) — 2 năm giảm mà active strategy không né được, và phí tích lũy ăn alpha.

Window cherry-picked phổ biến:
- Tech 2010-2020: bull run liên tục → ai cũng wins
- Crypto 2020-2021: 5× chỉ trong 1 năm → "lý thuyết X works"
- VN 2017-2026: 8/9 năm xanh → momentum strategy looks magic

Rule:
- Bất kỳ backtest nào nên span ≥ 2 chu kỳ kinh tế (~15 năm)
- Bất kỳ "10y CAGR > buy-hold benchmark" trong window được pick: skeptical
- Skin in the game: backtest qua 2008-2009 (US) hoặc 2022-2023 (VN BĐS) trước khi tin

Áp dụng:
- Marketing material show "10y returns": hỏi why 10y, không 20y?
- Newsletter "AI fund 5y CAGR 35%": hỏi why not since inception?
- Bot Telegram "win rate 78%": hỏi N sample, time horizon, fees.`,
		},
	}
	slugs := make([]string, len(topics))
	for i, c := range topics {
		slugs[i] = "insight-" + safeSlug(c.angle)
	}
	t := topics[pickFresh(slugs)]
	data = t.data
	title = t.headline
	slug = "insight-" + safeSlug(t.angle)
	return
}

// ---------- T8: comparison ----------

func buildComparison(ctx context.Context) (data, title, slug string) {
	cases := []struct {
		pair, headline, data string
	}{
		{
			pair:     "VCB-vs-CTG",
			headline: "VCB vs CTG: 2 bank lớn nhất VN — ai thắng 10 năm qua, và lý do",
			data: `VCB (Vietcombank) vs CTG (Vietinbank) — 2 SOE bank lớn nhất VN.

10-year performance (2016-2026):
- VCB: 24,000 → 90,000 VND adjusted (~+275%, CAGR +14.1%)
- CTG: 12,000 → 35,000 VND adjusted (~+192%, CAGR +11.3%)

Backtest 16 năm (2010-2026):
- VCB CAGR: ~+13.5%
- CTG CAGR: ~+11.95%

Khác biệt structural:
- VCB: bank trade hàng đầu, tỷ trọng FX + xuất nhập khẩu cao. Margin tốt nhất ngành.
- CTG: bank truyền thống tín dụng SOE, margin trung bình.

Lý do VCB thắng:
1. ROE consistent 18-20%, CTG 12-15%
2. NPL ratio VCB ~0.5-0.8%, CTG ~1.2-1.8%
3. Premium valuation (P/B 3.5-4× vs CTG 1.8-2.2×)

Insight cho retail:
- "Mua bank rẻ" (CTG) không thắng "mua bank quality" (VCB) trên 10y+
- Quality premium thực sự (không phải chỉ là perception)
- 2022 crash: VCB -25%, CTG -30%. VCB recovery nhanh hơn 4-6 tháng.

Cohort match hiện tại (cần check live):
- VCB nếu uptrend × RSI 50-70: cohort edge +1.45%
- CTG nếu downtrend × RSI 50-70: cohort edge -2.38%
- Mã nào ở cohort tốt hơn, đó là choice.

Khuyến nghị reproducible: chạy lmcli rate VCB và lmcli rate CTG, so sánh overall gauge.`,
		},
		{
			pair:     "HPG-vs-HSG",
			headline: "HPG vs HSG: ai là vua thép VN? Backtest 10 năm + sensitivity test",
			data: `HPG (Hòa Phát) vs HSG (Hoa Sen) — 2 nhà sản xuất thép lớn nhất VN.

10-year performance (2016-2026):
- HPG: 8,000 → 27,000 VND adjusted (~+237%, CAGR +12.9%)
- HSG: 15,000 → 22,000 VND adjusted (~+47%, CAGR +3.9%)

Big difference: HPG có lò cao tích hợp (vertical integration), HSG chỉ cán nguội + sơn.

Sensitivity với commodity:
- Quặng sắt (BHP) tăng 10%:
  - HPG fwd 30d: ~+5-7% (positive — HPG hưởng spread vertical)
  - HSG fwd 30d: ~-3-5% (negative — chi phí input tăng, không vertical)
- Giá thép HRC Shanghai tăng 15%:
  - HPG: ~+8%
  - HSG: ~+12% (revenue elastic hơn, không có buffer integration)

Macro position:
- HPG: net exporter, hưởng USD/VND lên
- HSG: net importer, mất USD/VND lên

Margin profile:
- HPG: gross 18-25%, net 10-12%
- HSG: gross 9-13%, net 3-6%

2022 crash (BĐS sập, demand thép VN -20%):
- HPG: -55% peak-to-trough
- HSG: -65% peak-to-trough

Recovery 2023-2024:
- HPG: +90% từ đáy
- HSG: +40% từ đáy

Insight: trong industry cyclical (thép), vertical integration = moat. HPG là proxy cho thép VN structural; HSG là leveraged play (high reward, high risk).`,
		},
		{
			pair:     "VNM-vs-MSN",
			headline: "VNM vs MSN: sữa cổ điển vs holding đa ngành — ai phù hợp retail?",
			data: `VNM (Vinamilk) vs MSN (Masan Group) — 2 consumer giants VN.

10-year (2016-2026):
- VNM: 80,000 → 58,000 VND adjusted (~-27%, CAGR -3.1%)
- MSN: 30,000 → 78,000 VND adjusted (~+160%, CAGR +10.0%)

VNM 10y NEGATIVE — sốc với retail VN coi đây là "blue chip ổn định".

Lý do VNM giảm:
- Sữa VN market saturated, growth 2-4%/năm
- Competition từ TH, Nestle, Vinasoy tăng
- Margin pressure từ giá feed lên (USD)

Lý do MSN tăng:
- 2014-2022: thâu tóm Vinacafé, Phúc Long, WinMart, Techcombank stake, etc.
- Đa ngành = đa risk = đa upside
- 2022 crash: -50% nhưng recovery nhanh

Risk profile:
- VNM: low volatility, predictable cash flow, dividend 5-7%/năm. KÊN — KÊN. Performance: ÂM 10y.
- MSN: high volatility, M&A driven, dividend yield thấp. Performance: +160% 10y.

Insight cho retail:
- "Mua VNM ngủ ngon" — kiểu khuyên cũ phổ biến trong VN gia đình. 10 năm: VỠ.
- VNM là cảnh báo: "blue chip an toàn" ≠ "đầu tư tốt"
- Define "tốt" rõ ràng:
  - Vốn protection? VN30 ETF + Treasury, không phải single stock
  - Growth? MSN-style holding, hoặc tech (FPT)
  - Income? Dividend ETF, KHÔNG phải mua single stock yield

VN30 buy-hold equal-weight CAGR 16y +15.37% beat cả VNM và MSN solo.`,
		},
	}
	slugs := make([]string, len(cases))
	for i, c := range cases {
		slugs[i] = "compare-" + safeSlug(c.pair)
	}
	t := cases[pickFresh(slugs)]
	data = t.data
	title = t.headline
	slug = "compare-" + safeSlug(t.pair)
	return
}

// ---------- T9: regime-now ----------

func buildRegimeNow(ctx context.Context) (data, title, slug string) {
	stocks, _ := fetchers.VPSMultiple(ctx, vn30)
	flow := signalsFlowSafe(stocks)
	globals := fetchers.YahooMultiple(ctx, []string{"^VIX", "^GSPC", "^N225", "^HSI"})

	var sb strings.Builder
	sb.WriteString("Today's market state (live snapshot):\n\n")
	fmt.Fprintf(&sb, "Domestic flow: %s (buy pressure %.1f%%)\n", flow.Signal, flow.BuyPressure)
	for _, g := range globals {
		fmt.Fprintf(&sb, "Global %s: %.2f (%+.2f%%)\n", g.Symbol, g.Price, g.ChangePct)
	}
	sb.WriteString("\nCohort context per regime (historical 9y):\n")
	sb.WriteString("- CRISIS: +8.53% fwd 60d, win 67%, edge +3.78%\n")
	sb.WriteString("- EUPHORIA: +6.12%, 59%, +1.37%\n")
	sb.WriteString("- VOLATILE: +6.34%, 56%, +1.59%\n")
	sb.WriteString("- STABLE: +3.20%, 48%, -1.55% (baseline +4.75%)\n\n")
	sb.WriteString("Regime classification rules (deterministic, no AI):\n")
	sb.WriteString("- CRISIS: VIX > 30 OR VN-Index -3%+ trong 1 phiên OR tier-1 news flag\n")
	sb.WriteString("- EUPHORIA: VIX < 15 AND VN-Index uptrend 60d AND breadth > 70%\n")
	sb.WriteString("- VOLATILE: VIX 20-30 OR contradicting global signals\n")
	sb.WriteString("- STABLE: default state, VIX 15-20\n\n")
	sb.WriteString("Sub-classification CRISIS:\n")
	sb.WriteString("- CRISIS_PANIC: contagion từ ngoài (VIX spike, global tier-1 news). Historically: recover 3-12 tháng.\n")
	sb.WriteString("- CRISIS_FUNDAMENTAL: vấn đề cấu trúc nội tại VN. Vd VTP 10/2022. May not recover.\n")

	data = sb.String()
	title = "Thị trường hôm nay (" + time.Now().Format("02/01") + ") đang ở regime nào? Phân loại deterministic + ý nghĩa"
	slug = "regime-now-" + time.Now().Format("01-02")
	return
}

// ---------- T10: news-impact ----------

func buildNewsImpact(ctx context.Context) (data, title, slug string) {
	globals := fetchers.YahooMultiple(ctx, []string{"^GSPC", "^DJI", "^IXIC", "^N225", "^HSI", "^VIX", "GC=F", "CL=F"})

	var sb strings.Builder
	sb.WriteString("Global market state past day:\n\n")
	for _, g := range globals {
		fmt.Fprintf(&sb, "%s: %.2f (%+.2f%% từ phiên trước, time %s)\n", g.Symbol, g.Price, g.ChangePct, g.MarketTime.Format("02/01"))
	}

	sb.WriteString("\nLịch sử cohort ảnh hưởng global → VN:\n")
	sb.WriteString("- S&P 500 -2%+ trong 1 phiên → VN-Index thường -1.5% phiên sau (lag 1 ngày, contagion via foreign flow)\n")
	sb.WriteString("- VIX > 25 → VN regime thường shift sang VOLATILE/CRISIS_PANIC\n")
	sb.WriteString("- Vàng (GC=F) +5%+ 1 tuần → safe haven mode, VN bank/BĐS thường yếu\n")
	sb.WriteString("- Dầu (CL=F) +10%+ 1 tháng → GAS/PLX hưởng lợi, MWG/VHM bất lợi (input cost)\n")
	sb.WriteString("- Hang Seng -3%+ → tâm lý APAC tiêu cực, VN30 cũng giảm theo (correlation 0.4-0.5)\n\n")

	sb.WriteString("Forward implication based on cohort:\n")
	sb.WriteString("- CRISIS_PANIC từ contagion ngoài: VN fwd 60d +8.53%, win 67% — vùng historically là buying opportunity.\n")
	sb.WriteString("- EUPHORIA global (VIX < 13, S&P all-time high): VN fwd thường +6.12%, win 59% — momentum continuation.\n")

	data = sb.String()
	title = "Tin nóng global (" + time.Now().Format("02/01") + ") — ảnh hưởng VN30 thế nào? Cohort lookup"
	slug = "news-impact-" + time.Now().Format("01-02")
	return
}

// ============================================================================
// Helpers
// ============================================================================

func safeSlug(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "/", "-")
	s = strings.ReplaceAll(s, ">", "")
	s = strings.ReplaceAll(s, "<", "")
	s = strings.ReplaceAll(s, "(", "")
	s = strings.ReplaceAll(s, ")", "")
	s = strings.ReplaceAll(s, "đ", "d")
	s = strings.ReplaceAll(s, "Đ", "d")
	// Strip non-ascii
	out := strings.Builder{}
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			out.WriteRune(r)
		}
	}
	res := out.String()
	if len(res) > 40 {
		res = res[:40]
	}
	return strings.Trim(res, "-")
}

// Keep these from previous version for any other callers
type domesticFlowResult struct {
	Signal      string
	BuyPressure float64
}

func signalsFlowSafe(stocks []types.StockData) domesticFlowResult {
	totalBid, totalAsk := int64(0), int64(0)
	for _, s := range stocks {
		totalBid += s.BidVol
		totalAsk += s.AskVol
	}
	total := totalBid + totalAsk
	if total == 0 {
		return domesticFlowResult{Signal: "không xác định (markets có thể đóng)", BuyPressure: 50}
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

// Stub to satisfy import — unused after refactor but keeps compatibility
func init() {
	_ = sort.Slice
	_ = fmt.Sprintf
}
