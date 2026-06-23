# Event-driven Market Blog Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the fixed-rotation auto-blog with an event-triggered generator that posts only when a real market event occurs (named ticker + number + date), plus a bot-filtered view counter to measure readership.

**Architecture:** A new `lmcli marketscan` subcommand fetches VN30+HNX30 quotes, global indices, and insider data into a `Snapshot`, runs a **pure** `DetectEvents` scorer, gates on a threshold, dedups against recent posts, and feeds the top event to Claude with an event-mode prompt. If no event clears the threshold and ≥3 days have passed since the last post, it falls back to the existing evergreen generator. A separate view counter in `lotusai-site` records non-bot reads.

**Tech Stack:** Go (stdlib + existing `fetchers`/`technical`/`ratings`/`ai`/`sentiment` packages). Tests are table-driven and fixture-based (no network).

**Scope:** This plan covers Phase 1 (data-only events) + the view counter + cron change. Phase 2 (CafeF RSS headlines as the "why" layer) is a separate plan written later.

---

## File Structure

| File | Responsibility |
|---|---|
| `go/cmd/lmcli/marketscan.go` | `runMarketScan` glue (fetch → detect → gate → write/backfill) + `buildSnapshot` (I/O) |
| `go/cmd/lmcli/events.go` | **Pure** types (`Snapshot`, `StockSnap`, `Event`) + `DetectEvents` + per-type scorers + ranking/threshold + dedup + backfill decision |
| `go/cmd/lmcli/events_test.go` | Table-driven tests for everything in `events.go` |
| `go/cmd/lmcli/universe.go` | `hnx30` list + `marketUniverse` helper |
| `go/cmd/lmcli/main.go` | Add `marketscan` subcommand dispatch (modify) |
| `go/cmd/lmcli/autoblog.go` | Reuse `matchCohort`, `safeSlug`, recency helpers; add event-mode prompt constant (modify) |
| `.github/workflows/auto-blog.yml` | Switch to post-close schedule + `marketscan` (modify) |
| `go/cmd/lotusai-site/views.go` | `viewStore`: bot-filtered per-slug counter, JSON persistence |
| `go/cmd/lotusai-site/views_test.go` | Tests for bot filtering + increment + persistence |
| `go/cmd/lotusai-site/server.go` | Wire `viewStore`, count in `handleBlog`, add `/blog/stats` (modify) |
| `go/cmd/lotusai-site/main.go` | Add `--data` flag (modify) |

Pure detection (`events.go`) is isolated from I/O (`marketscan.go`) so all scoring logic is unit-testable without network.

---

## Task 1: Pure types + empty DetectEvents

**Files:**
- Create: `go/cmd/lmcli/events.go`
- Test: `go/cmd/lmcli/events_test.go`

- [ ] **Step 1: Write the failing test**

```go
package main

import "testing"

func TestDetectEventsEmpty(t *testing.T) {
	got := DetectEvents(Snapshot{})
	if len(got) != 0 {
		t.Fatalf("empty snapshot: got %d events, want 0", len(got))
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd go && go test ./cmd/lmcli/ -run TestDetectEventsEmpty`
Expected: FAIL — `undefined: DetectEvents` / `Snapshot`.

- [ ] **Step 3: Write minimal implementation**

```go
package main

import "time"

// StockSnap is one ticker's same-day state, pre-computed for detection.
type StockSnap struct {
	Ticker    string
	Close     float64
	ChangePct float64
	Volume    int64
	AvgVol20  float64
	RSI       float64
	MATrend   string // uptrend|downtrend|mixed
}

// GlobalSnap is one global symbol's same-day state.
type GlobalSnap struct {
	Symbol    string
	Price     float64
	ChangePct float64
}

// InsiderSnap is a notable insider transaction (shares registered).
type InsiderSnap struct {
	Ticker     string
	Side       string // "buy" | "sell"
	Shares     int64
	InsiderName string
}

// Snapshot is the full pure input to DetectEvents (no network handles).
type Snapshot struct {
	Date            time.Time
	Stocks          []StockSnap
	Globals         []GlobalSnap
	Insiders        []InsiderSnap
	FlowBuyPressure float64 // aggregate bid/(bid+ask) %, 0-100
}

// Event is a scored, candidate news angle.
type Event struct {
	Type      string // mover|volume|flow|insider|global
	Ticker    string // "" for market-wide
	Magnitude float64
	Score     float64 // 0-100
	DataBlock string  // anchor numbers + cohort, fed to Claude
}

// DetectEvents runs every per-type detector over the snapshot and returns the
// union of candidate events (unranked, ungated).
func DetectEvents(s Snapshot) []Event {
	var out []Event
	return out
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd go && go test ./cmd/lmcli/ -run TestDetectEventsEmpty`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add go/cmd/lmcli/events.go go/cmd/lmcli/events_test.go
git commit -m "feat(marketscan): pure event types + empty DetectEvents"
```

---

## Task 2: `mover` detection + scoring

A big single-stock move: `|ChangePct| >= moverThresholdPct` (5%). Score scales 50→85
across 5%→7% (HOSE limit), with a +10 bonus when volume ≥ 2× AvgVol20 (capped 100).

**Files:**
- Modify: `go/cmd/lmcli/events.go`
- Test: `go/cmd/lmcli/events_test.go`

- [ ] **Step 1: Write the failing test**

```go
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd go && go test ./cmd/lmcli/ -run TestDetectMover`
Expected: FAIL — `undefined: detectMovers`.

- [ ] **Step 3: Write minimal implementation**

Add to `events.go`:

```go
import (
	"fmt"
	"math"
	"strings"
	"time"
)

const (
	moverThresholdPct = 5.0
	moverLimitPct     = 7.0 // HOSE daily limit
)

