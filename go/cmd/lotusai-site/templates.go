package main

import (
	"fmt"
	"html/template"
	"strings"
)

// parseTemplates builds one *Template per page by cloning the base + adding
// that page's body. Each page template defines a "body" block that the base
// renders via {{template "body" .}}. Without cloning, parsing multiple body
// blocks into the same template set would overwrite each other.
func parseTemplates() map[string]*template.Template {
	funcs := template.FuncMap{
		"pct":      pct,
		"pctSign":  pctSign,
		"pctClass": pctClass,
		"date":     dateFmt,
		"comma":    commaFmt,
		"divf":     func(a, b float64) float64 { return a / b },
	}
	base := template.New("base").Funcs(funcs)
	template.Must(base.Parse(baseTpl))
	template.Must(base.Parse(cssTpl))

	pages := map[string]string{
		"home":      homeTpl,
		"dashboard": dashboardTpl,
		"blog-list": blogListTpl,
		"blog-post": blogPostTpl,
		"docs":      docsTpl,
		"about":     aboutTpl,
		"not-found": notFoundTpl,
	}
	out := map[string]*template.Template{}
	for name, body := range pages {
		clone := template.Must(base.Clone())
		template.Must(clone.Parse(body))
		out[name] = clone
	}
	return out
}

func pct(v float64) string { return fmt.Sprintf("%+.2f%%", v) }
func pctSign(v float64) string {
	if v >= 0 {
		return "+"
	}
	return ""
}
func pctClass(v float64) string {
	if v > 0.05 {
		return "up"
	}
	if v < -0.05 {
		return "down"
	}
	return ""
}
func dateFmt(t interface{ Format(string) string }) string { return t.Format("02/01/2006") }
func commaFmt(n int64) string {
	s := fmt.Sprintf("%d", n)
	out := strings.Builder{}
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			out.WriteByte(',')
		}
		out.WriteRune(c)
	}
	return out.String()
}

// === Base layout ===
const baseTpl = `<!DOCTYPE html>
<html lang="vi">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{{.Title}}</title>
<meta name="description" content="{{.Description}}">
<link rel="canonical" href="{{.Canonical}}">
<meta name="robots" content="index,follow">
<meta name="author" content="ducnhd">
<!-- Open Graph -->
<meta property="og:type" content="{{.OGType}}">
<meta property="og:title" content="{{.Title}}">
<meta property="og:description" content="{{.Description}}">
<meta property="og:url" content="{{.Canonical}}">
<meta property="og:image" content="{{.OGImage}}">
<meta property="og:site_name" content="Lotus AI">
<meta property="og:locale" content="vi_VN">
<!-- Twitter -->
<meta name="twitter:card" content="summary_large_image">
<meta name="twitter:title" content="{{.Title}}">
<meta name="twitter:description" content="{{.Description}}">
<meta name="twitter:image" content="{{.OGImage}}">
<!-- Feeds -->
<link rel="alternate" type="application/rss+xml" title="Lotus AI Blog RSS" href="{{.BaseURL}}/feed.xml">
<link rel="icon" href="data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 100 100'%3E%3Ctext y='.9em' font-size='90'%3E%F0%9F%8C%B8%3C/text%3E%3C/svg%3E">
{{if .JSONLD}}<script type="application/ld+json">{{.JSONLD}}</script>{{end}}
<style>{{template "css" .}}</style>
</head>
<body>
<header class="topbar">
  <div class="container nav">
    <a class="brand" href="/"><span class="lotus">🌸</span> Lotus AI</a>
    <nav>
      <a href="/dashboard">Dashboard</a>
      <a href="/blog/">Blog</a>
      <a href="/docs">Docs</a>
      <a href="/about">About</a>
      <a class="cta" href="https://github.com/ducnhd/lotusmarket" rel="noopener" target="_blank">⭐ GitHub</a>
    </nav>
  </div>
</header>
<main class="container">{{template "body" .}}</main>
<footer>
  <div class="container">
    <p>
      <a href="/">Home</a> · <a href="/dashboard">Dashboard</a> · <a href="/blog/">Blog</a> ·
      <a href="/docs">Docs</a> · <a href="/feed.xml">RSS</a> ·
      <a href="https://github.com/ducnhd/lotusmarket" rel="noopener" target="_blank">GitHub</a> ·
      <a href="https://pypi.org/project/lotusmarket/" rel="noopener" target="_blank">PyPI</a> ·
      <a href="https://t.me/vnlotusmarket" rel="noopener" target="_blank">Telegram</a>
    </p>
    <p class="dim">© {{.Year}} lotusmarket · MIT License · Built with Go on a Raspberry Pi · Data is for educational purposes only, not investment advice.</p>
  </div>
</footer>
</body>
</html>`

