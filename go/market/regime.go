package market

import (
	"fmt"
	"strings"
)

// Regime constants.
const (
	RegimeStable   = "STABLE"
	RegimeVolatile = "VOLATILE"
	RegimeCrisis   = "CRISIS"
	RegimeEuphoria = "EUPHORIA"

	// Crisis sub-types.
	CrisisPanic       = "CRISIS_PANIC"       // contagion/external shock → buying opportunity
	CrisisFundamental = "CRISIS_FUNDAMENTAL" // VN structural problem → avoid
)

// RegimeSignals holds individual signal scores used for classification.
// All fields are public so callers can log/inspect intermediate state.
type RegimeSignals struct {
	VIX            float64 `json:"vix"`
	VIXScore       int     `json:"vix_score"`
	VNIndexChange  float64 `json:"vnindex_change_pct"`
	VNIndexScore   int     `json:"vnindex_score"`
	ForeignFlow    float64 `json:"foreign_flow_b"` // billions VND
	FlowScore      int     `json:"flow_score"`
	NewsTier       int     `json:"news_tier"` // 1, 2, 3, or 0
	NewsScore      int     `json:"news_score"`
	NewsDetail     string  `json:"news_detail"`
	GlobalDrop     float64 `json:"global_max_drop_pct"`
	GlobalScore    int     `json:"global_score"`
	CompositeScore int     `json:"composite_score"`
	Direction      string  `json:"direction"` // "down" | "up" | "mixed"
}

// ScoreSignals converts raw inputs into a RegimeSignals struct with all scores computed.
// This is the pure-logic equivalent of collectSignals in services/market_regime.go.
//
//   - vix: current VIX level
//   - vnIndexChange: VN-Index daily change percent (negative = down)
//   - foreignFlowBillions: net foreign flow in billions VND (negative = net sell)
//   - newsTier: highest active news tier (1 = arrest/circuit breaker, 2 = Fed/tariffs, 3 = GDP/CPI, 0 = none)
//   - newsDetail: first matching TIER 1/2 headline (used for crisis sub-type detection)
//   - globalMaxDropPct: most negative single-day drop among major global indices (≤0)
func ScoreSignals(vix, vnIndexChange, foreignFlowBillions float64, newsTier int, newsDetail string, globalMaxDropPct float64) RegimeSignals {
	var s RegimeSignals
	s.VIX = vix
	s.VNIndexChange = vnIndexChange
	s.ForeignFlow = foreignFlowBillions
	s.NewsTier = newsTier
	s.NewsDetail = newsDetail
	s.GlobalDrop = globalMaxDropPct

	// VIX score
	switch {
	case vix >= 30:
		s.VIXScore = 80
	case vix >= 25:
		s.VIXScore = 50
	case vix >= 20:
		s.VIXScore = 25
	default:
		s.VIXScore = 0
	}

	// VN-Index score (uses absolute change)
	absChange := vnIndexChange
	if absChange < 0 {
		absChange = -absChange
	}
	switch {
	case absChange >= 3.0:
		s.VNIndexScore = 80
	case absChange >= 2.0:
		s.VNIndexScore = 50
	case absChange >= 1.0:
		s.VNIndexScore = 25
	default:
		s.VNIndexScore = 0
	}

	// Direction from VN-Index change
	switch {
	case vnIndexChange < -1.0:
		s.Direction = "down"
	case vnIndexChange > 1.0:
		s.Direction = "up"
	default:
		s.Direction = "mixed"
	}

	// Foreign flow score
	switch {
	case foreignFlowBillions < -2000:
		s.FlowScore = 80
	case foreignFlowBillions < -1000:
		s.FlowScore = 50
	case foreignFlowBillions < -500:
		s.FlowScore = 25
	default:
		s.FlowScore = 0
	}

	// News score
	switch newsTier {
	case 1:
		s.NewsScore = 100
	case 2:
		s.NewsScore = 50
	default:
		s.NewsScore = 0
	}

	// Global contagion score
	switch {
	case globalMaxDropPct <= -5.0:
		s.GlobalScore = 80
	case globalMaxDropPct <= -3.0:
		s.GlobalScore = 50
	default:
		s.GlobalScore = 0
	}

	// Composite = max of individual scores (one strong signal is enough)
	s.CompositeScore = max(s.VIXScore, s.VNIndexScore, s.FlowScore, s.NewsScore, s.GlobalScore)

	return s
}

// fundamentalKeywords are VN-specific structural risk terms.
var fundamentalKeywords = []string{
	"bắt giữ", "khởi tố", "truy tố", "lừa đảo", "rút ruột",
	"bank run", "vỡ nợ", "phá sản", "hủy niêm yết",
	"kiểm soát đặc biệt", "đình chỉ",
}

// classifyCrisisType determines if a crisis is panic (opportunity) or fundamental (avoid).
func classifyCrisisType(s RegimeSignals) string {
	newsLower := strings.ToLower(s.NewsDetail)
	for _, kw := range fundamentalKeywords {
		if strings.Contains(newsLower, kw) {
			return CrisisFundamental
		}
	}
	if s.GlobalScore >= 80 {
		return CrisisPanic
	}
	if s.FlowScore >= 80 && s.NewsScore < 100 {
		return CrisisPanic
	}
	return CrisisPanic
}

// Classify determines the market regime and sub-type from pre-computed signals.
// Returns (regime, subType). subType is empty for non-crisis regimes.
func Classify(s RegimeSignals) (regime, subType string) {
	if s.CompositeScore >= 80 {
		if s.Direction == "up" {
			return RegimeEuphoria, ""
		}
		return RegimeCrisis, classifyCrisisType(s)
	}
	if s.CompositeScore >= 50 {
		return RegimeVolatile, ""
	}
	return RegimeStable, ""
}

// IdentifyTrigger returns the primary trigger (source, human-readable detail) for a regime.
func IdentifyTrigger(s RegimeSignals) (source, detail string) {
	maxScore := 0
	source = "composite"
	detail = ""

	if s.NewsScore > maxScore {
		maxScore = s.NewsScore
		source = "news_tier1"
		detail = s.NewsDetail
	}
	if s.VNIndexScore > maxScore {
		maxScore = s.VNIndexScore
		source = "vnindex_change"
		detail = fmt.Sprintf("VN-Index thay đổi %.1f%%", s.VNIndexChange)
	}
	if s.GlobalScore > maxScore {
		maxScore = s.GlobalScore
		source = "global_contagion"
		detail = fmt.Sprintf("Global index giảm %.1f%%", s.GlobalDrop)
	}
	if s.VIXScore > maxScore {
		maxScore = s.VIXScore
		source = "vix"
		detail = fmt.Sprintf("VIX = %.1f", s.VIX)
	}
	if s.FlowScore > maxScore {
		source = "foreign_flow"
		detail = fmt.Sprintf("Foreign net flow: %.0fB VND", s.ForeignFlow)
	}

	return source, detail
}
