"""Market pulse scoring, sector flow ranking, risk indicators, regime classification."""

from dataclasses import dataclass, field
from typing import List, Tuple
import pandas as pd
from lotusmarket.types import get_sector


@dataclass
class SectorFlow:
    sector: str = ""
    rank: int = 0
    signal: str = ""  # "inflow"/"neutral"/"outflow"
    net_foreign: float = 0.0
    breadth: str = ""
    top_stocks: list = field(default_factory=list)
    score: float = 0.0


def score_breadth(ratio: float) -> int:
    if ratio > 60:
        return 80
    if ratio >= 40:
        return 50
    return 20


def score_foreign_flow(net_total: float) -> int:
    threshold = 100e9
    if net_total > threshold:
        return 90
    if net_total > 0:
        return 60 + int((net_total / threshold) * 30)
    if net_total == 0:
        return 50
    abs_val = -net_total
    if abs_val > threshold:
        return 20
    return 50 - int((abs_val / threshold) * 30)


def pulse_score(
    breadth_score: int, foreign_flow_score: int, volume_score: int, risk_score: int
) -> Tuple[int, str]:
    s = (
        breadth_score * 30
        + foreign_flow_score * 25
        + volume_score * 20
        + risk_score * 25
    ) // 100
    if s > 60:
        signal = "green"
    elif s >= 40:
        signal = "yellow"
    else:
        signal = "red"
    return s, signal


def score_vix(vix: float) -> int:
    if vix < 15:
        return 10
    if vix < 20:
        return 30
    if vix < 25:
        return 50
    if vix < 30:
        return 70
    return 90


def score_yield_curve(spread: float) -> int:
    if spread > 1.0:
        return 10
    if spread > 0.5:
        return 30
    if spread > 0:
        return 50
    return 90


def score_news_sentiment(negative: int, total: int) -> int:
    if total == 0:
        return 20
    ratio = negative / total * 100
    if ratio > 60:
        return 80
    if ratio > 40:
        return 50
    return 20


def score_fed_trend(current: float, previous: float) -> int:
    diff = current - previous
    if diff > 0.05:
        return 70
    if diff < -0.05:
        return 15
    return 40


def sector_flow(df: pd.DataFrame) -> pd.DataFrame:
    """Rank sectors by composite flow score.
    Input: DataFrame with ticker, close, volume, change_percent, foreign_net_vol columns.
    Returns: DataFrame with sector, rank, signal, net_foreign, avg_change_pct,
    total_volume_vnd, advances_count, total_count, breadth, score columns.
    """
    if df.empty:
        return pd.DataFrame()
    df2 = df.copy()
    df2["sector"] = df2["ticker"].map(get_sector)
    df2 = df2[df2["sector"] != "Others"]
    if df2.empty:
        return pd.DataFrame()
    df2["flow_value"] = df2["foreign_net_vol"] * df2["close"]
    df2["trade_value"] = df2["volume"] * df2["close"]
    grouped = (
        df2.groupby("sector")
        .agg(
            total_flow=("flow_value", "sum"),
            advances=("change_percent", lambda x: (x > 0).sum()),
            total=("ticker", "count"),
            total_volume=("volume", "sum"),
            total_volume_vnd=("trade_value", "sum"),
            avg_change_pct=("change_percent", "mean"),
        )
        .reset_index()
    )
    # Normalize
    fr = grouped["total_flow"].max() - grouped["total_flow"].min()
    if fr == 0:
        grouped["norm_flow"] = 50.0
    else:
        grouped["norm_flow"] = (
            (grouped["total_flow"] - grouped["total_flow"].min()) / fr * 100
        )
    vr = grouped["total_volume"].max() - grouped["total_volume"].min()
    if vr == 0:
        grouped["norm_vol"] = 50.0
    else:
        grouped["norm_vol"] = (
            (grouped["total_volume"] - grouped["total_volume"].min()) / vr * 100
        )
    grouped["breadth_ratio"] = grouped["advances"] / grouped["total"]
    grouped["score"] = (
        grouped["norm_flow"] * 0.4
        + grouped["breadth_ratio"] * 100 * 0.3
        + grouped["norm_vol"] * 0.3
    )
    grouped = grouped.sort_values("score", ascending=False).reset_index(drop=True)
    grouped["rank"] = grouped.index + 1
    grouped["signal"] = grouped["score"].apply(
        lambda s: "inflow" if s > 60 else ("outflow" if s < 40 else "neutral")
    )
    grouped["breadth"] = grouped.apply(
        lambda r: f"{int(r['advances'])}/{int(r['total'])} tăng", axis=1
    )
    return grouped[
        [
            "sector",
            "rank",
            "signal",
            "total_flow",
            "avg_change_pct",
            "total_volume_vnd",
            "advances",
            "total",
            "breadth",
            "score",
        ]
    ].rename(
        columns={
            "total_flow": "net_foreign",
            "advances": "advances_count",
            "total": "total_count",
        }
    )