const cssTpl = `{{define "css"}}
:root {
  --bg:#0f1419; --bg2:#1a2128; --bg-hover:#232c35;
  --text:#e6edf3; --dim:#8b949e;
  --accent:#7ee8a3; --accent-dim:#4d8c66;
  --red:#ff7b72; --green:#7ee8a3; --yellow:#ffd16a;
  --border:#2a323c; --radius:10px;
}
@media (prefers-color-scheme: light) {
  :root { --bg:#f6f8fa; --bg2:#fff; --bg-hover:#f0f3f6;
    --text:#1f2328; --dim:#59636e; --accent:#1a7f4e; --border:#d0d7de; }
}
* { box-sizing:border-box; margin:0; padding:0; }
html { scroll-behavior:smooth; }
body { background:var(--bg); color:var(--text);
  font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",system-ui,sans-serif;
  font-size:15px; line-height:1.55; -webkit-font-smoothing:antialiased; }
a { color:var(--accent); text-decoration:none; }
a:hover { text-decoration:underline; }
.container { max-width:1080px; margin:0 auto; padding:0 20px; }
header.topbar { border-bottom:1px solid var(--border); padding:14px 0;
  position:sticky; top:0; background:var(--bg); z-index:10;
  backdrop-filter:saturate(180%) blur(8px); }
.nav { display:flex; align-items:center; gap:16px; }
.brand { font-weight:700; font-size:18px; color:var(--text); }
.brand .lotus { margin-right:4px; }
.nav nav { margin-left:auto; display:flex; gap:14px; align-items:center; flex-wrap:wrap; }
.nav nav a { color:var(--dim); font-size:14px; }
.nav nav a:hover { color:var(--text); text-decoration:none; }
.nav nav a.cta { padding:6px 12px; border-radius:999px; background:var(--accent);
  color:#0f1419; font-weight:600; }
main { padding:24px 0 60px; }
section { margin:48px 0; }
h1 { font-size:34px; font-weight:800; letter-spacing:-0.02em; line-height:1.15; }
h2 { font-size:24px; font-weight:700; margin:1.4em 0 .6em; }
h3 { font-size:18px; font-weight:600; margin:1.2em 0 .5em; }
p { margin:0.6em 0; }
.hero { padding:48px 0 24px; text-align:center;
  background:linear-gradient(180deg,rgba(126,232,163,0.06),transparent); }
.hero h1 { font-size:46px; }
.hero p.tagline { color:var(--dim); font-size:18px; margin:8px 0 22px; }
.cta-row { display:flex; gap:10px; justify-content:center; flex-wrap:wrap; }
.cta-row a { padding:10px 18px; border-radius:999px; font-weight:600; font-size:14px;
  background:var(--bg2); border:1px solid var(--border); color:var(--text); }
.cta-row a.primary { background:var(--accent); color:#0f1419; border-color:var(--accent); }
.badges { margin-bottom:18px; display:inline-flex; gap:6px; flex-wrap:wrap; justify-content:center; }
.badge { font-size:11px; padding:3px 9px; border-radius:999px;
  background:var(--bg2); border:1px solid var(--border); color:var(--dim); }
.grid { display:grid; gap:12px; }
.grid-2 { grid-template-columns:repeat(auto-fit,minmax(320px,1fr)); }
.grid-3 { grid-template-columns:repeat(auto-fit,minmax(240px,1fr)); }
.grid-4 { grid-template-columns:repeat(auto-fit,minmax(180px,1fr)); }
.card { background:var(--bg2); border:1px solid var(--border); border-radius:var(--radius); padding:16px; }
.card.dense { padding:12px 14px; }
.label { font-size:12px; color:var(--dim); text-transform:uppercase; letter-spacing:0.05em; }
.value { font-size:24px; font-weight:700; margin-top:4px; }
.sub { font-size:13px; color:var(--dim); margin-top:2px; }
.pct { font-weight:600; }
.pct.up { color:var(--green); }
.pct.down { color:var(--red); }
.row { display:flex; justify-content:space-between; align-items:center;
  padding:6px 0; border-bottom:1px solid var(--border); }
.row:last-child { border-bottom:none; }
.ticker { font-family:"SF Mono",Menlo,Consolas,monospace; font-weight:600; }
.pulse-banner { background:linear-gradient(135deg,rgba(126,232,163,0.12),rgba(126,232,163,0.02));
  border:1px solid var(--accent-dim); border-radius:var(--radius); padding:22px 26px;
  display:flex; align-items:center; gap:20px; flex-wrap:wrap; }
.pulse-num { font-size:44px; font-weight:800; line-height:1; }
.pulse-text { flex:1; min-width:240px; }
.pulse-text strong { font-size:17px; }
ul.features { list-style:none; }
ul.features li { padding:12px 0; border-bottom:1px solid var(--border);
  display:flex; gap:12px; align-items:flex-start; }
ul.features li:last-child { border-bottom:none; }
ul.features .icon { font-size:20px; }
ul.features strong { display:block; margin-bottom:2px; }
ul.features .desc { color:var(--dim); font-size:13px; }
pre { background:var(--bg); border:1px solid var(--border); border-radius:6px;
  padding:14px; overflow-x:auto; font-family:"SF Mono",Menlo,monospace; font-size:13px; }
code { background:rgba(126,232,163,0.10); padding:1px 6px; border-radius:4px; font-family:"SF Mono",Menlo,monospace; font-size:0.9em; }
pre code { background:transparent; padding:0; }
table { border-collapse:collapse; width:100%; margin:14px 0; }
th,td { border:1px solid var(--border); padding:7px 10px; text-align:left; }
th { background:rgba(126,232,163,0.08); font-weight:600; }
blockquote { border-left:3px solid var(--accent); padding:8px 14px; margin:14px 0;
  background:rgba(126,232,163,0.05); color:var(--dim); font-style:italic; }
hr { border:none; border-top:1px solid var(--border); margin:24px 0; }
article ul, article ol { padding-left:22px; margin:0.6em 0; }
article li { margin:4px 0; }
.post-meta { color:var(--dim); font-size:13px; margin-bottom:16px; }
.post-list a.post { display:block; padding:16px; border:1px solid var(--border);
  border-radius:var(--radius); background:var(--bg2); margin-bottom:10px; color:var(--text); }
.post-list a.post:hover { background:var(--bg-hover); text-decoration:none; }
.post-list .post-title { font-size:18px; font-weight:600; }
.post-list .post-date { color:var(--dim); font-size:12px; }
.post-list .post-excerpt { color:var(--dim); font-size:14px; margin-top:6px; }
footer { padding:30px 0; border-top:1px solid var(--border); color:var(--dim); text-align:center; font-size:13px; }
footer .dim { font-size:12px; margin-top:8px; opacity:0.7; }
@media (max-width:600px) {
  h1 { font-size:28px; }
  .hero h1 { font-size:32px; }
  .container { padding:0 14px; }
  .pulse-num { font-size:34px; }
  .nav nav a { font-size:13px; }
}
{{end}}`

