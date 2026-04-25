// Package exposure provides per-ticker mapping of external drivers (commodity /
// FX / peer stock) with correlations validated offline over 2 years of daily
// closes. Each entry carries historical Pearson r at the optimal horizon so
// callers can compute principled "estimated impact" numbers.
package exposure

import (
	"fmt"
	"sort"
	"strings"
)

// ExposureSide groups drivers for the "Input → Stock → Output" UI layout.
type ExposureSide string

const (
	SideInput   ExposureSide = "input"   // raw materials / import costs
	SidePeer    ExposureSide = "peer"    // competitors moving in sympathy
	SideOutput  ExposureSide = "output"  // end-demand signals / sector trends
	SideContext ExposureSide = "context" // market-wide beta, risk regime
)

// ExposureMapping is a single validated driver for a ticker.
type ExposureMapping struct {
	Driver       string       `json:"driver"`        // Yahoo symbol, e.g. "TIO=F"
	LabelVi      string       `json:"label_vi"`      // Vietnamese UI label
	Category     string       `json:"category"`      // INPUT, PEER, DEMAND, MARKET, RISK, FLOW, SECTOR, COST
	Side         ExposureSide `json:"side"`          // input / peer / output / context
	Sign         int          `json:"sign"`          // +1 positive r expected, -1 negative
	HistoricalR  float64      `json:"historical_r"`  // signed Pearson r from backtest
	HorizonDays  int          `json:"horizon_days"`  // return window that gave max r
	SampleN      int          `json:"sample_n"`      // observations in backtest
	Explanation  string       `json:"explanation"`   // plain-Vietnamese reason
	ConceptGloss string       `json:"concept_gloss"` // beginner-friendly definition
}

// DriverClose is a single daily closing price for a driver symbol.
type DriverClose struct {
	Date  string // YYYY-MM-DD
	Close float64
}

// HistoryProvider fetches daily closes for a driver symbol.
// symbol is the Yahoo Finance symbol; days is the lookback window.
// Returned slice must be ordered by date ascending (oldest first).
type HistoryProvider func(symbol string, days int) []DriverClose

// DriverImpact is the per-driver payload returned by Analyze.
type DriverImpact struct {
	ExposureMapping
	LatestClose     *float64 `json:"latest_close,omitempty"`
	PreviousClose   *float64 `json:"previous_close,omitempty"`
	ReturnPct       *float64 `json:"return_pct,omitempty"`     // % change over HorizonDays
	EstImpactPct    *float64 `json:"est_impact_pct,omitempty"` // return × r
	ImpactMagnitude string   `json:"impact_magnitude"`         // small / medium / large / none
	Direction       string   `json:"direction"`                // up / down / flat
	LastUpdated     string   `json:"last_updated,omitempty"`   // YYYY-MM-DD of latest close
}

// ExposureAnalysis is the top-level result for one ticker.
type ExposureAnalysis struct {
	Ticker     string         `json:"ticker"`
	Input      []DriverImpact `json:"input"`
	Peer       []DriverImpact `json:"peer"`
	Output     []DriverImpact `json:"output"`
	Context    []DriverImpact `json:"context"`
	NetImpact  float64        `json:"net_impact_pct"` // weighted average of est_impact
	Summary    string         `json:"summary_vi"`
	HasMapping bool           `json:"has_mapping"`
}

// exposureDefs is the canonical mapping derived from cmd/exposure_backtest.
var exposureDefs = buildExposureDefs()