# ============================================================
# Price Driver Attribution
# ============================================================


@dataclass
class DriverFeatures:
    """Raw daily features for one ticker — input to attribute_drivers()."""

    price_change_pct: float = 0.0
    foreign_net_vnd: float = 0.0
    market_delta_pct: float = 0.0
    sector_delta_pct: float = 0.0
    news_count: int = 0
    news_sentiment_avg: float = 0.0  # [-1, 1]
    rsi: float = 0.0
    ma_signal: str = ""  # "BUY" | "SELL" | "HOLD" | ""
    volume_ratio: float = 0.0  # today / MA20


@dataclass
class Attribution:
    """% contribution per driver. Sum + residual_pct = 100."""

    news_pct: float = 0.0
    foreign_pct: float = 0.0
    sector_pct: float = 0.0
    market_pct: float = 0.0
    technical_pct: float = 0.0
    volume_pct: float = 0.0
    residual_pct: float = 0.0
    dominant: str = "UNKNOWN"


def attribute_drivers(f: DriverFeatures) -> Attribution:
    """Decompose a stock's daily move into driver contributions.

    Algorithm (fully deterministic, zero external calls):
      1. Map each raw feature into a signed "impact" number.
      2. Aligned drivers (sign matches today's move) get full weight;
         misaligned get 30% credit (conflicting signal).
      3. Normalise to sum = 100%. Dominant driver = argmax(aligned weight).
      4. Residual = 100 - sum (unexplained portion).
    """
    import math

    move = f.price_change_pct
    if abs(move) < 0.05:
        return Attribution(dominant="UNKNOWN")

    def _clamp(v: float, limit: float) -> float:
        return max(-limit, min(limit, v))

    def _tech(rsi: float, ma: str) -> float:
        s = 2.0 if ma == "BUY" else (-2.0 if ma == "SELL" else 0.0)
        if rsi > 70:
            s -= 1
        elif 0 < rsi < 30:
            s += 1
        return s

    def _volume(ratio: float, mv: float) -> float:
        if ratio < 1.2:
            return 0.0
        impact = math.log(ratio) * 2
        return -impact if mv < 0 else impact

    signs = {
        "FOREIGN": _clamp(f.foreign_net_vnd / 10_000_000_000, 5),
        "MARKET": f.market_delta_pct,
        "SECTOR": f.sector_delta_pct - f.market_delta_pct,  # alpha
        "NEWS": f.news_sentiment_avg * math.log1p(f.news_count) * 3,
        "TECHNICAL": _tech(f.rsi, f.ma_signal),
        "VOLUME": _volume(f.volume_ratio, move),
    }

    move_sign = 1 if move > 0 else -1

    weights: dict = {}
    total = 0.0
    for k, v in signs.items():
        w = abs(v)
        if w == 0:
            continue
        vs = 1 if v > 0 else (-1 if v < 0 else 0)
        if vs != move_sign:
            w *= 0.3
        weights[k] = w
        total += w

    if total < 0.01:
        return Attribution(dominant="UNKNOWN")

    def _r1(v: float) -> float:
        return round(v * 10) / 10

    attr = Attribution(
        foreign_pct=_r1(weights.get("FOREIGN", 0) / total * 100),
        market_pct=_r1(weights.get("MARKET", 0) / total * 100),
        sector_pct=_r1(weights.get("SECTOR", 0) / total * 100),
        news_pct=_r1(weights.get("NEWS", 0) / total * 100),
        technical_pct=_r1(weights.get("TECHNICAL", 0) / total * 100),
        volume_pct=_r1(weights.get("VOLUME", 0) / total * 100),
    )
    attr.residual_pct = _r1(
        max(
            0,
            100
            - attr.news_pct
            - attr.foreign_pct
            - attr.sector_pct
            - attr.market_pct
            - attr.technical_pct
            - attr.volume_pct,
        )
    )

    best, best_w = "UNKNOWN", 0.0
    for k, w in weights.items():
        v = signs[k]
        vs = 1 if v > 0 else (-1 if v < 0 else 0)
        if vs != move_sign:
            continue
        if w > best_w:
            best_w, best = w, k
    attr.dominant = best
    return attr


# ============================================================
# Market Regime Classification
# ============================================================

REGIME_STABLE = "STABLE"
REGIME_VOLATILE = "VOLATILE"
REGIME_CRISIS = "CRISIS"
REGIME_EUPHORIA = "EUPHORIA"

CRISIS_PANIC = "CRISIS_PANIC"
CRISIS_FUNDAMENTAL = "CRISIS_FUNDAMENTAL"

_FUNDAMENTAL_KEYWORDS = [
    "bắt giữ", "khởi tố", "truy tố", "lừa đảo", "rút ruột",
    "bank run", "vỡ nợ", "phá sản", "hủy niêm yết",
    "kiểm soát đặc biệt", "đình chỉ",
]