// === Page templates ===
const homeTpl = `{{define "body"}}
<section class="hero">
  <h1><span style="display:inline-block;margin-right:8px">🌸</span>Lotus AI</h1>
  <p class="tagline">Toolkit phân tích cổ phiếu Việt Nam — miễn phí, open source, không hot picks</p>
  <div class="badges">
    <span class="badge">VN30 + HNX30</span>
    <span class="badge">Go + Python</span>
    <span class="badge">MIT License</span>
    <span class="badge">Cập nhật mỗi 10 phút</span>
  </div>
  <div class="cta-row">
    <a class="primary" href="/dashboard">📊 Live Dashboard</a>
    <a href="/blog/">📝 Blog</a>
    <a href="https://github.com/ducnhd/lotusmarket" rel="noopener" target="_blank">⭐ GitHub</a>
    <a href="https://t.me/vnlotusmarket" rel="noopener" target="_blank">💬 Telegram</a>
  </div>
</section>

{{if and .Snap .Snap.Stocks}}
<section>
  <h2>📊 Nhịp thị trường hôm nay</h2>
  <div class="pulse-banner">
    <div>
      <div class="pulse-num">{{printf "%.0f%%" .Snap.Flow.BuyPressure}}</div>
      <div class="label">áp lực mua</div>
    </div>
    <div class="pulse-text">
      <strong>{{.Snap.Flow.Signal}}</strong>
      <div class="sub">Tính từ tỷ lệ bid/ask volume trên 30 mã VN30. Cập nhật: {{.Snap.Updated.Format "15:04, 02/01/2006"}}</div>
    </div>
  </div>
</section>

<section>
  <h2>🏭 Dòng tiền sector (top 6)</h2>
  <div class="grid grid-3">
    {{range $i, $s := .Snap.Sectors}}{{if lt $i 6}}
    <div class="card">
      <div class="label">#{{$s.Rank}} {{$s.Sector}}</div>
      <div class="value pct {{pctClass $s.AvgChangePct}}">{{pct $s.AvgChangePct}}</div>
      <div class="sub">KL {{printf "%.0f" (divf $s.TotalVolumeVND 1e9)}} tỷ · {{$s.AdvancesCount}}/{{$s.TotalCount}} tăng</div>
    </div>
    {{end}}{{end}}
  </div>
</section>
{{end}}

<section>
  <h2>🧰 Vì sao Lotus AI?</h2>
  <div class="grid grid-2">
    <div class="card">
      <ul class="features">
        <li><span class="icon">🔬</span><div><strong>Deterministic</strong><div class="desc">Mọi signal tính từ raw price + volume. Không AI guess, không "vibe analysis", không hot pick.</div></div></li>
        <li><span class="icon">🆓</span><div><strong>Miễn phí toàn bộ</strong><div class="desc">Public API (VPS · Entrade · KBS · Yahoo · VCI · Cafef). Không paid feed. MIT license.</div></div></li>
        <li><span class="icon">🇻🇳</span><div><strong>VN-native</strong><div class="desc">Sentiment tiếng Việt, T+2 settlement, sector VN, VN30/HNX30 ticker mapping.</div></div></li>
        <li><span class="icon">📦</span><div><strong>2 ngôn ngữ</strong><div class="desc">Go (single binary) + Python (pip install). Test parity across both.</div></div></li>
      </ul>
    </div>
    <div class="card">
      <ul class="features">
        <li><span class="icon">📊</span><div><strong>Cohort analysis</strong><div class="desc">47,000 data point, 9 năm. Edge thật cho từng pattern (RSI × MA trend × regime).</div></div></li>
        <li><span class="icon">🌡️</span><div><strong>Regime classifier</strong><div class="desc">STABLE / VOLATILE / CRISIS / EUPHORIA — deterministic, không AI. CRISIS_PANIC vs CRISIS_FUNDAMENTAL sub-types.</div></div></li>
        <li><span class="icon">⭐</span><div><strong>Star ratings</strong><div class="desc">6-dim breakdown — price strength · trend · RSI · money flow · volatility · base range.</div></div></li>
        <li><span class="icon">🤖</span><div><strong>Auto daily report</strong><div class="desc">Cập nhật mỗi 15:30 VN qua GitHub Actions + Telegram. Zero server cost.</div></div></li>
      </ul>
    </div>
  </div>
</section>

<section>
  <h2>🚀 Quick start</h2>
  <div class="grid grid-2">
    <div class="card">
      <div class="label" style="margin-bottom:10px">Python</div>
      <pre><code>pip install lotusmarket==0.5.0

from lotusmarket.fetchers import vps_quote
from lotusmarket import technical, ratings

quote = vps_quote("ACB")
dash  = technical.dashboard(closes)
stars = ratings.compute(closes, vols,
                        dash.score, dash.rsi,
                        dash.ma20, dash.ma50, dash.ma200)
print(stars.overall_verdict)  # Outperform</code></pre>
    </div>
    <div class="card">
      <div class="label" style="margin-bottom:10px">Go</div>
      <pre><code>go get github.com/ducnhd/lotusmarket/go@v0.5.0

# CLI
go install github.com/ducnhd/lotusmarket/go/cmd/lmcli@v0.5.0

lmcli pulse              # daily VN30
lmcli rate ACB           # 6-dim ratings
lmcli screen --rsi=30-50 # filter VN30
lmcli dividends VNM      # corp actions</code></pre>
    </div>
  </div>
</section>

<section>
  <h2>📝 Bài viết gần đây</h2>
  <p class="sub">Phân tích bằng data thật. Backtest reproducible. Không opinion piece.</p>
  <p><a href="/blog/" class="cta">Xem toàn bộ blog →</a></p>
</section>
{{end}}`

