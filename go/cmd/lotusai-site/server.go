package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ducnhd/lotusmarket/go/fetchers"
	"github.com/ducnhd/lotusmarket/go/market"
	"github.com/ducnhd/lotusmarket/go/signals"
	"github.com/ducnhd/lotusmarket/go/types"
)

// ===== config & state =====

type serverConfig struct {
	BaseURL    string
	ContentDir string
	CacheTTL   time.Duration
	DataDir    string
}

type server struct {
	cfg       serverConfig
	tpl       map[string]*template.Template
	cache     liveCache
	startedAt time.Time
	views     *viewStore
}

func newServer(cfg serverConfig) *server {
	s := &server{
		cfg:       cfg,
		cache:     liveCache{ttl: cfg.CacheTTL},
		startedAt: time.Now(),
	}
	s.tpl = parseTemplates()
	if err := os.MkdirAll(cfg.DataDir, 0o755); err == nil {
		s.views = newViewStore(filepath.Join(cfg.DataDir, "blog_views.json"))
		go func() {
			t := time.NewTicker(2 * time.Minute)
			for range t.C {
				_ = s.views.flush()
			}
		}()
	}
	return s
}

// ===== routes =====

func (s *server) routes(mux *http.ServeMux) {
	mux.HandleFunc("/", s.handleRoot)
	mux.HandleFunc("/dashboard", s.handleDashboard)
	mux.HandleFunc("/dashboard/", s.handleDashboard) // accept trailing slash too — avoids 404
	// /blog/stats must be registered before /blog/ so ServeMux routes it to
	// its own handler rather than the catch-all handleBlog.
	mux.HandleFunc("/blog/stats", s.handleBlogStats)
	// Register both /blog and /blog/ as the same handler so Google never sees
	// a 301 redirect (GSC was flagging "Lỗi chuyển hướng" / Page with redirect
	// because Go's http.ServeMux auto-redirects "/blog" → "/blog/").
	mux.HandleFunc("/blog", s.handleBlog)
	mux.HandleFunc("/blog/", s.handleBlog)
	mux.HandleFunc("/docs", s.handleDocs)
	mux.HandleFunc("/docs/", s.handleDocs)
	mux.HandleFunc("/about", s.handleAbout)
	mux.HandleFunc("/about/", s.handleAbout)
	mux.HandleFunc("/sitemap.xml", s.handleSitemap)
	mux.HandleFunc("/robots.txt", s.handleRobots)
	mux.HandleFunc("/feed.xml", s.handleFeed)
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/og/", s.handleOGImage)
	// Static assets (CSS/img) served from content/static/
	staticDir := filepath.Join(s.cfg.ContentDir, "static")
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))
}

// ===== middleware =====

func withMiddleware(h http.Handler) http.Handler {
	return loggingMW(securityHeadersMW(h))
}

func loggingMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t0 := time.Now()
		sw := &statusWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(sw, r)
		log.Printf("%s %s %d %s %s", r.Method, r.URL.Path, sw.status, time.Since(t0), r.UserAgent())
	})
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (s *statusWriter) WriteHeader(c int) { s.status = c; s.ResponseWriter.WriteHeader(c) }

func securityHeadersMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		next.ServeHTTP(w, r)
	})
}

// ===== live data cache =====

type liveSnapshot struct {
	Updated     time.Time
	Flow        signals.DomesticFlow
	Sectors     []market.SectorFlow
	TopUp       []types.StockData
	TopDown     []types.StockData
	Stocks      []types.StockData
	GlobalQuote []fetchers.YahooQuote
}

type liveCache struct {
	mu   sync.Mutex
	ttl  time.Duration
	last time.Time
	snap *liveSnapshot
}

var vn30 = []string{
	"ACB", "BCM", "BID", "BVH", "CTG", "FPT", "GAS", "GVR", "HDB", "HPG",
	"MBB", "MSN", "MWG", "PLX", "POW", "SAB", "SHB", "SSB", "SSI", "STB",
	"TCB", "TPB", "VCB", "VHM", "VIB", "VIC", "VJC", "VNM", "VPB", "VRE",
}

