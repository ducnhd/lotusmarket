package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// writePosts creates empty .md files in dir for each "YYYY-MM-DD-<rest>" name.
func writePosts(t *testing.T, dir string, names ...string) {
	t.Helper()
	for _, n := range names {
		if err := os.WriteFile(filepath.Join(dir, n+".md"), []byte("---\n---\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func TestTopicOf(t *testing.T) {
	cases := map[string]string{
		"tech-classic-golden-cross": "tech-classic",
		"myth-c-tc-cao":             "myth-buster",
		"random-vhm":                "random-ticker",
		"career-fire-math-vn":       "career-investing",
		"supply-chain-hpg":          "supply-chain",
		"psychology-loss-aversion":  "psychology",
		"insight-joint-factor":      "data-insight",
		"compare-vnm-vs-msn":        "comparison",
		"regime-now-06-22":          "regime-now",
		"news-impact-06-21":         "news-impact",
	}
	for rest, want := range cases {
		if got := topicOf(rest); got != want {
			t.Errorf("topicOf(%q) = %q, want %q", rest, got, want)
		}
	}
}

func TestContentKey(t *testing.T) {
	// Timely topics collapse to base key regardless of date suffix.
	if got := contentKey("regime-now-06-22"); got != "regime-now" {
		t.Errorf("regime contentKey = %q, want regime-now", got)
	}
	if got := contentKey("news-impact-06-21"); got != "news-impact" {
		t.Errorf("news contentKey = %q, want news-impact", got)
	}
	// Catalog variants pass through unchanged.
	if got := contentKey("compare-vnm-vs-msn"); got != "compare-vnm-vs-msn" {
		t.Errorf("compare contentKey = %q, want compare-vnm-vs-msn", got)
	}
}

func TestRecentContentKeys(t *testing.T) {
	dir := t.TempDir()
	today := time.Now().Format("2006-01-02")
	old := time.Now().AddDate(0, 0, -90).Format("2006-01-02")
	writePosts(t, dir,
		today+"-compare-vnm-vs-msn",
		today+"-regime-now-"+time.Now().Format("01-02"),
		old+"-insight-joint-factor", // outside 60d window
	)
	keys := recentContentKeys(dir, 60)
	if _, ok := keys["compare-vnm-vs-msn"]; !ok {
		t.Error("expected compare-vnm-vs-msn in recent keys")
	}
	if _, ok := keys["regime-now"]; !ok {
		t.Error("expected regime-now (date-stripped) in recent keys")
	}
	if _, ok := keys["insight-joint-factor"]; ok {
		t.Error("90-day-old post should be outside the 60d window")
	}
}

func TestPickFreshAvoidsRecent(t *testing.T) {
	slugs := []string{"a", "b", "c"}
	// "a" and "b" used recently; only "c" is fresh → must pick index 2.
	recentVariantKeys = map[string]int{"a": 1, "b": 2}
	for i := 0; i < 50; i++ {
		if got := pickFresh(slugs); got != 2 {
			t.Fatalf("pickFresh chose %d, expected 2 (only fresh slug)", got)
		}
	}
}

func TestPickFreshFallsBackToOldest(t *testing.T) {
	slugs := []string{"a", "b", "c"}
	// All used; "c" is oldest (largest days-ago) → must be chosen.
	recentVariantKeys = map[string]int{"a": 1, "b": 2, "c": 40}
	for i := 0; i < 50; i++ {
		if got := pickFresh(slugs); got != 2 {
			t.Fatalf("pickFresh chose %d, expected 2 (oldest slug)", got)
		}
	}
}

func TestSelectTopicHonorsCooldown(t *testing.T) {
	dir := t.TempDir()
	// Publish every topic except news-impact inside the cooldown window.
	day := time.Now()
	names := []string{}
	for i, rest := range []string{
		"tech-classic-x", "myth-y", "random-vhm", "career-x",
		"supply-chain-hpg", "psychology-x", "insight-x", "compare-x", "regime-now-01-02",
	} {
		d := day.AddDate(0, 0, -i).Format("2006-01-02")
		names = append(names, d+"-"+rest)
	}
	writePosts(t, dir, names...)
	recentVariantKeys = recentContentKeys(dir, variantCooldownDays)
	// Only news-impact is outside the recent set → must be selected.
	got := selectTopic(dir, "auto")
	if got.Key != "news-impact" {
		t.Errorf("selectTopic = %q, want news-impact (only eligible topic)", got.Key)
	}
}

var _ = fmt.Sprintf
