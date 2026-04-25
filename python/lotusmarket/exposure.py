"""Exposure chain — per-ticker mapping of external drivers with validated correlations.

Each entry carries historical Pearson r at the optimal horizon so callers can
compute a principled estimated-impact number rather than guessing.

Usage::

    from lotusmarket.exposure import analyze, format_for_prompt, unique_driver_symbols

    def my_provider(symbol: str, days: int) -> list[DriverClose]:
        ...  # return historical closes for symbol

    analysis = analyze("HPG", my_provider)
    print(analysis.net_impact)
"""

from __future__ import annotations

from dataclasses import dataclass, field
from typing import Callable, Dict, List, Optional

# ---------------------------------------------------------------------------
# Types
# ---------------------------------------------------------------------------


@dataclass
class DriverClose:
    """Single daily closing price for one driver symbol."""

    date: str  # YYYY-MM-DD
    close: float


# Callable[[symbol, days], List[DriverClose]] — results must be ascending by date.
HistoryProvider = Callable[[str, int], List[DriverClose]]

SIDE_INPUT = "input"
SIDE_PEER = "peer"
SIDE_OUTPUT = "output"
SIDE_CONTEXT = "context"


@dataclass
class ExposureMapping:
    driver: str
    label_vi: str
    category: str
    side: str
    sign: int
    historical_r: float
    horizon_days: int
    sample_n: int
    explanation: str
    concept_gloss: str


@dataclass
class DriverImpact:
    mapping: ExposureMapping
    latest_close: Optional[float] = None
    previous_close: Optional[float] = None
    return_pct: Optional[float] = None
    est_impact_pct: Optional[float] = None
    impact_magnitude: str = "none"
    direction: str = "flat"
    last_updated: str = ""


@dataclass
class ExposureAnalysis:
    ticker: str
    input: List[DriverImpact] = field(default_factory=list)
    peer: List[DriverImpact] = field(default_factory=list)
    output: List[DriverImpact] = field(default_factory=list)
    context: List[DriverImpact] = field(default_factory=list)
    net_impact: float = 0.0
    summary: str = ""
    has_mapping: bool = False


# ---------------------------------------------------------------------------
# Exposure definitions (ported from services/stock_exposure.go)
# ---------------------------------------------------------------------------