func buildExposureDefs() map[string][]ExposureMapping {
	m := map[string][]ExposureMapping{}

	// ─── Thép (HPG) ─────────────────────────────────────────────────────────
	m["HPG"] = []ExposureMapping{
		{"000300.SS", "Chứng khoán Trung Quốc (CSI 300)", "DEMAND", SideOutput, +1, 0.461, 20, 456,
			"Kinh tế TQ mạnh → bất động sản & hạ tầng TQ tiêu thụ nhiều thép → cầu thép toàn cầu tăng → HPG hưởng lợi.",
			"CSI 300 là chỉ số 300 cổ phiếu lớn nhất Trung Quốc — dùng làm 'nhiệt kế' cho kinh tế TQ."},
		{"BHP", "Giá cổ phiếu BHP (công ty khai thác quặng)", "INPUT", SideInput, +1, 0.441, 20, 462,
			"BHP là công ty khai thác quặng sắt lớn nhất thế giới. Giá BHP phản ánh giá quặng sắt — nguyên liệu chính của HPG.",
			"'Quặng sắt' là đá chứa sắt, khi nấu trong lò cao sẽ ra thép. Chiếm ~30% chi phí sản xuất thép."},
		{"600019.SS", "Baoshan Iron & Steel (đối thủ TQ)", "PEER", SidePeer, +1, 0.420, 20, 456,
			"Baoshan là hãng thép lớn nhất TQ. Khi giá Baoshan tăng, thường cả ngành thép Châu Á đang lên theo cùng chu kỳ.",
			"'Đối thủ cùng ngành' là công ty làm cùng sản phẩm — thường giá cổ phiếu di chuyển cùng chiều theo chu kỳ ngành."},
		{"TIO=F", "Hợp đồng tương lai quặng sắt (SGX)", "INPUT", SideInput, +1, 0.278, 20, 463,
			"Giá quặng sắt giao dịch trên sàn Singapore — thước đo giá nguyên liệu đầu vào của HPG.",
			"'Hợp đồng tương lai' (futures) là thoả thuận mua/bán hàng hoá ở mức giá xác định trong tương lai."},
		{"CNY=X", "Tỷ giá USD/CNY", "DEMAND", SideOutput, -1, -0.216, 5, 431,
			"CNY yếu đi (USD/CNY tăng) thường đi kèm với việc kinh tế TQ chậm lại → cầu thép yếu → HPG giảm.",
			"Khi CNY yếu, hàng xuất khẩu TQ rẻ hơn nhưng đồng thời báo hiệu kinh tế nội địa TQ đang gặp khó."},
	}

	// ─── Ngân hàng — 14 mã VN30 ─────────────────────────────────────────────
	bankPack := []ExposureMapping{
		{"VNM", "Quỹ ETF VanEck Vietnam (VNM)", "MARKET", SideContext, +1, 0.70, 20, 460,
			"Là quỹ ETF tại Mỹ nắm giữ các cổ phiếu VN. Phản ánh xu hướng chung của toàn thị trường VN.",
			"ETF là quỹ mô phỏng một rổ cổ phiếu — tiền nước ngoài đầu tư vào VN thường đi qua ETF này."},
		{"^VIX", "Chỉ số sợ hãi (VIX)", "RISK", SideContext, -1, -0.30, 20, 455,
			"VIX tăng cao → nhà đầu tư toàn cầu rút tiền khỏi thị trường mới nổi (EM) → ngân hàng VN thường bị bán mạnh.",
			"VIX (thường gọi là 'chỉ số sợ hãi') đo mức biến động kỳ vọng của S&P 500. Cao = thị trường lo lắng."},
		{"EEM", "Quỹ ETF thị trường mới nổi (EEM)", "FLOW", SideContext, +1, 0.37, 20, 462,
			"EEM phản ánh dòng tiền quốc tế vào các thị trường mới nổi (VN, TQ, Ấn Độ, Brazil...). Dòng tiền vào EM = tốt cho ngân hàng VN.",
			"'Thị trường mới nổi' (Emerging Markets) là các nước đang phát triển — vốn nhạy với rủi ro toàn cầu."},
	}
	for _, t := range []string{"VCB", "TCB", "BID", "CTG", "ACB", "MBB", "HDB", "LPB", "SHB", "SSB", "STB", "TPB", "VIB", "VPB"} {
		m[t] = bankPack
	}

	// ─── Chứng khoán — SSI ─────────────────────────────────────────────────
	m["SSI"] = []ExposureMapping{
		{"VNM", "Quỹ ETF VanEck Vietnam (VNM)", "MARKET", SideContext, +1, 0.747, 20, 460,
			"Thanh khoản thị trường là nguồn thu chính của công ty chứng khoán. VN-Index lên = giao dịch sôi động = SSI thu phí nhiều.",
			"'Thanh khoản' = lượng giao dịch mua bán trong ngày. Cao = người ta giao dịch nhiều."},
		{"EEM", "Quỹ ETF thị trường mới nổi (EEM)", "FLOW", SideContext, +1, 0.415, 20, 455,
			"Khi dòng tiền quốc tế đổ vào EM, khối lượng giao dịch tại VN tăng → SSI hưởng lợi.", ""},
		{"^VIX", "Chỉ số sợ hãi (VIX)", "RISK", SideContext, -1, -0.243, 20, 456,
			"VIX cao → nhà đầu tư cá nhân VN ngại giao dịch → doanh thu môi giới SSI giảm.", ""},
	}

	// ─── Công nghệ — FPT ───────────────────────────────────────────────────
	m["FPT"] = []ExposureMapping{
		{"QQQ", "Nasdaq 100 ETF (QQQ)", "SECTOR", SideOutput, +1, 0.569, 20, 462,
			"FPT bán dịch vụ IT cho Mỹ, Nhật. Khi ngành công nghệ toàn cầu (Nasdaq) bùng nổ, khách hàng chi nhiều hơn → FPT tăng đơn hàng.",
			"QQQ là quỹ ETF theo 100 cổ phiếu công nghệ lớn nhất Mỹ (Apple, Microsoft, Nvidia...). Đo 'sức khoẻ' ngành công nghệ."},
		{"XLK", "US Tech Sector ETF (XLK)", "SECTOR", SideOutput, +1, 0.522, 20, 462,
			"Tương tự QQQ, đo xu hướng chi tiêu cho công nghệ.", ""},
		{"VNM", "Quỹ ETF VanEck Vietnam (VNM)", "MARKET", SideContext, +1, 0.377, 1, 481,
			"Beta thị trường — FPT là mã trụ trong VN-Index.", ""},
	}

	// ─── Bán lẻ — MWG ──────────────────────────────────────────────────────
	m["MWG"] = []ExposureMapping{
		{"VNM", "Quỹ ETF VanEck Vietnam (VNM)", "MARKET", SideContext, +1, 0.635, 20, 459,
			"Beta thị trường chung.", ""},
		{"EEM", "ETF thị trường mới nổi (EEM)", "FLOW", SideContext, +1, 0.565, 20, 459,
			"Tâm lý tiêu dùng EM tốt = MWG hưởng lợi.", ""},
		{"DX-Y.NYB", "Chỉ số USD (DXY)", "COST", SideInput, -1, -0.404, 20, 457,
			"MWG nhập iPhone, điện thoại, đồ điện tử — trả bằng USD. USD mạnh lên → giá nhập đắt → biên lợi nhuận giảm.",
			"DXY đo sức mạnh USD so với 6 đồng tiền lớn (EUR, JPY, GBP...). DXY tăng = USD mạnh lên toàn cầu."},
	}

	// ─── Bất động sản — Vingroup ───────────────────────────────────────────
	m["VHM"] = []ExposureMapping{
		{"VNM", "Quỹ ETF VanEck Vietnam (VNM)", "MARKET", SideContext, +1, 0.682, 20, 459,
			"Beta thị trường.", ""},
		{"FXI", "ETF cổ phiếu lớn Trung Quốc (FXI)", "DEMAND", SideOutput, +1, 0.156, 20, 459,
			"Tâm lý về BĐS TQ lan toả sang BĐS VN — khi BĐS TQ hồi phục, NĐT cũng tin vào BĐS VN hơn.", ""},
	}
	m["VIC"] = []ExposureMapping{
		{"VNM", "Quỹ ETF VanEck Vietnam (VNM)", "MARKET", SideContext, +1, 0.551, 5, 477,
			"Beta thị trường.", ""},
	}
	m["VRE"] = []ExposureMapping{
		{"VNM", "Quỹ ETF VanEck Vietnam (VNM)", "MARKET", SideContext, +1, 0.672, 20, 459,
			"Beta thị trường.", ""},
	}

	// ─── Năng lượng — GAS, PLX ─────────────────────────────────────────────
	m["GAS"] = []ExposureMapping{
		{"VNM", "Quỹ ETF VanEck Vietnam (VNM)", "MARKET", SideContext, +1, 0.40, 20, 460,
			"Beta thị trường. Lưu ý: giá bán khí nội địa có cơ chế riêng, ít bị ảnh hưởng trực tiếp bởi dầu thế giới.", ""},
	}
	m["PLX"] = []ExposureMapping{
		{"VNM", "Quỹ ETF VanEck Vietnam (VNM)", "MARKET", SideContext, +1, 0.40, 20, 460,
			"Beta thị trường. Giá xăng dầu bán lẻ VN do Bộ Công Thương quyết định theo chu kỳ 10 ngày, không biến động theo giờ như dầu WTI.", ""},
	}

	// ─── Tiêu dùng ─────────────────────────────────────────────────────────
	m["MSN"] = []ExposureMapping{
		{"VNM", "Quỹ ETF VanEck Vietnam (VNM)", "MARKET", SideContext, +1, 0.636, 20, 459,
			"Beta thị trường.", ""},
		{"FXI", "ETF cổ phiếu lớn TQ (FXI)", "DEMAND", SideOutput, +1, 0.448, 20, 459,
			"Masan là mảng tiêu dùng lớn; tâm lý tiêu dùng khu vực Châu Á (đại diện bởi TQ) ảnh hưởng đến MSN.", ""},
	}
	m["SAB"] = []ExposureMapping{
		{"VNM", "Quỹ ETF VanEck Vietnam (VNM)", "MARKET", SideContext, +1, 0.427, 5, 474,
			"Beta thị trường. SAB ít phụ thuộc commodity — bia là sản phẩm 'defensive'.", ""},
	}
	// Note: key "VNM" = Vinamilk VN stock code; driver "VNM" = VanEck ETF Yahoo symbol.
	m["VNM"] = []ExposureMapping{
		{"VNM", "Quỹ ETF VanEck Vietnam (VNM ETF)", "MARKET", SideContext, +1, 0.469, 5, 474,
			"Beta thị trường chung của Vinamilk.",
			"Trùng tên viết tắt — mã VNM là Vinamilk (VN), còn 'VNM ETF' là quỹ Mỹ đầu tư vào VN."},
		{"DX-Y.NYB", "Chỉ số USD (DXY)", "COST", SideInput, -1, -0.291, 20, 462,
			"Vinamilk nhập sữa bột (New Zealand, EU) trả bằng USD. USD mạnh = nguyên liệu đắt = biên lợi nhuận co lại.", ""},
	}

	// ─── Hoá chất — DGC ────────────────────────────────────────────────────
	m["DGC"] = []ExposureMapping{
		{"MOS", "Mosaic (công ty phân bón Mỹ)", "PEER", SidePeer, +1, 0.267, 5, 474,
			"Mosaic là hãng sản xuất phân bón phosphate lớn nhất thế giới. DGC cùng sản xuất phốt pho/phân bón nên di chuyển theo chu kỳ ngành hoá chất.",
			"'Phốt pho vàng' (yellow phosphorus) là nguyên liệu cho phân bón, pin lithium, dệt may — sản phẩm chủ lực của DGC."},
		{"NTR", "Nutrien (công ty phân bón Canada)", "PEER", SidePeer, +1, 0.229, 5, 473,
			"Nutrien là đối thủ toàn cầu trong mảng phân bón/hoá chất nông nghiệp.", ""},
		{"CNY=X", "Tỷ giá USD/CNY", "DEMAND", SideOutput, -1, -0.159, 5, 428,
			"TQ là thị trường xuất khẩu lớn của DGC. CNY yếu = TQ mua ít hơn.", ""},
	}

	// ─── Nông nghiệp — GVR (cao su) ────────────────────────────────────────
	m["GVR"] = []ExposureMapping{
		{"5108.T", "Bridgestone (hãng lốp Nhật)", "DEMAND", SideOutput, +1, 0.317, 20, 446,
			"GVR bán cao su thô. Khách hàng lớn là các hãng lốp xe. Bridgestone tăng = cầu lốp tốt = cầu cao su tốt.",
			"Cao su thiên nhiên chủ yếu dùng sản xuất lốp ô tô — chiếm 70% nhu cầu cao su toàn cầu."},
		{"GT", "Goodyear (hãng lốp Mỹ)", "DEMAND", SideOutput, +1, 0.188, 20, 454,
			"Tương tự Bridgestone — Goodyear là proxy cho cầu lốp xe thế giới.", ""},
	}

	// ─── Hàng không — VJC ──────────────────────────────────────────────────
	m["VJC"] = []ExposureMapping{
		{"VNM", "Quỹ ETF VanEck Vietnam (VNM)", "MARKET", SideContext, +1, 0.520, 20, 459,
			"Beta thị trường.", ""},
		{"CL=F", "Dầu WTI", "COST", SideInput, -1, -0.254, 20, 462,
			"Nhiên liệu máy bay chiếm 30-40% chi phí vận hành hãng hàng không. Dầu tăng = chi phí tăng = VJC giảm.",
			"WTI (West Texas Intermediate) là loại dầu tiêu chuẩn của Mỹ — benchmark giá dầu toàn cầu."},
		{"HO=F", "Dầu đốt (proxy cho nhiên liệu bay)", "COST", SideInput, -1, -0.251, 20, 462,
			"Heating oil là dẫn xuất gần với jet fuel (dầu máy bay) hơn WTI.", ""},
	}

	return m
}