const dashboardTpl = `{{define "body"}}
<section>
  <h1>📊 Dashboard thị trường Việt Nam</h1>
  <p class="sub">Cập nhật: <strong>{{.Snap.Updated.Format "15:04, 02/01/2006"}}</strong> — refresh page để load data mới (cache 10 phút).</p>
</section>

{{if .Snap.Stocks}}
<section>
  <div class="pulse-banner">
    <div>
      <div class="pulse-num">{{printf "%.0f%%" .Snap.Flow.BuyPressure}}</div>
      <div class="label">áp lực mua</div>
    </div>
    <div class="pulse-text">
      <strong>Dòng tiền nội địa: {{.Snap.Flow.Signal}}</strong>
      <div class="sub">Tính từ tỷ lệ bid/ask volume trên VN30. Deterministic, không có AI bias.</div>
    </div>
  </div>
</section>

<section>
  <h2>🏭 Sector flow ranking</h2>
  <div class="grid grid-3">
    {{range .Snap.Sectors}}
    <div class="card">
      <div class="label">#{{.Rank}} {{.Sector}}</div>
      <div class="value pct {{pctClass .AvgChangePct}}">{{pct .AvgChangePct}}</div>
      <div class="sub">KL {{printf "%.0f" (divf .TotalVolumeVND 1e9)}} tỷ · {{.AdvancesCount}}/{{.TotalCount}} tăng</div>
    </div>
    {{end}}
  </div>
</section>

<section>
  <h2>📈 Top tăng / 📉 Top giảm VN30</h2>
  <div class="grid grid-2">
    <div class="card">
      <div class="label" style="margin-bottom:10px">📈 Top tăng</div>
      {{range .Snap.TopUp}}
      <div class="row">
        <span class="ticker">{{.Ticker}}</span>
        <span class="sub">{{printf "%.0f" .Close}} · KL {{comma .Volume}}</span>
        <span class="pct up">{{pct .ChangePercent}}</span>
      </div>
      {{end}}
    </div>
    <div class="card">
      <div class="label" style="margin-bottom:10px">📉 Top giảm</div>
      {{range .Snap.TopDown}}
      <div class="row">
        <span class="ticker">{{.Ticker}}</span>
        <span class="sub">{{printf "%.0f" .Close}} · KL {{comma .Volume}}</span>
        <span class="pct down">{{pct .ChangePercent}}</span>
      </div>
      {{end}}
    </div>
  </div>
</section>
{{end}}

{{if .Snap.GlobalQuote}}
<section>
  <h2>🌏 Thị trường thế giới</h2>
  <div class="grid grid-4">
    {{range .Snap.GlobalQuote}}
    <div class="card dense">
      <div class="label">{{.Symbol}}</div>
      <div class="value">{{printf "%.2f" .Price}}</div>
      <div class="sub pct {{pctClass .ChangePct}}">{{pct .ChangePct}}</div>
    </div>
    {{end}}
  </div>
</section>
{{end}}

<section>
  <p class="sub">Dữ liệu được fetch realtime từ VPS HOSE feed (VN tickers) và Yahoo Finance (global indices), cache server-side 10 phút. Source code: <a href="https://github.com/ducnhd/lotusmarket">github.com/ducnhd/lotusmarket</a></p>
</section>
{{end}}`

