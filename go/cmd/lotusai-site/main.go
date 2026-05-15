// lotusai-site — SEO-friendly marketing website for the lotusmarket open-source
// toolkit. Single Go binary, server-side rendered HTML, no SPA, no client JS
// for content. Designed to run on a Raspberry Pi behind nginx + Let's Encrypt
// at https://lotusai.servehttp.com.
//
// Pages:
//
//	/                   landing (hero + live snapshot + features + CTA)
//	/dashboard          full live market data
//	/blog/              blog post list
//	/blog/{slug}        single blog post rendered from markdown
//	/docs/              install + quick start
//	/about              author bio + project background
//
// SEO endpoints:
//
//	/sitemap.xml        auto-generated from pages + blog posts
//	/robots.txt         points crawlers at sitemap
//	/feed.xml           RSS feed of blog + (optional) daily reports
//	/health             liveness probe for systemd / nginx
//
// Live data is pulled from the lotusmarket library and cached in-memory for
// 10 minutes — keeps Pi load low and avoids hammering VPS/Yahoo APIs.
package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	addr       = flag.String("addr", ":8083", "listen address")
	contentDir = flag.String("content", "./content", "directory containing blog markdown + static assets")
	baseURL    = flag.String("base-url", "https://lotusai.servehttp.com", "canonical base URL for sitemap/OG tags")
	cacheTTL   = flag.Duration("cache-ttl", 10*time.Minute, "TTL for live data cache")
)

func main() {
	flag.Parse()

	srv := newServer(serverConfig{
		BaseURL:    *baseURL,
		ContentDir: *contentDir,
		CacheTTL:   *cacheTTL,
	})

	mux := http.NewServeMux()
	srv.routes(mux)

	server := &http.Server{
		Addr:              *addr,
		Handler:           withMiddleware(mux),
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("[lotusai-site] listening on %s (base=%s, content=%s)", *addr, *baseURL, *contentDir)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server: %v", err)
	}
	_ = os.Exit
}