// UniqueDriverSymbols returns every Yahoo symbol referenced by at least one mapping (sorted).
func UniqueDriverSymbols() []string {
	seen := map[string]bool{}
	for _, list := range exposureDefs {
		for _, e := range list {
			seen[e.Driver] = true
		}
	}
	out := make([]string, 0, len(seen))
	for s := range seen {
		out = append(out, s)
	}
	sort.Strings(out)
	return out
}

// GetMappingsForTicker returns mappings for a VN ticker, or nil if unmapped.
func GetMappingsForTicker(ticker string) []ExposureMapping {
	return exposureDefs[ticker]
}

func absFloat(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}

func joinStrings(ss []string, sep string) string {
	return strings.Join(ss, sep)
}

// Analyze computes driver impacts for a ticker using the provided history.
// historyProvider is called with (symbol, days) — it must return closes ascending by date.
func Analyze(ticker string, hp HistoryProvider) ExposureAnalysis {
	mappings := GetMappingsForTicker(ticker)
	if len(mappings) == 0 {
		return ExposureAnalysis{
			Ticker:     ticker,
			HasMapping: false,
			Summary:    "Chưa có dữ liệu tương quan cho mã này. Chỉ các mã đã được backtest (HPG, ngân hàng VN30, FPT, MWG, DGC, GVR, VJC, VNM, MSN...) mới có phân tích chuỗi tác động.",
		}
	}

	result := ExposureAnalysis{Ticker: ticker, HasMapping: true}
	var netWeighted, totalWeight float64

	for _, m := range mappings {
		impact := DriverImpact{ExposureMapping: m, ImpactMagnitude: "none", Direction: "flat"}

		// Request generous window so weekends/holidays don't starve the lookup.
		history := hp(m.Driver, m.HorizonDays*3+10)
		if len(history) >= 2 {
			latest := history[len(history)-1]
			prevIdx := max(len(history)-1-m.HorizonDays, 0)
			prev := history[prevIdx]
			if prev.Close > 0 {
				lc, pc := latest.Close, prev.Close
				ret := (lc - pc) / pc * 100
				impact.LatestClose = &lc
				impact.PreviousClose = &pc
				impact.ReturnPct = &ret
				impact.LastUpdated = latest.Date

				est := ret * m.HistoricalR
				impact.EstImpactPct = &est

				switch {
				case est > 0.3:
					impact.Direction = "up"
				case est < -0.3:
					impact.Direction = "down"
				}
				switch mag := absFloat(est); {
				case mag > 1.5:
					impact.ImpactMagnitude = "large"
				case mag > 0.5:
					impact.ImpactMagnitude = "medium"
				case mag > 0.1:
					impact.ImpactMagnitude = "small"
				}

				w := absFloat(m.HistoricalR)
				netWeighted += est * w
				totalWeight += w
			}
		}

		switch m.Side {
		case SideInput:
			result.Input = append(result.Input, impact)
		case SidePeer:
			result.Peer = append(result.Peer, impact)
		case SideOutput:
			result.Output = append(result.Output, impact)
		case SideContext:
			result.Context = append(result.Context, impact)
		}
	}

	if totalWeight > 0 {
		result.NetImpact = netWeighted / totalWeight
	}
	result.Summary = buildSummary(ticker, &result)
	return result
}

