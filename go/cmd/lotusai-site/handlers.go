package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// pageData carries SEO meta + page-specific payload.
type pageData struct {
	Title       string
	Description string
	Canonical   string
	OGType      string
	OGImage     string
	PageType    string // "home" | "dashboard" | "blog-list" | "blog-post" | "docs" | "about"
	BaseURL     string
	Year        int
	JSONLD      template.JS

	// Live snapshot (home + dashboard)
	Snap *liveSnapshot

	// Blog
	Posts []blogPost
	Post  *blogPost

	// Docs / About
	Body template.HTML
}

func (s *server) renderPage(w http.ResponseWriter, name string, data pageData) {
	if data.BaseURL == "" {
		data.BaseURL = s.cfg.BaseURL
	}
	if data.OGType == "" {
		data.OGType = "website"
	}
	if data.OGImage == "" {
		data.OGImage = data.BaseURL + "/og/default.png"
	}
	data.Year = time.Now().Year()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=600") // 10 min — matches snap TTL
	tpl, ok := s.tpl[name]
	if !ok {
		http.Error(w, "template not found: "+name, http.StatusInternalServerError)
		return
	}
	if err := tpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// --- Home (/) ---

func (s *server) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		s.handleNotFound(w, r)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 25*time.Second)
	defer cancel()
	snap := s.cache.get(ctx)

	s.renderPage(w, "home", pageData{
		Title:       "Lotus AI — VN Stock Market Toolkit & Daily Dashboard",
		Description: "Free, open-source Vietnamese stock market analytics. Live VN30 + HNX30 dashboard, backtest harness, deterministic signals — no hot picks, no paid feed. Built by retail investors for retail investors.",
		Canonical:   s.cfg.BaseURL + "/",
		PageType:    "home",
		Snap:        snap,
		JSONLD:      jsonLDOrganization(s.cfg.BaseURL),
	})
}

// --- Dashboard (/dashboard) ---

func (s *server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 25*time.Second)
	defer cancel()
	snap := s.cache.get(ctx)

	updated := snap.Updated.Format("15:04, 02/01/2006")
	s.renderPage(w, "dashboard", pageData{
		Title:       fmt.Sprintf("Bảng điều khiển thị trường VN — Lotus AI (cập nhật %s)", updated),
		Description: "Bảng dữ liệu thị trường chứng khoán Việt Nam cập nhật mỗi 10 phút: VN30, HNX30, dòng tiền sector, top tăng/giảm, chỉ số toàn cầu. Miễn phí, open source.",
		Canonical:   s.cfg.BaseURL + "/dashboard",
		PageType:    "dashboard",
		Snap:        snap,
	})
}

// --- Blog list & detail (/blog/, /blog/{slug}) ---

func (s *server) handleBlog(w http.ResponseWriter, r *http.Request) {
	posts := s.loadBlogPosts()
	// Treat /blog and /blog/ identically — both render the list.
	// Strip both forms so /blog/<slug> still resolves correctly.
	path := r.URL.Path
	var rest string
	switch {
	case path == "/blog", path == "/blog/":
		rest = ""
	default:
		rest = strings.TrimPrefix(path, "/blog/")
		rest = strings.TrimSuffix(rest, "/")
	}

	if rest == "" {
		s.renderPage(w, "blog-list", pageData{
			Title:       "Blog — Lotus AI · Phân tích cổ phiếu Việt Nam bằng data thật",
			Description: "Bài viết về backtest, cohort analysis, regime classification, và lý thuyết đầu tư cổ phiếu Việt Nam. Code MIT open source, reproducible.",
			Canonical:   s.cfg.BaseURL + "/blog/",
			PageType:    "blog-list",
			Posts:       posts,
		})
		return
	}

	var match *blogPost
	for i := range posts {
		if posts[i].Slug == rest {
			match = &posts[i]
			break
		}
	}
	if match == nil {
		s.handleNotFound(w, r)
		return
	}
	if s.views != nil {
		s.views.hit(match.Slug, r.UserAgent())
	}
	s.renderPage(w, "blog-post", pageData{
		Title:       match.Title + " · Lotus AI",
		Description: match.Excerpt,
		Canonical:   s.cfg.BaseURL + "/blog/" + match.Slug,
		OGType:      "article",
		OGImage:     s.cfg.BaseURL + "/og/blog/" + match.Slug + ".png",
		PageType:    "blog-post",
		Post:        match,
		JSONLD:      jsonLDArticle(s.cfg.BaseURL, match),
	})
}

// --- Docs (/docs) ---

func (s *server) handleDocs(w http.ResponseWriter, r *http.Request) {
	docsPath := filepath.Join(s.cfg.ContentDir, "docs.md")
	body, err := os.ReadFile(docsPath)
	if err != nil {
		body = []byte(defaultDocsMarkdown)
	}
	s.renderPage(w, "docs", pageData{
		Title:       "Quick start & docs — Lotus AI (lotusmarket)",
		Description: "Cài đặt và sử dụng thư viện lotusmarket cho Go và Python: fetcher, technical, ratings, backtest, regime classifier. MIT license, no API key.",
		Canonical:   s.cfg.BaseURL + "/docs",
		PageType:    "docs",
		Body:        renderMarkdown(string(body)),
	})
}

// --- About (/about) ---

func (s *server) handleAbout(w http.ResponseWriter, _ *http.Request) {
	aboutPath := filepath.Join(s.cfg.ContentDir, "about.md")
	body, err := os.ReadFile(aboutPath)
	if err != nil {
		body = []byte(defaultAboutMarkdown)
	}
	s.renderPage(w, "about", pageData{
		Title:       "About — Lotus AI & lotusmarket project",
		Description: "Câu chuyện đằng sau lotusmarket: một lập trình viên + nhà đầu tư cá nhân build công cụ phân tích cổ phiếu VN open source vì retail VN deserve công cụ tốt hơn hot pick.",
		Canonical:   s.cfg.BaseURL + "/about",
		PageType:    "about",
		Body:        renderMarkdown(string(body)),
	})
}

// --- 404 ---

func (s *server) handleNotFound(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	s.renderPage(w, "not-found", pageData{
		Title:       "404 — Lotus AI",
		Description: "Trang không tồn tại.",
		Canonical:   s.cfg.BaseURL + "/",
	})
}
