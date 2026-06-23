package main

import (
	_ "embed"
	"net/http"
	"strings"
)

//go:embed assets/og-default.png
var ogDefaultPNG []byte

// handleOGImage serves dynamic 1200x630 PNGs for og:image previews.
//
//	/og/default.png       → generic Lotus AI card
//	/og/blog/{slug}.png   → per-post card with the post title
//
// Falls back to /static/og-default.png on render failure so social scrapers
// always receive a valid image.
func (s *server) handleOGImage(w http.ResponseWriter, r *http.Request) {
	rest := strings.TrimPrefix(r.URL.Path, "/og/")
	rest = strings.TrimSuffix(rest, ".png")

	var params ogParams
	cacheKey := rest

	switch {
	case rest == "" || rest == "default":
		params = ogParams{
			Title:    "Lotus AI — VN Stock Market Toolkit",
			Subtitle: "Phân tích cổ phiếu Việt Nam · MIT open source",
		}
		cacheKey = "default"
	case strings.HasPrefix(rest, "blog/"):
		slug := strings.TrimPrefix(rest, "blog/")
		var match *blogPost
		for _, p := range s.loadBlogPosts() {
			if p.Slug == slug {
				p := p
				match = &p
				break
			}
		}
		if match == nil {
			s.serveOGFallback(w, r)
			return
		}
		params = ogParams{
			Title:    match.Title,
			Subtitle: "Lotus AI · " + formatVnDate(match.Date),
		}
	default:
		s.serveOGFallback(w, r)
		return
	}

	png, err := ogImageFor(cacheKey, func() (ogParams, error) { return params, nil })
	if err != nil {
		s.serveOGFallback(w, r)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	_, _ = w.Write(png)
}

// serveOGFallback writes the embedded default PNG directly. Shipped inside the
// binary so a render bug or missing asset can never break social previews.
func (s *server) serveOGFallback(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	_, _ = w.Write(ogDefaultPNG)
}