def _build_exposure_defs() -> Dict[str, List[ExposureMapping]]:
    m: Dict[str, List[ExposureMapping]] = {}

    # ─── Thép (HPG) ──────────────────────────────────────────────────────────
    m["HPG"] = [
        ExposureMapping(
            "000300.SS",
            "Chứng khoán Trung Quốc (CSI 300)",
            "DEMAND",
            SIDE_OUTPUT,
            +1,
            0.461,
            20,
            456,
            "Kinh tế TQ mạnh → bất động sản & hạ tầng TQ tiêu thụ nhiều thép → cầu thép toàn cầu tăng → HPG hưởng lợi.",
            "CSI 300 là chỉ số 300 cổ phiếu lớn nhất Trung Quốc — dùng làm 'nhiệt kế' cho kinh tế TQ.",
        ),
        ExposureMapping(
            "BHP",
            "Giá cổ phiếu BHP (công ty khai thác quặng)",
            "INPUT",
            SIDE_INPUT,
            +1,
            0.441,
            20,
            462,
            "BHP là công ty khai thác quặng sắt lớn nhất thế giới. Giá BHP phản ánh giá quặng sắt — nguyên liệu chính của HPG.",
            "'Quặng sắt' là đá chứa sắt, khi nấu trong lò cao sẽ ra thép. Chiếm ~30% chi phí sản xuất thép.",
        ),
        ExposureMapping(
            "600019.SS",
            "Baoshan Iron & Steel (đối thủ TQ)",
            "PEER",
            SIDE_PEER,
            +1,
            0.420,
            20,
            456,
            "Baoshan là hãng thép lớn nhất TQ. Khi giá Baoshan tăng, thường cả ngành thép Châu Á đang lên theo cùng chu kỳ.",
            "'Đối thủ cùng ngành' là công ty làm cùng sản phẩm — thường giá cổ phiếu di chuyển cùng chiều theo chu kỳ ngành.",
        ),
        ExposureMapping(
            "TIO=F",
            "Hợp đồng tương lai quặng sắt (SGX)",
            "INPUT",
            SIDE_INPUT,
            +1,
            0.278,
            20,
            463,
            "Giá quặng sắt giao dịch trên sàn Singapore — thước đo giá nguyên liệu đầu vào của HPG.",
            "'Hợp đồng tương lai' (futures) là thoả thuận mua/bán hàng hoá ở mức giá xác định trong tương lai.",
        ),
        ExposureMapping(
            "CNY=X",
            "Tỷ giá USD/CNY",
            "DEMAND",
            SIDE_OUTPUT,
            -1,
            -0.216,
            5,
            431,
            "CNY yếu đi (USD/CNY tăng) thường đi kèm với việc kinh tế TQ chậm lại → cầu thép yếu → HPG giảm.",
            "Khi CNY yếu, hàng xuất khẩu TQ rẻ hơn nhưng đồng thời báo hiệu kinh tế nội địa TQ đang gặp khó.",
        ),
    ]

    # ─── Ngân hàng — 14 mã VN30 ──────────────────────────────────────────────
    bank_pack = [
        ExposureMapping(
            "VNM",
            "Quỹ ETF VanEck Vietnam (VNM)",
            "MARKET",
            SIDE_CONTEXT,
            +1,
            0.70,
            20,
            460,
            "Là quỹ ETF tại Mỹ nắm giữ các cổ phiếu VN. Phản ánh xu hướng chung của toàn thị trường VN.",
            "ETF là quỹ mô phỏng một rổ cổ phiếu — tiền nước ngoài đầu tư vào VN thường đi qua ETF này.",
        ),
        ExposureMapping(
            "^VIX",
            "Chỉ số sợ hãi (VIX)",
            "RISK",
            SIDE_CONTEXT,
            -1,
            -0.30,
            20,
            455,
            "VIX tăng cao → nhà đầu tư toàn cầu rút tiền khỏi thị trường mới nổi (EM) → ngân hàng VN thường bị bán mạnh.",
            "VIX (thường gọi là 'chỉ số sợ hãi') đo mức biến động kỳ vọng của S&P 500. Cao = thị trường lo lắng.",
        ),
        ExposureMapping(
            "EEM",
            "Quỹ ETF thị trường mới nổi (EEM)",
            "FLOW",
            SIDE_CONTEXT,
            +1,
            0.37,
            20,
            462,
            "EEM phản ánh dòng tiền quốc tế vào các thị trường mới nổi (VN, TQ, Ấn Độ, Brazil...). Dòng tiền vào EM = tốt cho ngân hàng VN.",
            "'Thị trường mới nổi' (Emerging Markets) là các nước đang phát triển — vốn nhạy với rủi ro toàn cầu.",
        ),
    ]
    for t in [
        "VCB",
        "TCB",
        "BID",
        "CTG",
        "ACB",
        "MBB",
        "HDB",
        "LPB",
        "SHB",
        "SSB",
        "STB",
        "TPB",
        "VIB",
        "VPB",
    ]:
        m[t] = bank_pack

    # ─── Chứng khoán — SSI ────────────────────────────────────────────────────
    m["SSI"] = [
        ExposureMapping(
            "VNM",
            "Quỹ ETF VanEck Vietnam (VNM)",
            "MARKET",
            SIDE_CONTEXT,
            +1,
            0.747,
            20,
            460,
            "Thanh khoản thị trường là nguồn thu chính của công ty chứng khoán. VN-Index lên = giao dịch sôi động = SSI thu phí nhiều.",
            "'Thanh khoản' = lượng giao dịch mua bán trong ngày. Cao = người ta giao dịch nhiều.",
        ),
        ExposureMapping(
            "EEM",
            "Quỹ ETF thị trường mới nổi (EEM)",
            "FLOW",
            SIDE_CONTEXT,
            +1,
            0.415,
            20,
            455,
            "Khi dòng tiền quốc tế đổ vào EM, khối lượng giao dịch tại VN tăng → SSI hưởng lợi.",
            "",
        ),
        ExposureMapping(
            "^VIX",
            "Chỉ số sợ hãi (VIX)",
            "RISK",
            SIDE_CONTEXT,
            -1,
            -0.243,
            20,
            456,
            "VIX cao → nhà đầu tư cá nhân VN ngại giao dịch → doanh thu môi giới SSI giảm.",
            "",
        ),
    ]

    # ─── Công nghệ — FPT ──────────────────────────────────────────────────────
    m["FPT"] = [
        ExposureMapping(
            "QQQ",
            "Nasdaq 100 ETF (QQQ)",
            "SECTOR",
            SIDE_OUTPUT,
            +1,
            0.569,
            20,
            462,
            "FPT bán dịch vụ IT cho Mỹ, Nhật. Khi ngành công nghệ toàn cầu (Nasdaq) bùng nổ, khách hàng chi nhiều hơn → FPT tăng đơn hàng.",
            "QQQ là quỹ ETF theo 100 cổ phiếu công nghệ lớn nhất Mỹ (Apple, Microsoft, Nvidia...). Đo 'sức khoẻ' ngành công nghệ.",
        ),
        ExposureMapping(
            "XLK",
            "US Tech Sector ETF (XLK)",
            "SECTOR",
            SIDE_OUTPUT,
            +1,
            0.522,
            20,
            462,
            "Tương tự QQQ, đo xu hướng chi tiêu cho công nghệ.",
            "",
        ),
        ExposureMapping(
            "VNM",
            "Quỹ ETF VanEck Vietnam (VNM)",
            "MARKET",
            SIDE_CONTEXT,
            +1,
            0.377,
            1,
            481,
            "Beta thị trường — FPT là mã trụ trong VN-Index.",
            "",
        ),
    ]

    # ─── Bán lẻ — MWG ────────────────────────────────────────────────────────
    m["MWG"] = [
        ExposureMapping(
            "VNM",
            "Quỹ ETF VanEck Vietnam (VNM)",
            "MARKET",
            SIDE_CONTEXT,
            +1,
            0.635,
            20,
            459,
            "Beta thị trường chung.",
            "",
        ),
        ExposureMapping(
            "EEM",
            "ETF thị trường mới nổi (EEM)",
            "FLOW",
            SIDE_CONTEXT,
            +1,
            0.565,
            20,
            459,
            "Tâm lý tiêu dùng EM tốt = MWG hưởng lợi.",
            "",
        ),
        ExposureMapping(
            "DX-Y.NYB",
            "Chỉ số USD (DXY)",
            "COST",
            SIDE_INPUT,
            -1,
            -0.404,
            20,
            457,
            "MWG nhập iPhone, điện thoại, đồ điện tử — trả bằng USD. USD mạnh lên → giá nhập đắt → biên lợi nhuận giảm.",
            "DXY đo sức mạnh USD so với 6 đồng tiền lớn (EUR, JPY, GBP...). DXY tăng = USD mạnh lên toàn cầu.",
        ),
    ]

    # ─── Bất động sản — Vingroup ──────────────────────────────────────────────
    m["VHM"] = [
        ExposureMapping(
            "VNM",
            "Quỹ ETF VanEck Vietnam (VNM)",
            "MARKET",
            SIDE_CONTEXT,
            +1,
            0.682,
            20,
            459,
            "Beta thị trường.",
            "",
        ),
        ExposureMapping(
            "FXI",
            "ETF cổ phiếu lớn Trung Quốc (FXI)",
            "DEMAND",
            SIDE_OUTPUT,
            +1,
            0.156,
            20,
            459,
            "Tâm lý về BĐS TQ lan toả sang BĐS VN — khi BĐS TQ hồi phục, NĐT cũng tin vào BĐS VN hơn.",
            "",
        ),
    ]
    m["VIC"] = [
        ExposureMapping(
            "VNM",
            "Quỹ ETF VanEck Vietnam (VNM)",
            "MARKET",
            SIDE_CONTEXT,
            +1,
            0.551,
            5,
            477,
            "Beta thị trường.",
            "",
        ),
    ]
    m["VRE"] = [
        ExposureMapping(
            "VNM",
            "Quỹ ETF VanEck Vietnam (VNM)",
            "MARKET",
            SIDE_CONTEXT,
            +1,
            0.672,
            20,
            459,
            "Beta thị trường.",
            "",
        ),
    ]

    # ─── Năng lượng — GAS, PLX ────────────────────────────────────────────────
    m["GAS"] = [
        ExposureMapping(
            "VNM",
            "Quỹ ETF VanEck Vietnam (VNM)",
            "MARKET",
            SIDE_CONTEXT,
            +1,
            0.40,
            20,
            460,
            "Beta thị trường. Lưu ý: giá bán khí nội địa có cơ chế riêng, ít bị ảnh hưởng trực tiếp bởi dầu thế giới.",
            "",
        ),
    ]
    m["PLX"] = [
        ExposureMapping(
            "VNM",
            "Quỹ ETF VanEck Vietnam (VNM)",
            "MARKET",
            SIDE_CONTEXT,
            +1,
            0.40,
            20,
            460,
            "Beta thị trường. Giá xăng dầu bán lẻ VN do Bộ Công Thương quyết định theo chu kỳ 10 ngày, không biến động theo giờ như dầu WTI.",
            "",
        ),
    ]

    # ─── Tiêu dùng ────────────────────────────────────────────────────────────
    m["MSN"] = [
        ExposureMapping(
            "VNM",
            "Quỹ ETF VanEck Vietnam (VNM)",
            "MARKET",
            SIDE_CONTEXT,
            +1,
            0.636,
            20,
            459,
            "Beta thị trường.",
            "",
        ),
        ExposureMapping(
            "FXI",
            "ETF cổ phiếu lớn TQ (FXI)",
            "DEMAND",
            SIDE_OUTPUT,
            +1,
            0.448,
            20,
            459,
            "Masan là mảng tiêu dùng lớn; tâm lý tiêu dùng khu vực Châu Á (đại diện bởi TQ) ảnh hưởng đến MSN.",
            "",
        ),
    ]
    m["SAB"] = [
        ExposureMapping(
            "VNM",
            "Quỹ ETF VanEck Vietnam (VNM)",
            "MARKET",
            SIDE_CONTEXT,
            +1,
            0.427,
            5,
            474,
            "Beta thị trường. SAB ít phụ thuộc commodity — bia là sản phẩm 'defensive'.",
            "",
        ),
    ]
    # Note: key "VNM" = Vinamilk VN stock code; driver "VNM" = VanEck ETF Yahoo symbol.
    m["VNM"] = [
        ExposureMapping(
            "VNM",
            "Quỹ ETF VanEck Vietnam (VNM ETF)",
            "MARKET",
            SIDE_CONTEXT,
            +1,
            0.469,
            5,
            474,
            "Beta thị trường chung của Vinamilk.",
            "Trùng tên viết tắt — mã VNM là Vinamilk (VN), còn 'VNM ETF' là quỹ Mỹ đầu tư vào VN.",
        ),
        ExposureMapping(
            "DX-Y.NYB",
            "Chỉ số USD (DXY)",
            "COST",
            SIDE_INPUT,
            -1,
            -0.291,
            20,
            462,
            "Vinamilk nhập sữa bột (New Zealand, EU) trả bằng USD. USD mạnh = nguyên liệu đắt = biên lợi nhuận co lại.",
            "",
        ),
    ]

    # ─── Hoá chất — DGC ──────────────────────────────────────────────────────
    m["DGC"] = [
        ExposureMapping(
            "MOS",
            "Mosaic (công ty phân bón Mỹ)",
            "PEER",
            SIDE_PEER,
            +1,
            0.267,
            5,
            474,
            "Mosaic là hãng sản xuất phân bón phosphate lớn nhất thế giới. DGC cùng sản xuất phốt pho/phân bón nên di chuyển theo chu kỳ ngành hoá chất.",
            "'Phốt pho vàng' (yellow phosphorus) là nguyên liệu cho phân bón, pin lithium, dệt may — sản phẩm chủ lực của DGC.",
        ),
        ExposureMapping(
            "NTR",
            "Nutrien (công ty phân bón Canada)",
            "PEER",
            SIDE_PEER,
            +1,
            0.229,
            5,
            473,
            "Nutrien là đối thủ toàn cầu trong mảng phân bón/hoá chất nông nghiệp.",
            "",
        ),
        ExposureMapping(
            "CNY=X",
            "Tỷ giá USD/CNY",
            "DEMAND",
            SIDE_OUTPUT,
            -1,
            -0.159,
            5,
            428,
            "TQ là thị trường xuất khẩu lớn của DGC. CNY yếu = TQ mua ít hơn.",
            "",
        ),
    ]

    # ─── Nông nghiệp — GVR (cao su) ───────────────────────────────────────────
    m["GVR"] = [
        ExposureMapping(
            "5108.T",
            "Bridgestone (hãng lốp Nhật)",
            "DEMAND",
            SIDE_OUTPUT,
            +1,
            0.317,
            20,
            446,
            "GVR bán cao su thô. Khách hàng lớn là các hãng lốp xe. Bridgestone tăng = cầu lốp tốt = cầu cao su tốt.",
            "Cao su thiên nhiên chủ yếu dùng sản xuất lốp ô tô — chiếm 70% nhu cầu cao su toàn cầu.",
        ),
        ExposureMapping(
            "GT",
            "Goodyear (hãng lốp Mỹ)",
            "DEMAND",
            SIDE_OUTPUT,
            +1,
            0.188,
            20,
            454,
            "Tương tự Bridgestone — Goodyear là proxy cho cầu lốp xe thế giới.",
            "",
        ),
    ]

    # ─── Hàng không — VJC ────────────────────────────────────────────────────
    m["VJC"] = [
        ExposureMapping(
            "VNM",
            "Quỹ ETF VanEck Vietnam (VNM)",
            "MARKET",
            SIDE_CONTEXT,
            +1,
            0.520,
            20,
            459,
            "Beta thị trường.",
            "",
        ),
        ExposureMapping(
            "CL=F",
            "Dầu WTI",
            "COST",
            SIDE_INPUT,
            -1,
            -0.254,
            20,
            462,
            "Nhiên liệu máy bay chiếm 30-40% chi phí vận hành hãng hàng không. Dầu tăng = chi phí tăng = VJC giảm.",
            "WTI (West Texas Intermediate) là loại dầu tiêu chuẩn của Mỹ — benchmark giá dầu toàn cầu.",
        ),
        ExposureMapping(
            "HO=F",
            "Dầu đốt (proxy cho nhiên liệu bay)",
            "COST",
            SIDE_INPUT,
            -1,
            -0.251,
            20,
            462,
            "Heating oil là dẫn xuất gần với jet fuel (dầu máy bay) hơn WTI.",
            "",
        ),
    ]

    return m