func buildSummary(ticker string, a *ExposureAnalysis) string {
	all := append(append(append([]DriverImpact{}, a.Input...), a.Peer...), a.Output...)
	all = append(all, a.Context...)

	var top *DriverImpact
	for i := range all {
		if all[i].EstImpactPct == nil {
			continue
		}
		if top == nil || absFloat(*all[i].EstImpactPct) > absFloat(*top.EstImpactPct) {
			top = &all[i]
		}
	}
	if top == nil || top.EstImpactPct == nil {
		return fmt.Sprintf("Chưa có dữ liệu gần nhất cho các yếu tố ảnh hưởng đến %s.", ticker)
	}
	dir := "tăng"
	if a.NetImpact < 0 {
		dir = "giảm"
	}
	return fmt.Sprintf("Dựa trên biến động %d ngày của các yếu tố liên quan, %s có áp lực %s khoảng %.2f%%. Yếu tố mạnh nhất: %s (%+.1f%% → tác động ước tính %+.2f%%).",
		top.HorizonDays, ticker, dir, a.NetImpact, top.LabelVi, *top.ReturnPct, *top.EstImpactPct)
}

// FormatForPrompt returns a compact multi-line text block for use in Claude prompts.
// Mirrors the logic of services.StockExposureService.FormatForPrompt.
func FormatForPrompt(tickers []string, hp HistoryProvider) string {
	if len(tickers) == 0 {
		return ""
	}
	var sb []string
	for _, t := range tickers {
		a := Analyze(t, hp)
		if !a.HasMapping {
			continue
		}
		all := append(append(append([]DriverImpact{}, a.Input...), a.Peer...), a.Output...)
		all = append(all, a.Context...)
		sort.Slice(all, func(i, j int) bool {
			ai, aj := 0.0, 0.0
			if all[i].EstImpactPct != nil {
				ai = absFloat(*all[i].EstImpactPct)
			}
			if all[j].EstImpactPct != nil {
				aj = absFloat(*all[j].EstImpactPct)
			}
			return ai > aj
		})
		var kept []string
		for _, d := range all {
			if d.EstImpactPct == nil || d.ReturnPct == nil {
				continue
			}
			if absFloat(*d.EstImpactPct) < 0.15 {
				continue
			}
			kept = append(kept, fmt.Sprintf("%s %+.1f%%→%+.2f%% (r=%+.2f)",
				d.LabelVi, *d.ReturnPct, *d.EstImpactPct, d.HistoricalR))
			if len(kept) >= 3 {
				break
			}
		}
		if len(kept) == 0 {
			continue
		}
		sign := "+"
		if a.NetImpact < 0 {
			sign = ""
		}
		sb = append(sb, fmt.Sprintf("- %s (ròng %s%.2f%%): %s", t, sign, a.NetImpact, joinStrings(kept, "; ")))
	}
	if len(sb) == 0 {
		return ""
	}
	return "\n\n## CHUỖI TÁC ĐỘNG NGOẠI SINH (biến động 20 ngày của yếu tố ngoài × tương quan lịch sử)\n" +
		joinStrings(sb, "\n") +
		"\nLưu ý: dùng các số này làm ngữ cảnh — NẾU tác động ròng một mã > ±1% thì phải nêu yếu tố chính trong phần rationale. KHÔNG bịa correlation ngoài danh sách trên."
}

