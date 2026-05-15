// HTML landing-page generator. Output is a single self-contained file:
// inline CSS, no JS, no external deps. Designed for GitHub Pages static
// serving so it loads fast on mobile.
package main

import (
	"context"
	"fmt"
	"html"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/ducnhd/lotusmarket/go/fetchers"
	"github.com/ducnhd/lotusmarket/go/market"
	"github.com/ducnhd/lotusmarket/go/signals"
	"github.com/ducnhd/lotusmarket/go/types"
)

func writeFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

func stderrW() io.Writer { return os.Stderr }

const htmlCSS = `
:root {
  --bg: #0f1419;
  --bg-card: #1a2128;
  --bg-card-hover: #232c35;
  --text: #e6edf3;
  --text-dim: #8b949e;
  --accent: #7ee8a3;
  --accent-dim: #4d8c66;
  --red: #ff7b72;
  --red-dim: #5c3a3a;
  --green: #7ee8a3;
  --green-dim: #3a5c46;
  --yellow: #ffd16a;
  --border: #2a323c;
  --radius: 10px;
}
@media (prefers-color-scheme: light) {
  :root {
    --bg: #f6f8fa;
    --bg-card: #ffffff;
    --bg-card-hover: #f0f3f6;
    --text: #1f2328;
    --text-dim: #59636e;
    --accent: #1a7f4e;
    --border: #d0d7de;
  }
}
* { box-sizing: border-box; margin: 0; padding: 0; }
body {
  background: var(--bg); color: var(--text);
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", "SF Pro Text", system-ui, sans-serif;
  font-size: 15px; line-height: 1.55;
  -webkit-font-smoothing: antialiased;
}
.container { max-width: 1080px; margin: 0 auto; padding: 24px 20px 60px; }
header.hero {
  padding: 56px 0 36px; text-align: center;
  background: linear-gradient(180deg, rgba(126,232,163,0.06) 0%, transparent 100%);
  border-bottom: 1px solid var(--border);
  margin-bottom: 32px;
}
h1.title { font-size: 44px; font-weight: 800; letter-spacing: -0.02em; margin-bottom: 6px; }
h1.title .lotus { display: inline-block; margin-right: 10px; }
.tagline { color: var(--text-dim); font-size: 17px; margin-bottom: 24px; }
.cta { display: flex; gap: 10px; justify-content: center; flex-wrap: wrap; }
.cta a {
  display: inline-flex; align-items: center; gap: 6px;
  padding: 9px 18px; border-radius: 999px; text-decoration: none;
  font-weight: 600; font-size: 14px;
  border: 1px solid var(--border); color: var(--text);
  background: var(--bg-card); transition: background 0.15s;
}
.cta a:hover { background: var(--bg-card-hover); }
.cta a.primary { background: var(--accent); color: #0f1419; border-color: var(--accent); }
.cta a.primary:hover { background: var(--green); }
section { margin: 48px 0; }
section h2 {
  font-size: 22px; font-weight: 700; margin-bottom: 18px;
  display: flex; align-items: center; gap: 10px;
}
section h2 .updated { font-size: 12px; color: var(--text-dim); font-weight: 400; margin-left: auto; }
.grid { display: grid; gap: 12px; }
.grid-3 { grid-template-columns: repeat(auto-fit, minmax(240px, 1fr)); }
.grid-4 { grid-template-columns: repeat(auto-fit, minmax(180px, 1fr)); }
.grid-2 { grid-template-columns: repeat(auto-fit, minmax(320px, 1fr)); }
.card {
  background: var(--bg-card); border: 1px solid var(--border);
  border-radius: var(--radius); padding: 16px;
}
.card.dense { padding: 12px 14px; }
.metric-label { font-size: 12px; color: var(--text-dim); text-transform: uppercase; letter-spacing: 0.05em; }
.metric-value { font-size: 24px; font-weight: 700; margin-top: 4px; }
.metric-sub { font-size: 13px; color: var(--text-dim); margin-top: 2px; }
.pct { font-weight: 600; }
.pct.up { color: var(--green); }
.pct.down { color: var(--red); }
.row { display: flex; justify-content: space-between; align-items: center; padding: 6px 0; border-bottom: 1px solid var(--border); }
.row:last-child { border-bottom: none; }
.row .name { font-weight: 500; }
.row .meta { font-size: 12px; color: var(--text-dim); }
.ticker { font-family: "SF Mono", Menlo, Consolas, monospace; font-weight: 600; letter-spacing: 0.02em; }
.pulse-banner {
  background: linear-gradient(135deg, rgba(126,232,163,0.10), rgba(126,232,163,0.02));
  border: 1px solid var(--accent-dim);
  border-radius: var(--radius); padding: 20px 24px;
  display: flex; align-items: center; gap: 18px; flex-wrap: wrap;
}
.pulse-banner .pulse-num { font-size: 40px; font-weight: 800; line-height: 1; }
.pulse-banner .pulse-label { font-size: 13px; color: var(--text-dim); }
.pulse-banner .pulse-text { flex: 1; min-width: 220px; }
.pulse-banner .pulse-headline { font-size: 17px; font-weight: 600; margin-bottom: 2px; }
.feature-list { list-style: none; }
.feature-list li {
  padding: 10px 0; border-bottom: 1px solid var(--border);
  display: flex; gap: 12px; align-items: flex-start;
}
.feature-list li:last-child { border-bottom: none; }
.feature-list .icon { font-size: 18px; flex-shrink: 0; }
.feature-list .ftext { flex: 1; }
.feature-list .ftext strong { display: block; margin-bottom: 2px; }
.feature-list .ftext .desc { font-size: 13px; color: var(--text-dim); }
pre.code {
  background: var(--bg); border: 1px solid var(--border);
  border-radius: 6px; padding: 14px 16px; overflow-x: auto;
  font-family: "SF Mono", Menlo, Consolas, monospace;
  font-size: 13px; line-height: 1.5;
}
pre.code .kw { color: var(--accent); }
pre.code .str { color: var(--yellow); }
pre.code .cmt { color: var(--text-dim); }
footer {
  margin-top: 60px; padding: 24px 0; text-align: center;
  border-top: 1px solid var(--border); color: var(--text-dim); font-size: 13px;
}
footer a { color: var(--accent); text-decoration: none; }
footer a:hover { text-decoration: underline; }
.badges { display: inline-flex; gap: 6px; margin: 0 0 18px; flex-wrap: wrap; justify-content: center; }
.badge {
  font-size: 11px; padding: 3px 9px; border-radius: 999px;
  background: var(--bg-card); border: 1px solid var(--border); color: var(--text-dim);
}
.archive-list { list-style: none; }
.archive-list li { padding: 6px 0; }
.archive-list a { color: var(--accent); text-decoration: none; font-family: "SF Mono", monospace; font-size: 14px; }
.archive-list a:hover { text-decoration: underline; }
@media (max-width: 600px) {
  h1.title { font-size: 32px; }
  .container { padding: 16px 14px 40px; }
  section { margin: 32px 0; }
  .pulse-banner .pulse-num { font-size: 32px; }
}
`