_EXPOSURE_DEFS: Dict[str, List[ExposureMapping]] = _build_exposure_defs()


# ---------------------------------------------------------------------------
# Public API
# ---------------------------------------------------------------------------


def unique_driver_symbols() -> List[str]:
    """Return every Yahoo Finance symbol referenced by at least one mapping (sorted)."""
    seen: set = set()
    for mappings in _EXPOSURE_DEFS.values():
        for m in mappings:
            seen.add(m.driver)
    return sorted(seen)


def get_mappings_for_ticker(ticker: str) -> List[ExposureMapping]:
    """Return mappings for a VN ticker, or [] if unmapped."""
    return _EXPOSURE_DEFS.get(ticker, [])


def _abs(x: float) -> float:
    return x if x >= 0 else -x


def _build_summary(ticker: str, analysis: ExposureAnalysis) -> str:
    all_impacts = analysis.input + analysis.peer + analysis.output + analysis.context
    top: Optional[DriverImpact] = None
    for d in all_impacts:
        if d.est_impact_pct is None:
            continue
        if top is None or _abs(d.est_impact_pct) > _abs(top.est_impact_pct):  # type: ignore[arg-type]
            top = d
    if top is None or top.est_impact_pct is None:
        return f"Chưa có dữ liệu gần nhất cho các yếu tố ảnh hưởng đến {ticker}."
    direction = "tăng" if analysis.net_impact >= 0 else "giảm"
    return (
        f"Dựa trên biến động {top.mapping.horizon_days} ngày của các yếu tố liên quan, "
        f"{ticker} có áp lực {direction} khoảng {_abs(analysis.net_impact):.2f}%. "
        f"Yếu tố mạnh nhất: {top.mapping.label_vi} "
        f"({top.return_pct:+.1f}% → tác động ước tính {top.est_impact_pct:+.2f}%)."
    )


