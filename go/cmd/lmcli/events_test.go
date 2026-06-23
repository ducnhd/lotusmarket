package main

import (
	"testing"
	"time"
)

// Task 1: empty snapshot produces no events.
func TestDetectEventsEmpty(t *testing.T) {
	got := DetectEvents(Snapshot{})
	if len(got) != 0 {
		t.Fatalf("empty snapshot: got %d events, want 0", len(got))
	}
}

// Task 2: mover detection + scoring.
func TestDetectMover(t *testing.T) {
	s := Snapshot{Stocks: []StockSnap{
		{Ticker: "HPG", ChangePct: -6.0, Close: 25000, Volume: 30e6, AvgVol20: 10e6, RSI: 38, MATrend: "downtrend"},
		{Ticker: "VCB", ChangePct: 1.2, Close: 90000, Volume: 5e6, AvgVol20: 6e6, RSI: 55, MATrend: "uptrend"},
	}}
	got := detectMovers(s)
	if len(got) != 1 {
		t.Fatalf("got %d movers, want 1 (only HPG >5%%)", len(got))
	}
	e := got[0]
	if e.Ticker != "HPG" || e.Type != "mover" {
		t.Fatalf("got %+v, want HPG mover", e)
	}
	// -6% is mid-band → score in (50,85); volume 3× → +10 bonus.
	if e.Score < 65 || e.Score > 90 {
		t.Errorf("HPG score = %.1f, want 65-90", e.Score)
	}
	if e.DataBlock == "" {
		t.Error("DataBlock must be populated")
	}
}

// Task 3: volume surge detection.
func TestDetectVolume(t *testing.T) {
	s := Snapshot{Stocks: []StockSnap{
		{Ticker: "SSI", ChangePct: 2.0, Close: 30000, Volume: 40e6, AvgVol20: 10e6, RSI: 60, MATrend: "uptrend"},    // 4× surge, small move
		{Ticker: "HPG", ChangePct: -6.0, Close: 25000, Volume: 30e6, AvgVol20: 10e6, RSI: 38, MATrend: "downtrend"}, // mover, excluded here
	}}
	got := detectVolume(s)
	if len(got) != 1 || got[0].Ticker != "SSI" {
		t.Fatalf("got %+v, want only SSI volume event", got)
	}
	if got[0].Type != "volume" || got[0].Score < 40 || got[0].Score > 70 {
		t.Errorf("SSI volume event = %+v, want type=volume score 40-70", got[0])
	}
}

// Task 4: flow detection.
func TestDetectFlow(t *testing.T) {
	if got := detectFlow(Snapshot{FlowBuyPressure: 52}); len(got) != 0 {
		t.Fatalf("neutral flow: got %d, want 0", len(got))
	}
	got := detectFlow(Snapshot{FlowBuyPressure: 64})
	if len(got) != 1 || got[0].Type != "flow" || got[0].Ticker != "" {
		t.Fatalf("got %+v, want one market-wide flow event", got)
	}
	if got[0].Score < 40 || got[0].Score > 60 {
		t.Errorf("flow score = %.1f, want 40-60", got[0].Score)
	}
}

// Task 5: global shock detection.
func TestDetectGlobal(t *testing.T) {
	calm := Snapshot{Globals: []GlobalSnap{{Symbol: "^VIX", Price: 14, ChangePct: 1}, {Symbol: "^GSPC", ChangePct: 0.3}}}
	if got := detectGlobal(calm); len(got) != 0 {
		t.Fatalf("calm: got %d, want 0", len(got))
	}
	shock := Snapshot{Globals: []GlobalSnap{{Symbol: "^VIX", Price: 32, ChangePct: 20}, {Symbol: "^GSPC", ChangePct: -2.5}}}
	got := detectGlobal(shock)
	if len(got) == 0 || got[0].Type != "global" {
		t.Fatalf("shock: got %+v, want a global event", got)
	}
	if got[0].Score < 45 {
		t.Errorf("global score = %.1f, want >=45", got[0].Score)
	}
}

// Task 6: insider detection.
func TestDetectInsider(t *testing.T) {
	s := Snapshot{Insiders: []InsiderSnap{
		{Ticker: "MWG", Side: "buy", Shares: 3_000_000, InsiderName: "CEO"},
		{Ticker: "FPT", Side: "sell", Shares: 200_000, InsiderName: "BoD"}, // below threshold
	}}
	got := detectInsider(s)
	if len(got) != 1 || got[0].Ticker != "MWG" || got[0].Type != "insider" {
		t.Fatalf("got %+v, want one MWG insider event", got)
	}
	if got[0].Score < 45 || got[0].Score > 65 {
		t.Errorf("insider score = %.1f, want 45-65", got[0].Score)
	}
}

// Task 7: ranking + threshold gate + tie-break.
func TestTopEventThreshold(t *testing.T) {
	if _, ok := TopEvent([]Event{{Type: "flow", Score: 50}}); ok {
		t.Error("score 50 < 55 must not pass the gate")
	}
}

func TestTopEventTieBreakPrefersNamed(t *testing.T) {
	evs := []Event{
		{Type: "global", Ticker: "", Score: 60},
		{Type: "mover", Ticker: "HPG", Score: 60},
	}
	got, ok := TopEvent(evs)
	if !ok || got.Ticker != "HPG" {
		t.Fatalf("got %+v ok=%v, want HPG (named beats market-wide at equal score)", got, ok)
	}
}

func TestTopEventHighestScoreWins(t *testing.T) {
	evs := []Event{{Type: "mover", Ticker: "HPG", Score: 60}, {Type: "global", Score: 75}}
	got, _ := TopEvent(evs)
	if got.Score != 75 {
		t.Fatalf("got score %.0f, want 75 (highest wins regardless of type)", got.Score)
	}
}

// Task 8: event slug + dedup.
func TestEventSlug(t *testing.T) {
	if s := eventSlug(Event{Type: "mover", Ticker: "HPG"}); s != "event-mover-hpg" {
		t.Errorf("got %q, want event-mover-hpg", s)
	}
	if s := eventSlug(Event{Type: "flow"}); s != "event-flow" {
		t.Errorf("got %q, want event-flow", s)
	}
}

func TestNotRecentlyPosted(t *testing.T) {
	recentVariantKeys = map[string]int{"event-mover-hpg": 2}
	if notRecentlyPosted(Event{Type: "mover", Ticker: "HPG"}) {
		t.Error("HPG mover posted 2 days ago must be considered recent")
	}
	if !notRecentlyPosted(Event{Type: "mover", Ticker: "VCB"}) {
		t.Error("VCB never posted must be fresh")
	}
}

// Task 9: backfill decision helpers.
func TestShouldBackfill(t *testing.T) {
	if shouldBackfill(2) {
		t.Error("2 days < 3 must not backfill")
	}
	if !shouldBackfill(3) {
		t.Error("3 days >= 3 must backfill")
	}
}

func TestDaysSinceLastPost(t *testing.T) {
	dir := t.TempDir()
	now := time.Date(2026, 6, 23, 12, 0, 0, 0, time.UTC)
	// writePosts from autoblog_test.go takes names WITHOUT the .md suffix.
	writePosts(t, dir, "2026-06-20-event-mover-hpg")
	if d := daysSinceLastPost(dir, now); d != 3 {
		t.Fatalf("got %d days, want 3", d)
	}
	if d := daysSinceLastPost(t.TempDir(), now); d < 9000 {
		t.Errorf("empty dir should be a large sentinel, got %d", d)
	}
}