// RegimeChangeFlags returns a text block flagging drivers that moved >10% in 20 days.
func RegimeChangeFlags(hp HistoryProvider) string {
	const threshold = 10.0

	type flag struct {
		symbol, labelVi string
		retPct          float64
	}

	labels := map[string]string{}
	for _, list := range exposureDefs {
		for _, m := range list {
			if _, ok := labels[m.Driver]; !ok {
				labels[m.Driver] = m.LabelVi
			}
		}
	}

	var flags []flag
	for sym, label := range labels {
		hist := hp(sym, 30)
		if len(hist) < 20 {
			continue
		}
		start := hist[len(hist)-20].Close
		end := hist[len(hist)-1].Close
		if start <= 0 {
			continue
		}
		ret := (end - start) / start * 100
		if absFloat(ret) < threshold {
			continue
		}
		flags = append(flags, flag{sym, label, ret})
	}
	if len(flags) == 0 {
		return ""
	}
	sort.Slice(flags, func(i, j int) bool { return absFloat(flags[i].retPct) > absFloat(flags[j].retPct) })

	var lines []string
	for _, f := range flags {
		dir := "TĂNG"
		if f.retPct < 0 {
			dir = "GIẢM"
		}
		lines = append(lines, fmt.Sprintf("- %s (%s): %s %.1f%% trong 20 ngày → regime shift",
			f.labelVi, f.symbol, dir, absFloat(f.retPct)))
	}
	return "\n\n## ⚡ REGIME SHIFT HÀNG HOÁ/TỶ GIÁ (biến động >10%% trong 20 ngày)\n" +
		joinStrings(lines, "\n") +
		"\n⚠️ Khi driver thay đổi regime, tương quan lịch sử có thể khuếch đại. Đánh giá lại các mã bị ảnh hưởng — đừng dựa vào tương quan thời kỳ ổn định."
}