const blogListTpl = `{{define "body"}}
<section>
  <h1>📝 Blog</h1>
  <p class="sub">Phân tích cổ phiếu Việt Nam bằng data thật. Backtest reproducible. Code MIT open source.</p>
</section>
<section class="post-list">
  {{range .Posts}}
  <a class="post" href="/blog/{{.Slug}}">
    <div class="post-title">{{.Title}}</div>
    <div class="post-date">{{.Date.Format "02/01/2006"}}</div>
    <div class="post-excerpt">{{.Excerpt}}</div>
  </a>
  {{else}}
  <p>Chưa có bài viết.</p>
  {{end}}
</section>
{{end}}`

const blogPostTpl = `{{define "body"}}
<article>
  <p class="post-meta"><a href="/blog/">← Tất cả bài viết</a> · {{.Post.Date.Format "02/01/2006"}}</p>
  {{.Post.HTML}}
  <hr>
  <p class="sub">
    Thích bài này? ⭐ <a href="https://github.com/ducnhd/lotusmarket">Star repo</a> ·
    💬 <a href="https://t.me/vnlotusmarket">Subscribe Telegram</a> để nhận daily report ·
    📡 <a href="/feed.xml">RSS feed</a>
  </p>
</article>
{{end}}`

const docsTpl = `{{define "body"}}
<article>
  <h1>📚 Docs &amp; Quick start</h1>
  {{.Body}}
</article>
{{end}}`