def analyze(ticker: str, history_provider: HistoryProvider) -> ExposureAnalysis:
    """Compute driver impacts for a VN ticker using provided history.

    history_provider(symbol, days) must return DriverClose list ascending by date.
    """
    mappings = get_mappings_for_ticker(ticker)
    if not mappings:
        return ExposureAnalysis(
            ticker=ticker,
            has_mapping=False,
            summary=(
                "Chưa có dữ liệu tương quan cho mã này. Chỉ các mã đã được backtest "
                "(HPG, ngân hàng VN30, FPT, MWG, DGC, GVR, VJC, VNM, MSN...) mới có phân tích chuỗi tác động."
            ),
        )

    result = ExposureAnalysis(ticker=ticker, has_mapping=True)
    net_weighted = 0.0
    total_weight = 0.0

    for em in mappings:
        impact = DriverImpact(mapping=em)
        history = history_provider(em.driver, em.horizon_days * 3 + 10)

        if len(history) >= 2:
            latest = history[-1]
            prev_idx = max(0, len(history) - 1 - em.horizon_days)
            prev = history[prev_idx]
            if prev.close > 0:
                ret = (latest.close - prev.close) / prev.close * 100
                impact.latest_close = latest.close
                impact.previous_close = prev.close
                impact.return_pct = ret
                impact.last_updated = latest.date

                est = ret * em.historical_r
                impact.est_impact_pct = est

                if est > 0.3:
                    impact.direction = "up"
                elif est < -0.3:
                    impact.direction = "down"

                mag = _abs(est)
                if mag > 1.5:
                    impact.impact_magnitude = "large"
                elif mag > 0.5:
                    impact.impact_magnitude = "medium"
                elif mag > 0.1:
                    impact.impact_magnitude = "small"

                w = _abs(em.historical_r)
                net_weighted += est * w
                total_weight += w

        if em.side == SIDE_INPUT:
            result.input.append(impact)
        elif em.side == SIDE_PEER:
            result.peer.append(impact)
        elif em.side == SIDE_OUTPUT:
            result.output.append(impact)
        else:
            result.context.append(impact)

    if total_weight > 0:
        result.net_impact = net_weighted / total_weight
    result.summary = _build_summary(ticker, result)
    return result