var globalSymbols = []string{"^GSPC", "^DJI", "^IXIC", "^N225", "^HSI", "^VIX", "GC=F", "CL=F", "DX-Y.NYB"}

func (c *liveCache) get(ctx context.Context) *liveSnapshot {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.snap != nil && time.Since(c.last) < c.ttl {
		return c.snap
	}
	snap := fetchLive(ctx)
	c.snap = snap
	c.last = time.Now()
	return snap
}

func fetchLive(ctx context.Context) *liveSnapshot {
	snap := &liveSnapshot{Updated: time.Now()}
	stocks, _ := fetchers.VPSMultiple(ctx, vn30)
	if len(stocks) == 0 {
		log.Println("[lotusai-site] WARN: empty stock snapshot")
		return snap
	}
	snap.Stocks = stocks
	snap.Flow = signals.ComputeDomesticFlow(stocks)
	snap.Sectors = market.RankSectorsByFlow(stocks)

	sorted := make([]types.StockData, len(stocks))
	copy(sorted, stocks)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].ChangePercent > sorted[j].ChangePercent })
	n := len(sorted)
	if n >= 5 {
		snap.TopUp = sorted[:5]
		snap.TopDown = make([]types.StockData, 5)
		for i := 0; i < 5; i++ {
			snap.TopDown[i] = sorted[n-1-i]
		}
	}
	snap.GlobalQuote = fetchers.YahooMultiple(ctx, globalSymbols)
	return snap
}

// ===== blog =====

type blogPost struct {
	Slug    string
	Title   string
	Date    time.Time
	Body    string
	HTML    template.HTML
	Excerpt string
}

func (s *server) loadBlogPosts() []blogPost {
	dir := filepath.Join(s.cfg.ContentDir, "blog")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	posts := []blogPost{}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		body, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		slug := strings.TrimSuffix(e.Name(), ".md")
		title, date := parseFrontmatter(string(body))
		if title == "" {
			title = slug
		}
		clean := stripFrontmatter(string(body))
		posts = append(posts, blogPost{
			Slug:    slug,
			Title:   title,
			Date:    date,
			Body:    clean,
			HTML:    renderMarkdown(clean),
			Excerpt: firstParagraph(clean, 300),
		})
	}
	sort.Slice(posts, func(i, j int) bool { return posts[i].Date.After(posts[j].Date) })
	return posts
}

func parseFrontmatter(body string) (title string, date time.Time) {
	date = time.Now()
	if !strings.HasPrefix(body, "---") {
		return "", date
	}
	end := strings.Index(body[3:], "---")
	if end < 0 {
		return "", date
	}
	for _, line := range strings.Split(body[3:end+3], "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "title:") {
			title = strings.Trim(strings.TrimSpace(strings.TrimPrefix(line, "title:")), `"`)
		}
		if strings.HasPrefix(line, "date:") {
			d := strings.TrimSpace(strings.TrimPrefix(line, "date:"))
			if t, err := time.Parse("2006-01-02", d); err == nil {
				date = t
			}
		}
	}
	return
}

func stripFrontmatter(body string) string {
	if !strings.HasPrefix(body, "---") {
		return body
	}
	end := strings.Index(body[3:], "---")
	if end < 0 {
		return body
	}
	return strings.TrimSpace(body[end+6:])
}

func firstParagraph(body string, max int) string {
	// Skip leading # headers and blockquotes
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ">") || strings.HasPrefix(line, "---") {
			continue
		}
		if len(line) > max {
			line = line[:max] + "…"
		}
		// Strip basic markdown
		line = strings.ReplaceAll(line, "**", "")
		line = strings.ReplaceAll(line, "*", "")
		return line
	}
	return ""
}

