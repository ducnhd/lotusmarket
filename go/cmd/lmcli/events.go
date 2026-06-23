package main

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"
)

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
	Ticker      string
	Side        string // "buy" | "sell"
	Shares      int64
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
	out = append(out, detectMovers(s)...)
	out = append(out, detectVolume(s)...)
	out = append(out, detectFlow(s)...)
	out = append(out, detectGlobal(s)...)
	out = append(out, detectInsider(s)...)
	return out
}

// ============================================================================
// Per-type detectors
// ============================================================================

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

func detectFlow(s Snapshot) []Event {
	if s.FlowBuyPressure <= 0 { // zero-value means not populated
		return nil
	}
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

// ============================================================================
// Ranking + threshold
// ============================================================================

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

// ============================================================================
// Dedup
// ============================================================================

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

// ============================================================================
// Backfill helpers
// ============================================================================

const backfillGapDays = 3

func shouldBackfill(daysSince int) bool { return daysSince >= backfillGapDays }

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