const aboutTpl = `{{define "body"}}
<article>
  <h1>About</h1>
  {{.Body}}
</article>
{{end}}`

const notFoundTpl = `{{define "body"}}
<section style="text-align:center; padding:80px 0">
  <h1>404</h1>
  <p class="sub">Trang không tồn tại.</p>
  <p style="margin-top:20px"><a href="/">← Về trang chủ</a></p>
</section>
{{end}}`

// === Default markdown content (used when content/docs.md missing) ===

const defaultDocsMarkdown = `## Install

### Python

` + "```bash\npip install lotusmarket==0.5.0\n```" + `

Optional extras:

` + "```bash\npip install lotusmarket[fetchers]   # adds httpx for API fetchers\npip install lotusmarket[ai]         # adds anthropic for AI helpers\npip install lotusmarket[all]\n```" + `

### Go

` + "```bash\ngo get github.com/ducnhd/lotusmarket/go@v0.5.0\n```" + `

CLI binary:

` + "```bash\ngo install github.com/ducnhd/lotusmarket/go/cmd/lmcli@v0.5.0\n\nlmcli pulse           # daily VN30 markdown\nlmcli quote ACB       # real-time quote\nlmcli rate ACB        # 6-dim star ratings\nlmcli screen --rsi=30-50\nlmcli sectors\nlmcli global\nlmcli dividends VNM\nlmcli report --out=today.md\n```" + `

## Modules

- **fetchers** — VPS, Entrade, KBS, Cafef, Yahoo Finance, VCI dividends, FRED macro (all free, no API key)
- **technical** — RSI, MA, MACD, Bollinger, ATR, momentum, dashboard signal + score
- **ratings** — 6-dim star ratings (price strength · trend · RSI · money flow · volatility · base range)
- **signals** — volume surge classifier, domestic flow (bid/ask pressure)
- **market** — pulse score, regime classifier, sector flow, driver attribution
- **portfolio** — TWR returns, confidence scoring
- **sentiment** — Vietnamese keyword analyzer (70+ financial terms)
- **nlu** — intent parsing, ticker extraction
- **exposure** — per-ticker external driver mapping (commodity / FX / peer) with 2-year backtest correlations
- **backtest** — strategy harness (buyhold · momentum · low-vol · dual momentum) with realistic VN retail fees
- **historical** — generic cohort analysis on (ticker, date, features, fwd_returns) panels
- **earnings** — Vietnamese earnings headline parser
- **ai** — optional Claude API helper for narrative summaries

## License

MIT — fork, modify, commercialize, no obligation. PRs welcome.

## Support

If lotusmarket saved you time, consider supporting:

- [Star the GitHub repo](https://github.com/ducnhd/lotusmarket)
- [Subscribe Telegram channel](https://t.me/vnlotusmarket)
- [Sponsor via PayPal](https://www.paypal.com/paypalme/ducnhd)
- Send Telegram Stars to [@vnlotusmarketbot](https://t.me/vnlotusmarketbot) ` + "`/donate`"