// renderMarkdown — minimal subset suitable for our content (headings, lists,
// bold/italic, code, tables, links, paragraphs). No external dep — keeps the
// binary tiny.
func renderMarkdown(md string) template.HTML {
	var out strings.Builder
	lines := strings.Split(md, "\n")
	inCode := false
	inTable := false
	inList := false
	flushList := func() {
		if inList {
			out.WriteString("</ul>\n")
			inList = false
		}
	}
	flushTable := func() {
		if inTable {
			out.WriteString("</tbody></table>\n")
			inTable = false
		}
	}
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		trim := strings.TrimSpace(line)

		// Code fences
		if strings.HasPrefix(trim, "```") {
			flushList()
			flushTable()
			if inCode {
				out.WriteString("</code></pre>\n")
				inCode = false
			} else {
				out.WriteString("<pre><code>")
				inCode = true
			}
			continue
		}
		if inCode {
			out.WriteString(template.HTMLEscapeString(line))
			out.WriteByte('\n')
			continue
		}

		if trim == "" {
			flushList()
			flushTable()
			continue
		}

		// Headings
		if strings.HasPrefix(trim, "### ") {
			flushList()
			flushTable()
			fmt.Fprintf(&out, "<h3>%s</h3>\n", inline(strings.TrimPrefix(trim, "### ")))
			continue
		}
		if strings.HasPrefix(trim, "## ") {
			flushList()
			flushTable()
			fmt.Fprintf(&out, "<h2>%s</h2>\n", inline(strings.TrimPrefix(trim, "## ")))
			continue
		}
		if strings.HasPrefix(trim, "# ") {
			flushList()
			flushTable()
			fmt.Fprintf(&out, "<h1>%s</h1>\n", inline(strings.TrimPrefix(trim, "# ")))
			continue
		}

		// Blockquotes
		if strings.HasPrefix(trim, "> ") {
			flushList()
			flushTable()
			fmt.Fprintf(&out, "<blockquote>%s</blockquote>\n", inline(strings.TrimPrefix(trim, "> ")))
			continue
		}

		// HR
		if trim == "---" || trim == "***" {
			flushList()
			flushTable()
			out.WriteString("<hr>\n")
			continue
		}

		// Tables — detect `| col | col |` followed by `|---|---|`
		if strings.HasPrefix(trim, "|") && strings.HasSuffix(trim, "|") {
			if !inTable {
				// Look ahead for separator
				if i+1 < len(lines) && strings.Contains(lines[i+1], "---") {
					out.WriteString("<table><thead><tr>")
					for _, c := range splitTableRow(trim) {
						fmt.Fprintf(&out, "<th>%s</th>", inline(c))
					}
					out.WriteString("</tr></thead><tbody>\n")
					inTable = true
					i++ // skip separator
					continue
				}
			} else {
				out.WriteString("<tr>")
				for _, c := range splitTableRow(trim) {
					fmt.Fprintf(&out, "<td>%s</td>", inline(c))
				}
				out.WriteString("</tr>\n")
				continue
			}
		}
		flushTable()

		// Lists
		if strings.HasPrefix(trim, "- ") || strings.HasPrefix(trim, "* ") {
			if !inList {
				out.WriteString("<ul>\n")
				inList = true
			}
			item := strings.TrimSpace(trim[2:])
			fmt.Fprintf(&out, "<li>%s</li>\n", inline(item))
			continue
		}
		flushList()

		// Paragraph
		fmt.Fprintf(&out, "<p>%s</p>\n", inline(trim))
	}
	flushList()
	flushTable()
	if inCode {
		out.WriteString("</code></pre>\n")
	}
	return template.HTML(out.String())
}

func splitTableRow(s string) []string {
	s = strings.TrimPrefix(s, "|")
	s = strings.TrimSuffix(s, "|")
	parts := strings.Split(s, "|")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

// inline handles **bold**, *italic*, `code`, [text](url) — order matters.
func inline(s string) string {
	// Escape HTML first, then re-introduce our tags via marker tokens to avoid double-escape.
	s = template.HTMLEscapeString(s)
	// Code spans
	s = replaceWrap(s, "`", "<code>", "</code>")
	// Bold
	s = replaceWrap(s, "**", "<strong>", "</strong>")
	// Italic
	s = replaceWrap(s, "*", "<em>", "</em>")
	// Links [text](url) — pattern simple, no nested
	for {
		o := strings.Index(s, "[")
		if o < 0 {
			break
		}
		c := strings.Index(s[o:], "](")
		if c < 0 {
			break
		}
		e := strings.Index(s[o+c:], ")")
		if e < 0 {
			break
		}
		text := s[o+1 : o+c]
		url := s[o+c+2 : o+c+e]
		s = s[:o] + `<a href="` + url + `">` + text + `</a>` + s[o+c+e+1:]
	}
	return s
}

func replaceWrap(s, marker, openTag, closeTag string) string {
	out := strings.Builder{}
	open := true
	parts := strings.Split(s, marker)
	for i, p := range parts {
		out.WriteString(p)
		if i == len(parts)-1 {
			break
		}
		if open {
			out.WriteString(openTag)
		} else {
			out.WriteString(closeTag)
		}
		open = !open
	}
	// If we ended in an unclosed wrapper, undo
	if !open {
		// Unbalanced; fall back to raw
		return s
	}
	return out.String()
}

// ===== handlers =====

func (s *server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "ok uptime=%s\n", time.Since(s.startedAt).Round(time.Second))
}