def format_for_prompt(tickers: List[str], history_provider: HistoryProvider) -> str:
    """Return a compact text block for Claude prompts describing top driver moves."""
    if not tickers:
        return ""
    lines: List[str] = []
    for t in tickers:
        a = analyze(t, history_provider)
        if not a.has_mapping:
            continue
        all_impacts = a.input + a.peer + a.output + a.context
        all_impacts.sort(
            key=lambda d: _abs(d.est_impact_pct)
            if d.est_impact_pct is not None
            else 0.0,
            reverse=True,
        )
        kept: List[str] = []
        for d in all_impacts:
            if d.est_impact_pct is None or d.return_pct is None:
                continue
            if _abs(d.est_impact_pct) < 0.15:
                continue
            kept.append(
                f"{d.mapping.label_vi} {d.return_pct:+.1f}%→{d.est_impact_pct:+.2f}% (r={d.mapping.historical_r:+.2f})"
            )
            if len(kept) >= 3:
                break
        if not kept:
            continue
        sign = "+" if a.net_impact >= 0 else ""
        lines.append(f"- {t} (ròng {sign}{a.net_impact:.2f}%): {'; '.join(kept)}")
    if not lines:
        return ""
    return (
        "\n\n## CHUỖI TÁC ĐỘNG NGOẠI SINH (biến động 20 ngày của yếu tố ngoài × tương quan lịch sử)\n"
        + "\n".join(lines)
        + "\nLưu ý: dùng các số này làm ngữ cảnh — NẾU tác động ròng một mã > ±1% thì phải nêu yếu tố chính trong phần rationale. KHÔNG bịa correlation ngoài danh sách trên."
    )


