package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestGenerateStaticOGDefault writes content/static/og-default.png as a build
// artifact. Run with: go test ./cmd/lotusai-site/ -run TestGenerateStaticOGDefault -gen-og
// It's gated by an env flag so it doesn't pollute normal test runs.
func TestGenerateStaticOGDefault(t *testing.T) {
	if os.Getenv("GEN_OG") != "1" {
		t.Skip("set GEN_OG=1 to (re)generate static og-default.png")
	}
	png, err := renderOGImage(ogParams{
		Title:    "Lotus AI — VN Stock Market Toolkit",
		Subtitle: "Phân tích cổ phiếu Việt Nam · MIT open source",
	})
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	out := filepath.Join("assets", "og-default.png")
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(out, png, 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	t.Logf("wrote %s (%d bytes)", out, len(png))
}

// TestRenderOGImageBasic ensures the renderer produces a non-trivial PNG. Runs
// in normal test mode — guards against future regressions in the font/embed.
func TestRenderOGImageBasic(t *testing.T) {
	png, err := renderOGImage(ogParams{
		Title:    "Thị trường hôm nay đang ở regime nào? Phân loại deterministic + ý nghĩa",
		Subtitle: "Lotus AI · 19/05/2026",
	})
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	if len(png) < 1000 {
		t.Fatalf("png too small: %d bytes", len(png))
	}
	if png[0] != 0x89 || png[1] != 'P' || png[2] != 'N' || png[3] != 'G' {
		t.Fatalf("not a PNG header: % x", png[:4])
	}
}
