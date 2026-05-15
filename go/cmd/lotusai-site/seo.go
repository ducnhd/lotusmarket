package main

import (
	"encoding/json"
	"fmt"
	"html/template"
)

// jsonLDOrganization — schema.org WebSite + Organization, used on home page
// to help Google understand the brand + provide site-name and search box.
func jsonLDOrganization(baseURL string) template.JS {
	doc := map[string]any{
		"@context": "https://schema.org",
		"@graph": []any{
			map[string]any{
				"@type":       "Organization",
				"@id":         baseURL + "#org",
				"name":        "Lotus AI / lotusmarket",
				"url":         baseURL + "/",
				"sameAs":      []string{"https://github.com/ducnhd/lotusmarket", "https://ducnhd.github.io/lotusmarket/", "https://pypi.org/project/lotusmarket/", "https://t.me/vnlotusmarket"},
				"description": "Vietnamese stock market open-source toolkit — free, deterministic, no hot picks.",
				"logo": map[string]any{
					"@type": "ImageObject",
					"url":   baseURL + "/static/og-default.png",
				},
			},
			map[string]any{
				"@type":       "WebSite",
				"@id":         baseURL + "#website",
				"url":         baseURL + "/",
				"name":        "Lotus AI",
				"description": "Toolkit phân tích cổ phiếu Việt Nam — miễn phí, open source, không hot picks.",
				"publisher":   map[string]any{"@id": baseURL + "#org"},
				"inLanguage":  "vi-VN",
			},
			map[string]any{
				"@type":               "SoftwareApplication",
				"name":                "lotusmarket",
				"url":                 "https://github.com/ducnhd/lotusmarket",
				"operatingSystem":     "Linux, macOS, Windows",
				"applicationCategory": "FinanceApplication",
				"offers": map[string]any{
					"@type":         "Offer",
					"price":         "0",
					"priceCurrency": "USD",
				},
				"description": "Vietnamese stock market analytics library for Go and Python. Real-time quotes, technical indicators, sentiment, regime classifier, backtest harness.",
				"license":     "https://opensource.org/licenses/MIT",
			},
		},
	}
	b, _ := json.Marshal(doc)
	return template.JS(b)
}

// jsonLDArticle — schema.org BlogPosting for blog post pages.
func jsonLDArticle(baseURL string, p *blogPost) template.JS {
	doc := map[string]any{
		"@context":      "https://schema.org",
		"@type":         "BlogPosting",
		"headline":      p.Title,
		"datePublished": p.Date.Format("2006-01-02"),
		"dateModified":  p.Date.Format("2006-01-02"),
		"author": map[string]any{
			"@type": "Person",
			"name":  "ducnhd",
			"url":   "https://github.com/ducnhd",
		},
		"publisher": map[string]any{
			"@type": "Organization",
			"name":  "Lotus AI",
			"logo": map[string]any{
				"@type": "ImageObject",
				"url":   baseURL + "/static/og-default.png",
			},
		},
		"description":      p.Excerpt,
		"mainEntityOfPage": fmt.Sprintf("%s/blog/%s", baseURL, p.Slug),
		"inLanguage":       "vi-VN",
	}
	b, _ := json.Marshal(doc)
	return template.JS(b)
}