// BounceSignals returns an oversold-bounce block for the plan prompt.
// rsiByTicker: map of VN ticker → current RSI. hp is used to check VIX trend.
func BounceSignals(rsiByTicker map[string]float64, hp HistoryProvider) string {
	if len(rsiByTicker) == 0 {
		return ""
	}

	// Check if VIX is falling over past 5 trading days.
	vixFalling := false
	hist := hp("^VIX", 10)
	if len(hist) >= 5 {
		start := hist[len(hist)-5].Close
		end := hist[len(hist)-1].Close
		if start > 0 && end < start {
			vixFalling = true
		}
	}

	var candidates []string
	for ticker, rsi := range rsiByTicker {
		if rsi >= 40 {
			continue
		}
		a := Analyze(ticker, hp)
		if !a.HasMapping {
			continue
		}
		if a.NetImpact < 0.5 {
			continue
		}
		note := fmt.Sprintf("- %s: RSI=%.0f (oversold), driver ròng %+.2f%%",
			ticker, rsi, a.NetImpact)
		if vixFalling {
			note += ", VIX đang giảm"
		}
		note += " → KHẢ NĂNG CAO HỒI KỸ THUẬT"
		candidates = append(candidates, note)
	}
	if len(candidates) == 0 {
		return ""
	}
	return "\n\n## 🔄 TÍN HIỆU OVERSOLD BOUNCE (từ audit lịch sử)\n" +
		joinStrings(candidates, "\n") +
		"\n⚠️ CẤM CUT_LOSS các mã trên — audit 30 ngày cho thấy bán đáy trong tình huống này sai 3/3 lần, lỗ cơ hội trung bình 6%."
}