func detectMovers(s Snapshot) []Event {
	var out []Event
	for _, st := range s.Stocks {
		mag := math.Abs(st.ChangePct)
		if mag < moverThresholdPct {
			continue
		}
		// 5%→50, 7%→85, linear, clamped.
		frac := (mag - moverThresholdPct) / (moverLimitPct - moverThresholdPct)
		if frac > 1 {
			frac = 1
		}
		score := 50 + frac*35
		if st.AvgVol20 > 0 && float64(st.Volume) >= 2*st.AvgVol20 {
			score += 10
		}
		if score > 100 {
			score = 100
		}
		out = append(out, Event{
			Type:      "mover",
			Ticker:    st.Ticker,
			Magnitude: st.ChangePct,
			Score:     score,
			DataBlock: stockDataBlock(st),
		})
	}
	return out
}

// stockDataBlock renders the anchor numbers + historical cohort for a ticker.
func stockDataBlock(st StockSnap) string {
	volMult := 0.0
	if st.AvgVol20 > 0 {
		volMult = float64(st.Volume) / st.AvgVol20
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "Mã: %s\n", st.Ticker)
	fmt.Fprintf(&sb, "Giá đóng cửa: %.0f VND\n", st.Close)
	fmt.Fprintf(&sb, "Thay đổi hôm nay: %+.2f%%\n", st.ChangePct)
	fmt.Fprintf(&sb, "Volume: %d (%.1f× MA20 volume)\n", st.Volume, volMult)
	fmt.Fprintf(&sb, "RSI(14): %.1f\n", st.RSI)
	fmt.Fprintf(&sb, "MA trend: %s\n\n", st.MATrend)
	fmt.Fprintf(&sb, "Cohort lịch sử khớp:\n%s\n", matchCohort(st.MATrend, st.RSI))
	return sb.String()
}
```

Wire into `DetectEvents`:

```go
func DetectEvents(s Snapshot) []Event {
	var out []Event
	out = append(out, detectMovers(s)...)
	return out
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd go && go test ./cmd/lmcli/ -run 'TestDetectMover|TestDetectEventsEmpty'`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add go/cmd/lmcli/events.go go/cmd/lmcli/events_test.go
git commit -m "feat(marketscan): mover detection + scoring + cohort data block"
```

---

## Task 3: `volume` detection + scoring

Volume surge without an extreme price move: `Volume >= volumeSurgeMult × AvgVol20`
(2×) **and** `|ChangePct| < moverThresholdPct` (so it does not double-count a mover).
Score 40→70 across 2×→5×.

**Files:** Modify `go/cmd/lmcli/events.go`; Test `go/cmd/lmcli/events_test.go`

- [ ] **Step 1: Write the failing test**

```go
func TestDetectVolume(t *testing.T) {
	s := Snapshot{Stocks: []StockSnap{
		{Ticker: "SSI", ChangePct: 2.0, Close: 30000, Volume: 40e6, AvgVol20: 10e6, RSI: 60, MATrend: "uptrend"},   // 4× surge, small move
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd go && go test ./cmd/lmcli/ -run TestDetectVolume`
Expected: FAIL — `undefined: detectVolume`.

- [ ] **Step 3: Write minimal implementation**

```go
const volumeSurgeMult = 2.0

func detectVolume(s Snapshot) []Event {
	var out []Event
	for _, st := range s.Stocks {
		if st.AvgVol20 <= 0 || math.Abs(st.ChangePct) >= moverThresholdPct {
			continue
		}
		mult := float64(st.Volume) / st.AvgVol20
		if mult < volumeSurgeMult {
			continue
		}
		// 2×→40, 5×→70, clamped.
		frac := (mult - volumeSurgeMult) / (5 - volumeSurgeMult)
		if frac > 1 {
			frac = 1
		}
		out = append(out, Event{
			Type:      "volume",
			Ticker:    st.Ticker,
			Magnitude: mult,
			Score:     40 + frac*30,
			DataBlock: stockDataBlock(st),
		})
	}
	return out
}
```

Add to `DetectEvents`: `out = append(out, detectVolume(s)...)`.

- [ ] **Step 4: Run test to verify it passes**

Run: `cd go && go test ./cmd/lmcli/ -run TestDetectVolume`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add go/cmd/lmcli/events.go go/cmd/lmcli/events_test.go
git commit -m "feat(marketscan): volume-surge detection"
```

---

## Task 4: `flow` detection + scoring

Market-wide aggregate buy-pressure extreme: `FlowBuyPressure >= 60` or `<= 40`.
Score 40→60 across 60→70 (and symmetric 40→30). Ticker is "".

**Files:** Modify `go/cmd/lmcli/events.go`; Test `go/cmd/lmcli/events_test.go`

- [ ] **Step 1: Write the failing test**

```go
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd go && go test ./cmd/lmcli/ -run TestDetectFlow`
Expected: FAIL — `undefined: detectFlow`.

- [ ] **Step 3: Write minimal implementation**

```go
func detectFlow(s Snapshot) []Event {
	dev := math.Abs(s.FlowBuyPressure - 50)
	if dev < 10 { // within 40-60 is neutral
		return nil
	}
	// dev 10→40, dev 20→60, clamped.
	frac := (dev - 10) / 10
	if frac > 1 {
		frac = 1
	}
	dir := "mua mạnh"
	if s.FlowBuyPressure < 50 {
		dir = "bán mạnh"
	}
	db := fmt.Sprintf("Dòng tiền nội tổng VN30+HNX30: buy-pressure %.1f%% (%s)\n\n"+
		"Cohort theo regime (9 năm):\n- CRISIS: +8.53%% fwd 60d, win 67%%\n"+
		"- STABLE: +3.20%%, 48%% (baseline +4.75%%)\n", s.FlowBuyPressure, dir)
	return []Event{{Type: "flow", Magnitude: s.FlowBuyPressure, Score: 40 + frac*20, DataBlock: db}}
}
```

Add to `DetectEvents`: `out = append(out, detectFlow(s)...)`.

- [ ] **Step 4: Run test to verify it passes**

Run: `cd go && go test ./cmd/lmcli/ -run TestDetectFlow`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add go/cmd/lmcli/events.go go/cmd/lmcli/events_test.go
git commit -m "feat(marketscan): aggregate flow detection"
```

---

## Task 5: `global` detection + scoring

Global shock spilling into VN: VIX `> 25` (and `>30` stronger), or `|ChangePct| >= 2`
on `^GSPC`/`^HSI`. Score 45→75.

**Files:** Modify `go/cmd/lmcli/events.go`; Test `go/cmd/lmcli/events_test.go`

- [ ] **Step 1: Write the failing test**

```go
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd go && go test ./cmd/lmcli/ -run TestDetectGlobal`
Expected: FAIL — `undefined: detectGlobal`.

- [ ] **Step 3: Write minimal implementation**

```go
func detectGlobal(s Snapshot) []Event {
	var vix, sp, hsi float64
	have := map[string]bool{}
	for _, g := range s.Globals {
		switch g.Symbol {
		case "^VIX":
			vix = g.Price
			have["vix"] = true
		case "^GSPC":
			sp = g.ChangePct
		case "^HSI":
			hsi = g.ChangePct
		}
	}
	score := 0.0
	switch {
	case have["vix"] && vix > 30:
		score = 75
	case have["vix"] && vix > 25:
		score = 60
	}
	if a := math.Max(math.Abs(sp), math.Abs(hsi)); a >= 2 {
		s2 := 45 + math.Min((a-2)/2, 1)*25 // 2%→45, 4%→70
		if s2 > score {
			score = s2
		}
	}
	if score == 0 {
		return nil
	}
	db := fmt.Sprintf("Trạng thái global hôm nay:\nVIX: %.1f\nS&P 500: %+.2f%%\nHang Seng: %+.2f%%\n\n"+
		"Cohort lan toả global → VN:\n- S&P -2%%+ → VN-Index thường -1.5%% phiên sau (foreign flow)\n"+
		"- VIX>25 → VN regime shift VOLATILE/CRISIS_PANIC, fwd 60d historically +8.53%%\n", vix, sp, hsi)
	return []Event{{Type: "global", Magnitude: vix, Score: score, DataBlock: db}}
}
```

Add to `DetectEvents`: `out = append(out, detectGlobal(s)...)`.

- [ ] **Step 4: Run test to verify it passes**

Run: `cd go && go test ./cmd/lmcli/ -run TestDetectGlobal`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add go/cmd/lmcli/events.go go/cmd/lmcli/events_test.go
git commit -m "feat(marketscan): global-shock detection"
```

---

## Task 6: `insider` detection + scoring

A notable insider transaction: `Shares >= insiderShareThreshold` (1,000,000).
Score 45→65 across 1M→5M shares.

**Files:** Modify `go/cmd/lmcli/events.go`; Test `go/cmd/lmcli/events_test.go`

- [ ] **Step 1: Write the failing test**

```go
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd go && go test ./cmd/lmcli/ -run TestDetectInsider`
Expected: FAIL — `undefined: detectInsider`.

- [ ] **Step 3: Write minimal implementation**

```go
const insiderShareThreshold = 1_000_000

func detectInsider(s Snapshot) []Event {
	var out []Event
	for _, in := range s.Insiders {
		if in.Shares < insiderShareThreshold {
			continue
		}
		frac := math.Min(float64(in.Shares-insiderShareThreshold)/4_000_000, 1) // 1M→0, 5M→1
		db := fmt.Sprintf("Giao dịch nội bộ: %s %s %d cổ phiếu (%s)\n\n"+
			"Bối cảnh: insider mua lớn thường là tín hiệu tin tưởng nội bộ; bán lớn có thể là chốt lời hoặc nhu cầu cá nhân — không tự động hàm ý xấu/tốt.\n",
			in.Ticker, in.Side, in.Shares, in.InsiderName)
		out = append(out, Event{Type: "insider", Ticker: in.Ticker, Magnitude: float64(in.Shares), Score: 45 + frac*20, DataBlock: db})
	}
	return out
}
```

Add to `DetectEvents`: `out = append(out, detectInsider(s)...)`.

- [ ] **Step 4: Run test to verify it passes**

Run: `cd go && go test ./cmd/lmcli/ -run TestDetectInsider`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add go/cmd/lmcli/events.go go/cmd/lmcli/events_test.go
git commit -m "feat(marketscan): insider-transaction detection"
```

---

## Task 7: Ranking + threshold gate + tie-break

`TopEvent` returns the single highest-scoring event ≥ `eventThreshold` (55), or
`(Event{}, false)` if none. Tie-break: prefer a named single-stock event
(`mover`/`volume`/`insider`) over a market-wide one (`flow`/`global`).

**Files:** Modify `go/cmd/lmcli/events.go`; Test `go/cmd/lmcli/events_test.go`

- [ ] **Step 1: Write the failing test**

```go
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd go && go test ./cmd/lmcli/ -run TestTopEvent`
Expected: FAIL — `undefined: TopEvent`.

- [ ] **Step 3: Write minimal implementation**

```go
const eventThreshold = 55.0

func isNamed(t string) bool { return t == "mover" || t == "volume" || t == "insider" }

// TopEvent returns the most newsworthy event at or above the threshold.
func TopEvent(evs []Event) (Event, bool) {
	best := Event{}
	found := false
	for _, e := range evs {
		if e.Score < eventThreshold {
			continue
		}
		if !found || e.Score > best.Score ||
			(e.Score == best.Score && isNamed(e.Type) && !isNamed(best.Type)) {
			best = e
			found = true
		}
	}
	return best, found
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd go && go test ./cmd/lmcli/ -run TestTopEvent`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add go/cmd/lmcli/events.go go/cmd/lmcli/events_test.go
git commit -m "feat(marketscan): ranking, threshold gate, tie-break"
```

---

## Task 8: Event slug + dedup against recent posts

`eventSlug(e)` → `event-<type>-<ticker>` (lowercased; market-wide uses just the
type). `notRecentlyPosted` reuses the recency map from `autoblog.go`
(`recentVariantKeys`, populated by `recentContentKeys`).

**Files:** Modify `go/cmd/lmcli/events.go`; Test `go/cmd/lmcli/events_test.go`

- [ ] **Step 1: Write the failing test**

```go
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd go && go test ./cmd/lmcli/ -run 'TestEventSlug|TestNotRecentlyPosted'`
Expected: FAIL — `undefined: eventSlug` / `notRecentlyPosted`.

- [ ] **Step 3: Write minimal implementation**

```go
func eventSlug(e Event) string {
	if e.Ticker == "" {
		return "event-" + e.Type
	}
	return "event-" + e.Type + "-" + strings.ToLower(e.Ticker)
}

// notRecentlyPosted reports whether this event's slug is absent from the recent
// content map (populated by recentContentKeys in autoblog.go).
func notRecentlyPosted(e Event) bool {
	_, used := recentVariantKeys[eventSlug(e)]
	return !used
}
```

Note: `eventCooldownDays = 5` is enforced by populating `recentVariantKeys` with a
5-day window in Task 10 (`recentContentKeys(outDir, 5)`).

- [ ] **Step 4: Run test to verify it passes**

Run: `cd go && go test ./cmd/lmcli/ -run 'TestEventSlug|TestNotRecentlyPosted'`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add go/cmd/lmcli/events.go go/cmd/lmcli/events_test.go
git commit -m "feat(marketscan): event slug + recency dedup"
```

---

## Task 9: Backfill decision (pure)

`shouldBackfill(daysSinceLastPost)` → true when `>= backfillGapDays` (3).
`daysSinceLastPost(dir, now)` scans `docs/blog/` for the newest `YYYY-MM-DD-*.md`.

**Files:** Modify `go/cmd/lmcli/events.go`; Test `go/cmd/lmcli/events_test.go`

- [ ] **Step 1: Write the failing test**

```go
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
	mustWrite(t, dir, "2026-06-20-event-mover-hpg.md")
	if d := daysSinceLastPost(dir, now); d != 3 {
		t.Fatalf("got %d days, want 3", d)
	}
	if d := daysSinceLastPost(t.TempDir(), now); d < 9000 {
		t.Errorf("empty dir should be a large sentinel, got %d", d)
	}
}

// mustWrite creates an empty post file. Shared test helper.
func mustWrite(t *testing.T, dir, name string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte("---\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
}
```

Note: if `writePosts` from `autoblog_test.go` already exists, reuse it instead of
`mustWrite` to avoid a duplicate helper — pick whichever is already defined.

- [ ] **Step 2: Run test to verify it fails**

Run: `cd go && go test ./cmd/lmcli/ -run 'TestShouldBackfill|TestDaysSinceLastPost'`
Expected: FAIL — `undefined: shouldBackfill` / `daysSinceLastPost`.

- [ ] **Step 3: Write minimal implementation**

```go
import "os"
import "path/filepath"

const backfillGapDays = 3

func shouldBackfill(daysSinceLastPost int) bool { return daysSinceLastPost >= backfillGapDays }

// daysSinceLastPost returns whole days between now and the newest dated post in
// dir. Returns a large sentinel (99999) when no dated post exists.
func daysSinceLastPost(dir string, now time.Time) int {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 99999
	}
	var newest time.Time
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		name := strings.TrimSuffix(e.Name(), ".md")
		if len(name) < 11 {
			continue
		}
		d, err := time.Parse("2006-01-02", name[:10])
		if err != nil {
			continue
		}
		if d.After(newest) {
			newest = d
		}
	}
	if newest.IsZero() {
		return 99999
	}
	return int(now.Sub(newest).Hours() / 24)
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd go && go test ./cmd/lmcli/ -run 'TestShouldBackfill|TestDaysSinceLastPost'`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add go/cmd/lmcli/events.go go/cmd/lmcli/events_test.go
git commit -m "feat(marketscan): backfill decision helpers"
```

---

## Task 10: Universe — add HNX30

**Files:**
- Create: `go/cmd/lmcli/universe.go`
- Test: `go/cmd/lmcli/events_test.go`

- [ ] **Step 1: Write the failing test**

```go
func TestMarketUniverse(t *testing.T) {
	u := marketUniverse()
	if len(u) < 55 {
		t.Fatalf("universe = %d tickers, want ~60 (VN30+HNX30)", len(u))
	}
	seen := map[string]bool{}
	for _, tk := range u {
		if seen[tk] {
			t.Errorf("duplicate ticker in universe: %s", tk)
		}
		seen[tk] = true
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd go && go test ./cmd/lmcli/ -run TestMarketUniverse`
Expected: FAIL — `undefined: marketUniverse`.

- [ ] **Step 3: Write minimal implementation**

Create `universe.go`. **Before merging, verify `hnx30` against the current
official HNX30 constituents (https://hnx.vn) — index membership changes.**

```go
package main

// hnx30 — HNX30 index constituents. VERIFY against the current official list
// before merge; membership is reviewed periodically by HNX.
var hnx30 = []string{
	"SHS", "PVS", "CEO", "MBS", "IDC", "VCS", "HUT", "TNG", "PVI", "DTD",
	"L14", "NVB", "BVS", "TIG", "LAS", "VGS", "DXP", "NRC", "PVC", "TVC",
	"API", "IDV", "NTP", "PLC", "MST", "S99", "DDG", "AMV", "ART", "VC3",
}

// marketUniverse returns the dedup'd union of VN30 and HNX30 tickers.
func marketUniverse() []string {
	seen := map[string]bool{}
	out := []string{}
	for _, tk := range append(append([]string{}, vn30...), hnx30...) {
		if !seen[tk] {
			seen[tk] = true
			out = append(out, tk)
		}
	}
	return out
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd go && go test ./cmd/lmcli/ -run TestMarketUniverse`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add go/cmd/lmcli/universe.go go/cmd/lmcli/events_test.go
git commit -m "feat(marketscan): VN30+HNX30 universe"
```

---

## Task 11: Event-mode prompt constant

Add the event-mode prompt to `autoblog.go` (next to `blogPromptBase`), reusing the
shared discipline rules but with the event title formula + narrative structure.

**Files:** Modify `go/cmd/lmcli/autoblog.go`

- [ ] **Step 1: Add the constant** (no test — it is a string constant consumed in Task 12)

```go
const eventPromptBase = `Bạn là 1 quant analyst viết blog tiếng Việt cho Lotus AI / lotusmarket. KHÔNG tiết lộ mình là AI.

Đây là bài SỰ KIỆN — bám đúng diễn biến thị trường hôm nay trong DATA.

OUTPUT spec:
- Markdown tiếng Việt 700-1100 chữ. Mạch kể, có cao trào, KHÔNG khô khan.
- KHÔNG include YAML frontmatter (đã có sẵn). Bắt đầu từ heading hoặc body.
- TIÊU ĐỀ đã được set sẵn ở frontmatter — bài viết KHÔNG lặp lại nguyên văn tiêu đề ở dòng đầu.
- Lede (2-3 câu): nói NGAY chuyện gì xảy ra hôm nay + 1 câu tò mò mà data sẽ trả lời.
- Mạch bài: ① Chuyện gì (số liệu thật từ DATA) → ② Vì sao có thể (bối cảnh ngành/global/dòng tiền trong DATA) → ③ Cohort lịch sử tương tự đã đi tiếp ra sao → ④ Cần theo dõi gì.
- KHÔNG bịa số. CHỈ dùng số trong DATA. KHÔNG viết "có thể, có lẽ, dự đoán" kiểu mơ hồ — khẳng định từ data.
- KHÔNG khuyến nghị mua/bán. Phrasing: "data cho thấy cohort thắng X%", "lịch sử cho thấy...".
- Cuối bài: 1 đoạn "Verify reproducible" (pip install lotusmarket + 1-3 dòng code, hoặc 1 lệnh lmcli) + 1 dòng disclaimer không phải lời khuyên đầu tư.
- Link 1 lần: https://lotusai.servehttp.com hoặc https://github.com/ducnhd/lotusmarket.
- Tone: thẳng, không hype. Emoji tối đa 2. Ngôi đơn (mình/impersonal), KHÔNG "chúng ta".

DATA:
`
```

- [ ] **Step 2: Verify it compiles**

Run: `cd go && go build ./cmd/lmcli/`
Expected: success (constant unused yet is fine — it is referenced in Task 12; if the
linter flags unused, proceed directly to Task 12 in the same commit).

- [ ] **Step 3: Commit**

```bash
git add go/cmd/lmcli/autoblog.go
git commit -m "feat(marketscan): event-mode prompt"
```

---

## Task 12: Fetch + glue (`runMarketScan`, `buildSnapshot`)

The I/O layer: fetch the universe, compute per-ticker `StockSnap` (RSI, MA trend,
AvgVol20 via `EntradeHistory`), assemble `Snapshot`, run detection/gate/dedup, and
either write an event post or backfill. Thin glue — not unit-tested (detection is
already covered); validated by build + a live smoke run.

**Files:**
- Create: `go/cmd/lmcli/marketscan.go`

- [ ] **Step 1: Write the implementation**

```go
package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/ducnhd/lotusmarket/go/ai"
	"github.com/ducnhd/lotusmarket/go/fetchers"
	"github.com/ducnhd/lotusmarket/go/technical"
)

const eventCooldownDays = 5

func runMarketScan(ctx context.Context, outDir string) {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		log.Fatalf("mkdir %s: %v", outDir, err)
	}
	// Recency for event dedup (5-day window).
	recentVariantKeys = recentContentKeys(outDir, eventCooldownDays)

	snap := buildSnapshot(ctx)
	events := DetectEvents(snap)

	top, ok := TopEvent(events)
	for ok && !notRecentlyPosted(top) {
		// Drop the deduped top event and re-rank the rest.
		filtered := events[:0]
		for _, e := range events {
			if eventSlug(e) != eventSlug(top) {
				filtered = append(filtered, e)
			}
		}
		events = filtered
		top, ok = TopEvent(events)
	}

	if !ok {
		// No event — backfill an evergreen post if the gap is long enough.
		if shouldBackfill(daysSinceLastPost(outDir, time.Now())) {
			log.Printf("[marketscan] no event; backfilling evergreen")
			runAutoBlog(ctx, outDir, "auto")
			return
		}
		log.Printf("[marketscan] no event above threshold and gap < %dd — no post", backfillGapDays)
		return
	}

	log.Printf("[marketscan] top event: %s %s score=%.0f", top.Type, top.Ticker, top.Score)
	writeEventPost(ctx, outDir, top)
}

func writeEventPost(ctx context.Context, outDir string, e Event) {
	apiKey := os.Getenv("CLAUDE_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("ANTHROPIC_API_KEY")
	}
	client, err := ai.New(ai.Config{APIKey: apiKey})
	if err != nil {
		log.Fatalf("ai: %v", err)
	}
	title := eventTitle(e)
	result, err := client.AnalyzeWithContext(ctx, e.DataBlock, eventPromptBase)
	if err != nil {
		log.Fatalf("claude: %v", err)
	}
	date := time.Now().Format("2006-01-02")
	slug := eventSlug(e)
	front := fmt.Sprintf("---\ntitle: %q\ndate: %s\ntopic: %s\n---\n\n", title, date, e.Type)
	full := filepath.Join(outDir, date+"-"+slug+".md")
	if err := os.WriteFile(full, []byte(front+result.Text+"\n"), 0o644); err != nil {
		log.Fatalf("write %s: %v", full, err)
	}
	fmt.Println(full)
}

// eventTitle builds the SEO/Telegram headline (ticker + number + date).
func eventTitle(e Event) string {
	d := time.Now().Format("02/01")
	switch e.Type {
	case "mover":
		dir := "tăng"
		if e.Magnitude < 0 {
			dir = "giảm"
		}
		return fmt.Sprintf("%s %s %.1f%% hôm nay %s — cohort lịch sử nói gì?", e.Ticker, dir, math.Abs(e.Magnitude), d)
	case "volume":
		return fmt.Sprintf("%s volume gấp %.1f lần hôm nay %s — smart money hay noise?", e.Ticker, e.Magnitude, d)
	case "insider":
		return fmt.Sprintf("Nội bộ %s giao dịch lớn (%s) — đọc tín hiệu thế nào?", e.Ticker, d)
	case "flow":
		return fmt.Sprintf("Dòng tiền VN30 %s hôm nay %s — regime nào đang dẫn dắt?", flowDir(e.Magnitude), d)
	case "global":
		return fmt.Sprintf("Sốc global hôm nay %s — VN30 chịu ảnh hưởng ra sao? Cohort lookup", d)
	}
	return fmt.Sprintf("Diễn biến thị trường %s", d)
}

func flowDir(buyPressure float64) string {
	if buyPressure < 50 {
		return "nghiêng bán"
	}
	return "nghiêng mua"
}

// buildSnapshot performs all network I/O and assembles the pure Snapshot.
func buildSnapshot(ctx context.Context) Snapshot {
	universe := marketUniverse()
	quotes, _ := fetchers.VPSMultiple(ctx, universe)

	stocks := make([]StockSnap, 0, len(quotes))
	for _, q := range quotes {
		hist, err := fetchers.EntradeHistory(ctx, q.Ticker, 30)
		avg := 0.0
		rsi := 0.0
		mat := "mixed"
		if err == nil && len(hist) >= 20 {
			var vsum int64
			closes := make([]float64, len(hist))
			for i, h := range hist {
				closes[i] = h.Close
				if i >= len(hist)-20 {
					vsum += h.Volume
				}
			}
			avg = float64(vsum) / 20
			d := technical.Dashboard(closes)
			rsi = d.RSI
			if d.MA200 != nil && d.MA50 != nil {
				if q.Close >= *d.MA200 && *d.MA50 >= *d.MA200 {
					mat = "uptrend"
				} else if q.Close < *d.MA200 && *d.MA50 < *d.MA200 {
					mat = "downtrend"
				}
			}
		}
		stocks = append(stocks, StockSnap{
			Ticker: q.Ticker, Close: q.Close, ChangePct: q.ChangePercent,
			Volume: q.Volume, AvgVol20: avg, RSI: rsi, MATrend: mat,
		})
	}

	globals := []GlobalSnap{}
	for _, g := range fetchers.YahooMultiple(ctx, []string{"^VIX", "^GSPC", "^N225", "^HSI"}) {
		globals = append(globals, GlobalSnap{Symbol: g.Symbol, Price: g.Price, ChangePct: g.ChangePct})
	}

	flow := signalsFlowSafe(quotes) // existing helper in autoblog.go

	insiders := []InsiderSnap{}
	for _, st := range stocks {
		if math.Abs(st.ChangePct) < moverThresholdPct {
			continue // only check insiders for tickers already moving, to bound API calls
		}
		txs, err := fetchers.CafefInsider(ctx, st.Ticker)
		if err != nil {
			continue
		}
		for _, tx := range txs {
			if tx.BuyResult >= insiderShareThreshold {
				insiders = append(insiders, InsiderSnap{Ticker: tx.Ticker, Side: "buy", Shares: int64(tx.BuyResult), InsiderName: tx.InsiderName})
			} else if tx.SellResult >= insiderShareThreshold {
				insiders = append(insiders, InsiderSnap{Ticker: tx.Ticker, Side: "sell", Shares: int64(tx.SellResult), InsiderName: tx.InsiderName})
			}
		}
	}

	return Snapshot{
		Date: time.Now(), Stocks: stocks, Globals: globals,
		Insiders: insiders, FlowBuyPressure: flow.BuyPressure,
	}
}
```

- [ ] **Step 2: Verify build**

Run: `cd go && go build ./cmd/lmcli/`
Expected: success.

- [ ] **Step 3: Commit**

```bash
git add go/cmd/lmcli/marketscan.go
git commit -m "feat(marketscan): fetch layer + event post writer + backfill glue"
```

---

## Task 13: Wire `marketscan` into main.go

**Files:** Modify `go/cmd/lmcli/main.go`

- [ ] **Step 1: Add the dispatch case** — mirror the existing `autoblog` case (`main.go:90-103`). Add after that case:

```go
	case "marketscan":
		outDir := "docs/blog"
		for i := 0; i < len(args); i++ {
			if args[i] == "--out" && i+1 < len(args) {
				outDir = args[i+1]
				i++
			}
		}
		runMarketScan(ctx, outDir)
```

- [ ] **Step 2: Add a help line** — in the usage text (near `main.go:131`):

```go
  lmcli marketscan [--out path]  Event-triggered market blog (posts only on a real event)
```

- [ ] **Step 3: Verify build + full lmcli tests**

Run: `cd go && go build ./cmd/lmcli/ && go test ./cmd/lmcli/ -count=1`
Expected: build OK; all tests PASS.

- [ ] **Step 4: Commit**

```bash
git add go/cmd/lmcli/main.go
git commit -m "feat(marketscan): wire subcommand into CLI"
```

---

## Task 14: Switch the workflow to post-close `marketscan`

**Files:** Modify `.github/workflows/auto-blog.yml`

- [ ] **Step 1: Change the schedule** — replace the `schedule` cron (currently `0 2 * * 1,3,5`) with a weekday post-close run:

```yaml
on:
  schedule:
    - cron: "0 9 * * 1-5"   # 16:00 VN (after 15:00 close), Mon-Fri
  workflow_dispatch:
```

- [ ] **Step 2: Change the generate step** — replace the `lmcli autoblog ...` invocation (`auto-blog.yml:59`) with:

```bash
          OUT=$(/tmp/lmcli marketscan --out docs/blog)
          if [ -z "$OUT" ]; then
            echo "no-event=true" >> "${GITHUB_OUTPUT}"
            echo "::notice::marketscan produced no post today (no event, no backfill)"
            exit 0
          fi
          echo "path=${OUT}" >> "${GITHUB_OUTPUT}"
          echo "filename=$(basename "$OUT")" >> "${GITHUB_OUTPUT}"
          echo "slug=$(basename "$OUT" .md)" >> "${GITHUB_OUTPUT}"
```

- [ ] **Step 3: Guard downstream steps** — ensure the quality-gate, commit, deploy-wait, and Telegram steps only run when a post exists. Add to each of those steps' `if:` (the commit step has none yet — add one):

```yaml
        if: steps.gen.outputs.path != ''
```

(For steps that already have an `if:`, AND the existing condition with `steps.gen.outputs.path != ''`.)

- [ ] **Step 4: Validate YAML**

Run: `cd /Users/datducnguyenhuu/Desktop/lotusmarket && python3 -c "import yaml,sys; yaml.safe_load(open('.github/workflows/auto-blog.yml')); print('YAML OK')"`
Expected: `YAML OK`.

- [ ] **Step 5: Commit**

```bash
git add .github/workflows/auto-blog.yml
git commit -m "ci(autoblog): post-close marketscan schedule, skip when no post"
```

- [ ] **Step 6: Update the Pi fallback trigger (manual note)**

`trigger-workflows.sh` on the Pi fires `auto-blog.yml` only on Mon/Wed/Fri. After
this change it should fire Mon–Fri. Update the Pi script in the deploy step
(documented in the final deploy task), not in the repo — it lives on the Pi.

---

## Task 15: View counter store (`lotusai-site`)

**Files:**
- Create: `go/cmd/lotusai-site/views.go`
- Test: `go/cmd/lotusai-site/views_test.go`

- [ ] **Step 1: Write the failing test**

```go
package main

import (
	"path/filepath"
	"testing"
)

func TestIsBot(t *testing.T) {
	bots := []string{
		"Mozilla/5.0 (compatible; Googlebot/2.1)",
		"python-requests/2.31",
		"curl/8.1",
		"Go-http-client/1.1",
	}
	for _, ua := range bots {
		if !isBot(ua) {
			t.Errorf("isBot(%q) = false, want true", ua)
		}
	}
	human := "Mozilla/5.0 (Linux; Android 14; Pixel 8) Chrome/135 Mobile Safari/537.36"
	if isBot(human) {
		t.Errorf("isBot(human) = true, want false")
	}
}

func TestViewStoreIncrementAndPersist(t *testing.T) {
	path := filepath.Join(t.TempDir(), "views.json")
	vs := newViewStore(path)
	vs.hit("2026-06-23-event-mover-hpg", "Chrome/135 Mobile")
	vs.hit("2026-06-23-event-mover-hpg", "Chrome/135 Mobile")
	vs.hit("2026-06-23-event-mover-hpg", "Googlebot/2.1") // bot — ignored
	if err := vs.flush(); err != nil {
		t.Fatal(err)
	}
	reloaded := newViewStore(path)
	if got := reloaded.count("2026-06-23-event-mover-hpg"); got != 2 {
		t.Fatalf("count = %d, want 2 (bot excluded, persisted)", got)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd go && go test ./cmd/lotusai-site/ -run 'TestIsBot|TestViewStore'`
Expected: FAIL — `undefined: isBot` / `newViewStore`.

- [ ] **Step 3: Write minimal implementation**

```go
package main

import (
	"encoding/json"
	"os"
	"strings"
	"sync"
)

var botMarkers = []string{
	"bot", "crawl", "spider", "slurp", "bing", "google", "yandex",
	"facebookexternal", "semrush", "ahrefs", "petal", "gpt", "claude",
	"python-requests", "curl", "go-http", "libredtail", "wget", "scan",
}

func isBot(ua string) bool {
	l := strings.ToLower(ua)
	if l == "" {
		return true
	}
	for _, m := range botMarkers {
		if strings.Contains(l, m) {
			return true
		}
	}
	return false
}

type viewStore struct {
	path string
	mu   sync.Mutex
	data map[string]int
}

func newViewStore(path string) *viewStore {
	vs := &viewStore{path: path, data: map[string]int{}}
	if b, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(b, &vs.data)
	}
	return vs
}

func (vs *viewStore) hit(slug, ua string) {
	if isBot(ua) {
		return
	}
	vs.mu.Lock()
	vs.data[slug]++
	vs.mu.Unlock()
}

func (vs *viewStore) count(slug string) int {
	vs.mu.Lock()
	defer vs.mu.Unlock()
	return vs.data[slug]
}

func (vs *viewStore) flush() error {
	vs.mu.Lock()
	b, err := json.MarshalIndent(vs.data, "", "  ")
	vs.mu.Unlock()
	if err != nil {
		return err
	}
	return os.WriteFile(vs.path, b, 0o644)
}

func (vs *viewStore) snapshot() map[string]int {
	vs.mu.Lock()
	defer vs.mu.Unlock()
	out := make(map[string]int, len(vs.data))
	for k, v := range vs.data {
		out[k] = v
	}
	return out
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd go && go test ./cmd/lotusai-site/ -run 'TestIsBot|TestViewStore'`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add go/cmd/lotusai-site/views.go go/cmd/lotusai-site/views_test.go
git commit -m "feat(site): bot-filtered view-counter store"
```

---

## Task 16: Wire view counter into the site

**Files:**
- Modify: `go/cmd/lotusai-site/main.go` (add `--data` flag)
- Modify: `go/cmd/lotusai-site/server.go` (hold `viewStore`, count in `handleBlog`, add `/blog/stats`, periodic flush)

- [ ] **Step 1: Add the `--data` flag** in `main.go` (near `main.go:35-38`):

```go
	dataDir = flag.String("data", "./data", "directory for runtime data (view counts)")
```

And pass it into the server config (extend the `serverConfig{...}` literal in
`main.go` and the `serverConfig` struct in `server.go` with `DataDir string`).

- [ ] **Step 2: Initialize the store** in `newServer` (`server.go:~38`), after the
struct is built:

```go
	if err := os.MkdirAll(cfg.DataDir, 0o755); err == nil {
		s.views = newViewStore(filepath.Join(cfg.DataDir, "blog_views.json"))
		go func() {
			t := time.NewTicker(2 * time.Minute)
			for range t.C {
				_ = s.views.flush()
			}
		}()
	}
```

Add `views *viewStore` to the `server` struct (`server.go:30`). Ensure `os`,
`filepath`, `time` are imported (they already are in `server.go`).

- [ ] **Step 3: Count on post detail** — in `handleBlog` (`handlers.go:136`, the
branch that renders a matched post), before `s.renderPage(...)`:

```go
	if s.views != nil {
		s.views.hit(match.Slug, r.UserAgent())
	}
```

- [ ] **Step 4: Add `/blog/stats`** — register in `newServer` mux (`server.go:~57`):

```go
	mux.HandleFunc("/blog/stats", s.handleBlogStats)
```

And add the handler in `server.go`:

```go
func (s *server) handleBlogStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if s.views == nil {
		w.Write([]byte("{}"))
		return
	}
	b, _ := json.MarshalIndent(s.views.snapshot(), "", "  ")
	w.Write(b)
}
```

Ensure `encoding/json` and `net/http` are imported in `server.go`.

Note: `/blog/stats` must be registered before `/blog/` so the more specific
pattern wins — Go's `ServeMux` longest-match handles this regardless of order, but
verify `/blog/stats` does not get treated as a post slug by `handleBlog`. Since it
is its own registered handler, it takes precedence.

- [ ] **Step 5: Verify build + tests**

Run: `cd go && go build ./cmd/lotusai-site/ && go test ./cmd/lotusai-site/ -count=1`
Expected: build OK; tests PASS.

- [ ] **Step 6: Commit**

```bash
git add go/cmd/lotusai-site/main.go go/cmd/lotusai-site/server.go go/cmd/lotusai-site/handlers.go
git commit -m "feat(site): count non-bot blog reads, expose /blog/stats"
```

---

## Task 17: Full verification + smoke test

- [ ] **Step 1: Full module build + vet + test**

Run: `cd go && go build ./... && go vet ./... && go test ./... -count=1`
Expected: all green.

- [ ] **Step 2: Smoke-run marketscan against live data** (requires `CLAUDE_API_KEY`)

Run:
```bash
cd go && go build -o /tmp/lmcli ./cmd/lmcli
mkdir -p /tmp/blogtest
CLAUDE_API_KEY=$KEY /tmp/lmcli marketscan --out /tmp/blogtest || true
ls -la /tmp/blogtest
```
Expected: either a single `YYYY-MM-DD-event-*.md` (or backfill `*.md`) printed, OR
log "no event above threshold" and no file. Inspect the generated file: title
contains a ticker/number/date; body follows the 4-part structure; no fabricated
numbers beyond the DataBlock.

- [ ] **Step 3: Commit any fixups**, then deploy the site (next task).

---

## Task 18: Deploy the site (view counter) to the Pi

The `marketscan` generator runs in GitHub Actions (no deploy). Only the view
counter needs a Pi deploy. Use the deploy-pi pattern.

- [ ] **Step 1: Confirm the systemd unit's start command** so the new `--data`
flag points at a persistent dir:

Run: `ssh pi 'cat /etc/systemd/system/lotusai-site.service'`
Expected: note the `ExecStart` line and its `--content` path (e.g.
`/home/pi/lotusai/content`). The new `--data` should be `/home/pi/lotusai/data`.

- [ ] **Step 2: If needed, add `--data /home/pi/lotusai/data` to `ExecStart`** and
reload:

```bash
ssh pi 'sudo systemctl daemon-reload'
```

- [ ] **Step 3: Build ARM64 + deploy** (stop before scp to avoid "text file busy"):

```bash
cd go && GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o /tmp/lotusai-site ./cmd/lotusai-site
ssh pi 'sudo systemctl stop lotusai-site'
scp /tmp/lotusai-site pi:/home/pi/lotusai/lotusai-site
ssh pi 'mkdir -p /home/pi/lotusai/data && sudo systemctl start lotusai-site && sleep 3 && sudo systemctl status lotusai-site --no-pager | head -6'
```

- [ ] **Step 4: Verify**

```bash
curl -s -o /dev/null -w "%{http_code}\n" https://lotusai.servehttp.com/blog/
curl -s https://lotusai.servehttp.com/blog/stats | head
```
Expected: `200`; stats returns JSON (likely `{}` until reads accrue).

- [ ] **Step 5: Update the Pi fallback trigger** for the new Mon–Fri schedule:

Edit `/home/pi/lotusai/trigger-workflows.sh` so the auto-blog block fires on
`dow` 1–5 (not just 1,3,5). Keep the 18h staleness guard.

---

## Self-Review Notes

- **Spec coverage:** event taxonomy (Tasks 2–6), scoring/threshold/tie-break (7),
  dedup (8), backfill enabled gap=3 (9), VN30+HNX30 universe (10), engaging
  event prompt (11), fetch glue + title formula (12–13), post-close cron (14),
  bot-filtered view counter + `/blog/stats` (15–16), Pi deploy (18). Phase 2
  (CafeF RSS) is intentionally deferred to a separate plan.
- **Tunables** live as named consts: `moverThresholdPct`, `volumeSurgeMult`,
  `insiderShareThreshold`, `eventThreshold`, `eventCooldownDays`, `backfillGapDays`.
- **Reuse:** `matchCohort`, `safeSlug`, `signalsFlowSafe`, `recentContentKeys`,
  `recentVariantKeys`, `runAutoBlog` are all existing in `autoblog.go`.
- **Open data item:** the `hnx30` constituent list must be verified against the
  current official HNX30 index before merge (flagged in Task 10).