const defaultAboutMarkdown = `## Lotus AI là gì?

Lotus AI là **trang quảng bá + dashboard** của project [lotusmarket](https://github.com/ducnhd/lotusmarket) — một thư viện open source MIT để phân tích cổ phiếu Việt Nam bằng Go và Python.

Lotusmarket tồn tại vì 1 lý do: **retail VN deserve công cụ tốt hơn "hot pick từ thầy".**

7 triệu+ tài khoản retail tại VN hầu hết dùng:

- App broker với UI rối, không có AI analysis
- Telegram group "phím hàng" trả phí 1-3 triệu/tháng, chất lượng phập phù
- "Khóa học đầu tư" 5-50 triệu, đa số content tái chế

Không có lựa chọn ở giữa: **một bộ tool data-driven, deterministic, miễn phí, dành cho người tự ra quyết định.**

## Ai làm cái này?

Mình là **ducnhd** — backend developer 8 năm kinh nghiệm, đầu tư cá nhân chứng khoán VN từ 2017. Năm 2022 mình mất gần một nửa danh mục do bán đáy panic (Vạn Thịnh Phát + BĐS sập). Sau đó mình ngồi build 1 bot Telegram cá nhân để tự đầu tư cho mình.

Bot đó (vnstock-bot) đang chạy production 24/7 trên Raspberry Pi của mình từ 2024. Tuần này mình tách phần engine ra thành thư viện open source: **lotusmarket**.

## Tech stack

- **Go 1.25** — core library, CLI, daily report generator, web site
- **Python 3.9+** — parity package trên PyPI
- **Raspberry Pi 4** — production server (web, bot, cron)
- **GitHub Actions** — daily/weekly automation, releases
- **MIT License** — fork freely

## Liên hệ

- GitHub: [@ducnhd](https://github.com/ducnhd)
- Telegram: [@vnlotusmarket](https://t.me/vnlotusmarket) (channel) · [@vnlotusmarketbot](https://t.me/vnlotusmarketbot) (bot)

## Disclaimer

Tất cả data + analysis trên trang này là **educational**. Không phải lời khuyên đầu tư. Past performance does not guarantee future results. Bạn tự chịu trách nhiệm với quyết định đầu tư của mình.
`