@dataclass
class RegimeSignals:
    vix: float = 0.0
    vix_score: int = 0
    vnindex_change: float = 0.0
    vnindex_score: int = 0
    foreign_flow: float = 0.0   # billions VND
    flow_score: int = 0
    news_tier: int = 0
    news_score: int = 0
    news_detail: str = ""
    global_drop: float = 0.0
    global_score: int = 0
    composite_score: int = 0
    direction: str = "mixed"    # "down" | "up" | "mixed"


def score_regime_signals(
    vix: float,
    vnindex_change: float,
    foreign_flow_billions: float,
    news_tier: int,
    news_detail: str,
    global_max_drop_pct: float,
) -> RegimeSignals:
    """Convert raw inputs into a RegimeSignals struct with all scores computed.

    Args:
        vix: Current VIX level.
        vnindex_change: VN-Index daily change percent (negative = down).
        foreign_flow_billions: Net foreign flow in billions VND (negative = net sell).
        news_tier: Highest active news tier (1 = arrest/circuit breaker, 2 = Fed/tariffs, 3+ = lower, 0 = none).
        news_detail: First matching TIER 1/2 headline (used for crisis sub-type detection).
        global_max_drop_pct: Most negative single-day drop among major global indices (<=0).

    Returns:
        RegimeSignals with composite_score and direction set.
    """
    s = RegimeSignals(
        vix=vix,
        vnindex_change=vnindex_change,
        foreign_flow=foreign_flow_billions,
        news_tier=news_tier,
        news_detail=news_detail,
        global_drop=global_max_drop_pct,
    )

    # VIX score
    if vix >= 30:
        s.vix_score = 80
    elif vix >= 25:
        s.vix_score = 50
    elif vix >= 20:
        s.vix_score = 25
    else:
        s.vix_score = 0

    # VN-Index score (absolute change)
    abs_change = abs(vnindex_change)
    if abs_change >= 3.0:
        s.vnindex_score = 80
    elif abs_change >= 2.0:
        s.vnindex_score = 50
    elif abs_change >= 1.0:
        s.vnindex_score = 25
    else:
        s.vnindex_score = 0

    if vnindex_change < -1.0:
        s.direction = "down"
    elif vnindex_change > 1.0:
        s.direction = "up"
    else:
        s.direction = "mixed"

    # Foreign flow score
    if foreign_flow_billions < -2000:
        s.flow_score = 80
    elif foreign_flow_billions < -1000:
        s.flow_score = 50
    elif foreign_flow_billions < -500:
        s.flow_score = 25
    else:
        s.flow_score = 0

    # News score
    if news_tier == 1:
        s.news_score = 100
    elif news_tier == 2:
        s.news_score = 50
    else:
        s.news_score = 0

    # Global contagion score
    if global_max_drop_pct <= -5.0:
        s.global_score = 80
    elif global_max_drop_pct <= -3.0:
        s.global_score = 50
    else:
        s.global_score = 0

    s.composite_score = max(s.vix_score, s.vnindex_score, s.flow_score, s.news_score, s.global_score)
    return s


def _classify_crisis_type(s: RegimeSignals) -> str:
    news_lower = s.news_detail.lower()
    for kw in _FUNDAMENTAL_KEYWORDS:
        if kw in news_lower:
            return CRISIS_FUNDAMENTAL
    if s.global_score >= 80:
        return CRISIS_PANIC
    if s.flow_score >= 80 and s.news_score < 100:
        return CRISIS_PANIC
    return CRISIS_PANIC


def classify_regime(s: RegimeSignals) -> Tuple[str, str]:
    """Determine (regime, sub_type) from pre-computed RegimeSignals.

    Returns:
        (regime, sub_type) — sub_type is "" for non-crisis regimes.
    """
    if s.composite_score >= 80:
        if s.direction == "up":
            return REGIME_EUPHORIA, ""
        return REGIME_CRISIS, _classify_crisis_type(s)
    if s.composite_score >= 50:
        return REGIME_VOLATILE, ""
    return REGIME_STABLE, ""


def identify_trigger(s: RegimeSignals) -> Tuple[str, str]:
    """Return (source, detail) for the primary trigger of the regime."""
    max_score = 0
    source = "composite"
    detail = ""

    if s.news_score > max_score:
        max_score = s.news_score
        source = "news_tier1"
        detail = s.news_detail
    if s.vnindex_score > max_score:
        max_score = s.vnindex_score
        source = "vnindex_change"
        detail = f"VN-Index thay đổi {s.vnindex_change:.1f}%"
    if s.global_score > max_score:
        max_score = s.global_score
        source = "global_contagion"
        detail = f"Global index giảm {s.global_drop:.1f}%"
    if s.vix_score > max_score:
        max_score = s.vix_score
        source = "vix"
        detail = f"VIX = {s.vix:.1f}"
    if s.flow_score > max_score:
        source = "foreign_flow"
        detail = f"Foreign net flow: {s.foreign_flow:.0f}B VND"

    return source, detail