func (s *server) handleRobots(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "User-agent: *\nAllow: /\nSitemap: %s/sitemap.xml\n", s.cfg.BaseURL)
}

func (s *server) handleSitemap(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	now := time.Now().Format("2006-01-02")
	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	sb.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n")
	pages := []struct {
		Loc      string
		Priority string
	}{
		{"/", "1.0"}, {"/dashboard", "0.9"}, {"/blog/", "0.8"},
		{"/docs", "0.7"}, {"/about", "0.5"},
	}
	for _, p := range pages {
		fmt.Fprintf(&sb, "  <url><loc>%s%s</loc><lastmod>%s</lastmod><priority>%s</priority></url>\n",
			s.cfg.BaseURL, p.Loc, now, p.Priority)
	}
	for _, p := range s.loadBlogPosts() {
		fmt.Fprintf(&sb, "  <url><loc>%s/blog/%s</loc><lastmod>%s</lastmod><priority>0.6</priority></url>\n",
			s.cfg.BaseURL, p.Slug, p.Date.Format("2006-01-02"))
	}
	sb.WriteString("</urlset>\n")
	_, _ = w.Write([]byte(sb.String()))
}

func (s *server) handleFeed(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
	posts := s.loadBlogPosts()
	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	sb.WriteString(`<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom"><channel>` + "\n")
	fmt.Fprintf(&sb, "<title>Lotus AI · VN Stock Market Toolkit</title>\n")
	fmt.Fprintf(&sb, "<link>%s/</link>\n", s.cfg.BaseURL)
	fmt.Fprintf(&sb, `<atom:link href="%s/feed.xml" rel="self" type="application/rss+xml" />`+"\n", s.cfg.BaseURL)
	fmt.Fprintf(&sb, "<description>Vietnamese stock market analysis, backtest, and open-source toolkit — free, deterministic, no hot picks.</description>\n")
	fmt.Fprintf(&sb, "<language>vi</language>\n")
	fmt.Fprintf(&sb, "<lastBuildDate>%s</lastBuildDate>\n", time.Now().UTC().Format(time.RFC1123Z))
	for _, p := range posts {
		fmt.Fprintf(&sb, "<item>\n")
		fmt.Fprintf(&sb, "  <title>%s</title>\n", template.HTMLEscapeString(p.Title))
		fmt.Fprintf(&sb, "  <link>%s/blog/%s</link>\n", s.cfg.BaseURL, p.Slug)
		fmt.Fprintf(&sb, "  <guid isPermaLink=\"true\">%s/blog/%s</guid>\n", s.cfg.BaseURL, p.Slug)
		fmt.Fprintf(&sb, "  <pubDate>%s</pubDate>\n", p.Date.UTC().Format(time.RFC1123Z))
		fmt.Fprintf(&sb, "  <description>%s</description>\n", template.HTMLEscapeString(p.Excerpt))
		fmt.Fprintf(&sb, "</item>\n")
	}
	sb.WriteString("</channel></rss>\n")
	_, _ = w.Write([]byte(sb.String()))
}

func (s *server) handleBlogStats(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if s.views == nil {
		_, _ = w.Write([]byte("{}"))
		return
	}
	b, _ := json.MarshalIndent(s.views.snapshot(), "", "  ")
	_, _ = w.Write(b)
}