def regime_change_flags(history_provider: HistoryProvider) -> str:
    """Return a text block flagging drivers that moved >10% in 20 days."""
    threshold = 10.0
    labels: Dict[str, str] = {}
    for mappings in _EXPOSURE_DEFS.values():
        for em in mappings:
            if em.driver not in labels:
                labels[em.driver] = em.label_vi

    flags: List[tuple] = []
    for sym, label in labels.items():
        hist = history_provider(sym, 30)
        if len(hist) < 20:
            continue
        start = hist[-20].close
        end = hist[-1].close
        if start <= 0:
            continue
        ret = (end - start) / start * 100
        if _abs(ret) < threshold:
            continue
        flags.append((sym, label, ret))

    if not flags:
        return ""
    flags.sort(key=lambda x: _abs(x[2]), reverse=True)
    lines = []
    for sym, label, ret in flags:
        direction = "TĂNG" if ret > 0 else "GIẢM"
        lines.append(
            f"- {label} ({sym}): {direction} {_abs(ret):.1f}% trong 20 ngày → regime shift"
        )
    return (
        "\n\n## ⚡ REGIME SHIFT HÀNG HOÁ/TỶ GIÁ (biến động >10% trong 20 ngày)\n"
        + "\n".join(lines)
        + "\n⚠️ Khi driver thay đổi regime, tương quan lịch sử có thể khuếch đại. Đánh giá lại các mã bị ảnh hưởng — đừng dựa vào tương quan thời kỳ ổn định."
    )