func pctClass(v float64) string {
	if v > 0.05 {
		return "up"
	}
	if v < -0.05 {
		return "down"
	}
	return ""
}

func pctSign(v float64) string {
	if v >= 0 {
		return "+"
	}
	return ""
}

func runHTML(ctx context.Context, outPath string) {
	stocks, _ := fetchers.VPSMultiple(ctx, vn30)
	flow := signals.ComputeDomesticFlow(stocks)
	sectors := market.RankSectorsByFlow(stocks)
	globals := fetchers.YahooMultiple(ctx, []string{"^GSPC", "^DJI", "^IXIC", "^N225", "^HSI", "^VIX", "GC=F", "CL=F", "DX-Y.NYB"})

	// Top movers (top 5 each side)
	sorted := make([]types.StockData, len(stocks))
	copy(sorted, stocks)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].ChangePercent > sorted[j].ChangePercent })

	var sb strings.Builder
	sb.WriteString(`<!DOCTYPE html>
<html lang="vi">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<meta name="description" content="Vietnamese Stock Market Toolkit — daily auto-generated reports + open-source library for Go &amp; Python. VN30 / HNX30 / global markets / sector flow / star ratings.">
<meta property="og:title" content="lotusmarket — VN stock toolkit + daily reports">
<meta property="og:description" content="Free open-source Vietnamese stock market library + daily auto-generated reports. Real data, no hot picks.">
<meta property="og:type" content="website">
<meta property="og:url" content="https://ducnhd.github.io/lotusmarket/">
<link rel="canonical" href="https://ducnhd.github.io/lotusmarket/">
<title>lotusmarket — VN Stock Market Toolkit · Daily Reports</title>
<script type="application/ld+json">{"@context":"https://schema.org","@graph":[{"@type":"Organization","@id":"https://lotusai.servehttp.com/#org","name":"Lotus AI / lotusmarket","url":"https://lotusai.servehttp.com/","sameAs":["https://github.com/ducnhd/lotusmarket","https://pypi.org/project/lotusmarket/","https://t.me/vnlotusmarket","https://ducnhd.github.io/lotusmarket/"]},{"@type":"WebSite","url":"https://ducnhd.github.io/lotusmarket/","name":"lotusmarket daily reports","publisher":{"@id":"https://lotusai.servehttp.com/#org"},"inLanguage":"vi-VN"}]}</script>
<style>` + htmlCSS + `</style>
</head>
<body>
<header class="hero">
  <div class="container">
    <h1 class="title"><span class="lotus">🌸</span>lotusmarket</h1>
    <p class="tagline">Vietnamese Stock Market Toolkit · Free · Open Source · No API key</p>
    <div class="badges">
      <span class="badge">v0.5.0</span>
      <span class="badge">Go + Python</span>
      <span class="badge">MIT License</span>
      <span class="badge">No hot picks</span>
    </div>
    <div class="cta">
      <a class="primary" href="https://lotusai.servehttp.com/">🌸 Main site</a>
      <a href="https://github.com/ducnhd/lotusmarket">⭐ GitHub</a>
      <a href="https://pypi.org/project/lotusmarket/">📦 PyPI</a>
      <a href="https://t.me/vnlotusmarket">💬 Telegram</a>
      <a href="latest.html">📄 Today's report</a>
    </div>
  </div>
</header>

<main class="container">
`)

	// === SECTION: Today's snapshot ===
	now := time.Now().Format("02/01/2006 15:04")
	fmt.Fprintf(&sb, `<section id="today">
  <h2>📊 Hôm nay <span class="updated">cập nhật %s</span></h2>
  <div class="pulse-banner">
    <div>
      <div class="pulse-num">%.0f%%</div>
      <div class="pulse-label">áp lực mua</div>
    </div>
    <div class="pulse-text">
      <div class="pulse-headline">Dòng tiền nội địa: <strong>%s</strong></div>
      <div class="metric-sub">Tính từ tỷ lệ bid/ask volume trên 30 mã VN30 — không phải hot pick, chỉ deterministic signal từ raw data.</div>
    </div>
  </div>
</section>
`, html.EscapeString(now), flow.BuyPressure, html.EscapeString(flow.Signal))

	// === SECTION: Sector flow ===
	sb.WriteString(`<section id="sectors">
  <h2>🏭 Dòng tiền sector</h2>
  <div class="grid grid-3">
`)
	maxS := 6
	if len(sectors) < maxS {
		maxS = len(sectors)
	}
	for i := 0; i < maxS; i++ {
		s := sectors[i]
		fmt.Fprintf(&sb, `    <div class="card">
      <div class="metric-label">#%d %s</div>
      <div class="metric-value pct %s">%s%.2f%%</div>
      <div class="metric-sub">KL %.0f tỷ · %d/%d tăng</div>
    </div>
`, s.Rank, html.EscapeString(s.Sector), pctClass(s.AvgChangePct),
			pctSign(s.AvgChangePct), s.AvgChangePct,
			s.TotalVolumeVND/1e9, s.AdvancesCount, s.TotalCount)
	}
	sb.WriteString(`  </div>
</section>
`)

	// === SECTION: Top movers ===
	sb.WriteString(`<section id="movers">
  <h2>📈 Top tăng / 📉 Top giảm VN30</h2>
  <div class="grid grid-2">
    <div class="card">
      <div class="metric-label" style="margin-bottom:10px">📈 Top tăng</div>
`)
	for i := 0; i < 5 && i < len(sorted); i++ {
		s := sorted[i]
		fmt.Fprintf(&sb, `      <div class="row"><span class="ticker">%s</span><span class="meta">%.0f · KL %d</span><span class="pct up">%s%.2f%%</span></div>
`, html.EscapeString(s.Ticker), s.Close, s.Volume, pctSign(s.ChangePercent), s.ChangePercent)
	}
	sb.WriteString(`    </div>
    <div class="card">
      <div class="metric-label" style="margin-bottom:10px">📉 Top giảm</div>
`)
	for i := 0; i < 5 && i < len(sorted); i++ {
		s := sorted[len(sorted)-1-i]
		fmt.Fprintf(&sb, `      <div class="row"><span class="ticker">%s</span><span class="meta">%.0f · KL %d</span><span class="pct down">%s%.2f%%</span></div>
`, html.EscapeString(s.Ticker), s.Close, s.Volume, pctSign(s.ChangePercent), s.ChangePercent)
	}
	sb.WriteString(`    </div>
  </div>
</section>
`)

	// === SECTION: Global ===
	if len(globals) > 0 {
		sb.WriteString(`<section id="global">
  <h2>🌏 Thị trường thế giới</h2>
  <div class="grid grid-4">
`)
		nameByIdx := map[string]string{
			"^GSPC": "S&P 500", "^DJI": "Dow Jones", "^IXIC": "Nasdaq",
			"^N225": "Nikkei", "^HSI": "Hang Seng", "^VIX": "VIX",
			"GC=F": "Vàng", "CL=F": "Dầu WTI", "DX-Y.NYB": "USD Index",
		}
		for _, q := range globals {
			name, ok := nameByIdx[q.Symbol]
			if !ok {
				name = q.Symbol
			}
			fmt.Fprintf(&sb, `    <div class="card dense">
      <div class="metric-label">%s</div>
      <div class="metric-value">%.2f</div>
      <div class="metric-sub pct %s">%s%.2f%%</div>
    </div>
`, html.EscapeString(name), q.Price, pctClass(q.ChangePct), pctSign(q.ChangePct), q.ChangePct)
		}
		sb.WriteString(`  </div>
</section>
`)
	}

	// === SECTION: Features ===
	sb.WriteString(`<section id="features">
  <h2>🧰 Library features</h2>
  <div class="grid grid-2">
    <div class="card">
      <ul class="feature-list">
        <li><span class="icon">📡</span><div class="ftext"><strong>Fetchers</strong><div class="desc">VPS real-time · Entrade history · KBS fundamentals · Cafef insider · VCI dividends · Yahoo global · FRED macro — all free, no API key needed (FRED optional).</div></div></li>
        <li><span class="icon">📐</span><div class="ftext"><strong>Technical analysis</strong><div class="desc">RSI · MA20/50/200 · MACD · Bollinger Bands · ATR · momentum · weekly aggregation · signal + score 0-100.</div></div></li>
        <li><span class="icon">⭐</span><div class="ftext"><strong>Star ratings</strong><div class="desc">6-dim breakdown — price strength · trend · short-term · money flow · volatility · base range — deterministic, no AI.</div></div></li>
        <li><span class="icon">📊</span><div class="ftext"><strong>Market analysis</strong><div class="desc">Pulse score · regime classifier (STABLE/VOLATILE/CRISIS/EUPHORIA) · sector flow · driver attribution · risk indicators.</div></div></li>
      </ul>
    </div>
    <div class="card">
      <ul class="feature-list">
        <li><span class="icon">🎯</span><div class="ftext"><strong>Signals</strong><div class="desc">Volume surge classifier · domestic flow (bid/ask pressure) · news clustering · earnings reaction.</div></div></li>
        <li><span class="icon">🔗</span><div class="ftext"><strong>Exposure chain</strong><div class="desc">Per-ticker external driver mapping with 2-year backtest-validated Pearson r — commodity / FX / peer correlations for 30+ VN30 tickers.</div></div></li>
        <li><span class="icon">🧪</span><div class="ftext"><strong>Backtest + cohort</strong><div class="desc">Strategy harness (buy-hold · momentum · low-vol · dual momentum) with realistic fees. Cohort analysis by RSI / MA trend / regime.</div></div></li>
        <li><span class="icon">🇻🇳</span><div class="ftext"><strong>Vietnamese NLP</strong><div class="desc">Sentiment (70+ financial keywords) · NLU intent parsing · ticker extraction · number parser.</div></div></li>
      </ul>
    </div>
  </div>
</section>
`)

	// === SECTION: Quick start ===
	sb.WriteString(`<section id="quickstart">
  <h2>🚀 Quick start</h2>
  <div class="grid grid-2">
    <div class="card">
      <div class="metric-label" style="margin-bottom:10px">Python</div>
<pre class="code"><span class="cmt"># pip install</span>
pip install lotusmarket==0.5.0

<span class="cmt"># Real-time quote + technical analysis</span>
<span class="kw">from</span> lotusmarket.fetchers <span class="kw">import</span> vps_quote
<span class="kw">from</span> lotusmarket <span class="kw">import</span> technical, ratings

quote = vps_quote(<span class="str">"ACB"</span>)
dash  = technical.dashboard(closes)
stars = ratings.compute(closes, vols, dash.score,
                        dash.rsi, dash.ma20, dash.ma50, dash.ma200)
<span class="kw">print</span>(stars.overall_verdict)  <span class="cmt"># Outperform</span></pre>
    </div>
    <div class="card">
      <div class="metric-label" style="margin-bottom:10px">Go</div>
<pre class="code"><span class="cmt">// go get</span>
go get github.com/ducnhd/lotusmarket/go@v0.5.0

<span class="cmt">// CLI binary</span>
go install github.com/ducnhd/lotusmarket/go/cmd/lmcli@v0.5.0

lmcli pulse              <span class="cmt"># daily VN30 markdown</span>
lmcli rate ACB           <span class="cmt"># 6-dim star ratings</span>
lmcli screen --rsi=30-50 <span class="cmt"># filter VN30</span>
lmcli dividends VNM      <span class="cmt"># corporate actions</span>
lmcli global             <span class="cmt"># 14 world indices</span></pre>
    </div>
  </div>
</section>
`)

	// === SECTION: How it works ===
	sb.WriteString(`<section id="how">
  <h2>⚙️ Cách trang này hoạt động</h2>
  <div class="card">
    <div class="metric-sub" style="font-size:15px; line-height:1.7">
      Trang này là <strong>100% bot, zero-cost, zero server</strong>:
      <br><br>
      • <strong>GitHub Actions</strong> chạy <code>lmcli</code> mỗi ngày 15:30 VN (sau khi thị trường VN đóng cửa) — free trên public repo.<br>
      • Output được commit vào <code>docs/</code>, GitHub Pages serve miễn phí.<br>
      • Bot Telegram <code>@vnlotusmarketbot</code> post tóm tắt lên <a href="https://t.me/vnlotusmarket">@vnlotusmarket</a>.<br>
      • Code generate trang này nằm trong <a href="https://github.com/ducnhd/lotusmarket/blob/main/go/cmd/lmcli/html.go">lmcli/html.go</a> — bạn có thể fork và customize.<br>
      <br>
      Tất cả dữ liệu lấy từ <strong>public API miễn phí</strong> (VPS, Entrade, Yahoo, VCI, Cafef). Không paid feed, không hot pick từ thầy, không vibe analysis — chỉ deterministic signal từ raw data.
    </div>
  </div>
</section>
`)

	// === SECTION: Archive ===
	sb.WriteString(`<section id="archive">
  <h2>📚 Archive</h2>
  <div class="card">
    <div class="metric-sub" style="margin-bottom:12px">Daily reports tự động cập nhật mỗi ngày làm việc.</div>
    <ul class="archive-list">
      <li>📄 <a href="latest.html">Latest report</a> — refresh để xem bản mới nhất</li>
      <li>📁 <a href="reports/">All reports</a> — archive theo ngày</li>
      <li>📝 <a href="blog/16-nam-backtest-vn-stocks.html">Blog: 16 năm backtest VN — 5 điều ngược với khuyên cũ</a></li>
    </ul>
  </div>
</section>
`)

	sb.WriteString(`</main>

<footer>
  <p>
    Built by <a href="https://github.com/ducnhd">@ducnhd</a> ·
    <a href="https://github.com/ducnhd/lotusmarket">GitHub</a> ·
    <a href="https://pypi.org/project/lotusmarket/">PyPI</a> ·
    <a href="https://t.me/vnlotusmarket">Telegram</a>
  </p>
  <p style="margin-top:8px">
    MIT License · Past performance does not guarantee future results · Not investment advice.
  </p>
</footer>

</body>
</html>
`)

	if outPath == "" {
		fmt.Print(sb.String())
		return
	}
	if err := writeFile(outPath, sb.String()); err != nil {
		panic(err)
	}
	fmt.Fprintf(stderrW(), "wrote %s (%d bytes)\n", outPath, sb.Len())
}