def bounce_signals(
    rsi_by_ticker: Dict[str, float], history_provider: HistoryProvider
) -> str:
    """Return an oversold-bounce text block for the plan prompt."""
    if not rsi_by_ticker:
        return ""

    vix_falling = False
    vix_hist = history_provider("^VIX", 10)
    if len(vix_hist) >= 5:
        start = vix_hist[-5].close
        end = vix_hist[-1].close
        if start > 0 and end < start:
            vix_falling = True

    candidates: List[str] = []
    for ticker, rsi in rsi_by_ticker.items():
        if rsi >= 40:
            continue
        a = analyze(ticker, history_provider)
        if not a.has_mapping or a.net_impact < 0.5:
            continue
        note = f"- {ticker}: RSI={rsi:.0f} (oversold), driver ròng {a.net_impact:+.2f}%"
        if vix_falling:
            note += ", VIX đang giảm"
        note += " → KHẢ NĂNG CAO HỒI KỸ THUẬT"
        candidates.append(note)

    if not candidates:
        return ""
    return (
        "\n\n## 🔄 TÍN HIỆU OVERSOLD BOUNCE (từ audit lịch sử)\n"
        + "\n".join(candidates)
        + "\n⚠️ CẤM CUT_LOSS các mã trên — audit 30 ngày cho thấy bán đáy trong tình huống này sai 3/3 lần, lỗ cơ hội trung bình 6%."
    )
