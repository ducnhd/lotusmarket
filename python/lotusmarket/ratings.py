"""6-dimensional star ratings for a stock — deterministic.

Mirrors SSI iResearch's "Đánh giá kỹ thuật" panel:
  - Price strength (1-5)
  - Trend strength (1-5)
  - Short-term position (RSI) (1-5)
  - Money flow (1-5)
  - Volatility (1-5)
  - Base range (1-5)

Plus an overall 0-100 gauge and verdict (Outperform / Neutral / Underperform).

All math is deterministic. No AI, no external data, no I/O.
"""

from __future__ import annotations

import math
from dataclasses import dataclass
from typing import Optional, Sequence


@dataclass
class Ratings:
    price_strength: int
    trend_strength: int
    short_term_pos: int
    money_flow: int
    volatility_rating: int
    base_range: int
    overall_gauge: int
    overall_verdict: str


def compute(
    closes: Sequence[float],
    volumes: Optional[Sequence[int]],
    score: float,
    rsi: float,
    ma20: Optional[float],
    ma50: Optional[float],
    ma200: Optional[float],
) -> Ratings:
    """Compute 6-dim star ratings.

    Args:
        closes: oldest-first daily close prices (need >= 20 for full output).
        volumes: parallel daily volumes (oldest-first); may be None.
        score: external 0-100 technical signal score (e.g. from technical.score).
        rsi: 14-day RSI (0-100).
        ma20 / ma50 / ma200: moving averages; may be None for short histories.
    """
    if len(closes) < 20:
        return Ratings(3, 3, 3, 3, 3, 3, 50, "Neutral")

    current = closes[-1]
    ps = _price_strength_stars(score)
    ts = _trend_stars(current, ma20, ma50, ma200)
    sp = _rsi_stars(rsi)
    mf = _money_flow_stars(closes, volumes)
    vol = _volatility_stars(closes)
    br = _base_range_stars(current, ma20)

    total = ps + ts + sp + mf + vol + br
    gauge = total * 100 // 30
    if gauge >= 60:
        verdict = "Outperform"
    elif gauge >= 40:
        verdict = "Neutral"
    else:
        verdict = "Underperform"
    return Ratings(ps, ts, sp, mf, vol, br, gauge, verdict)


def _price_strength_stars(score: float) -> int:
    if score >= 80:
        return 5
    if score >= 60:
        return 4
    if score >= 40:
        return 3
    if score >= 20:
        return 2
    return 1


def _trend_stars(
    price: float,
    ma20: Optional[float],
    ma50: Optional[float],
    ma200: Optional[float],
) -> int:
    if ma20 is None:
        return 3
    above20 = price > ma20
    above50 = ma50 is not None and price > ma50
    above200 = ma200 is not None and price > ma200
    stacked = ma50 is not None and ma200 is not None and ma20 > ma50 > ma200
    if above20 and above50 and above200 and stacked:
        return 5
    if above20 and above50 and above200:
        return 4
    if above20 and above50:
        return 3
    if above20:
        return 2
    return 1


def _rsi_stars(rsi: float) -> int:
    if 50 <= rsi <= 65:
        return 5
    if (45 <= rsi < 50) or (65 < rsi <= 70):
        return 4
    if (35 <= rsi < 45) or (70 < rsi <= 75):
        return 3
    if (25 <= rsi < 35) or (75 < rsi <= 80):
        return 2
    return 1


def _money_flow_stars(closes: Sequence[float], volumes: Optional[Sequence[int]]) -> int:
    if volumes is None or len(volumes) < 21 or len(closes) < 6:
        return 3
    avg20 = sum(volumes[-21:-1]) / 20
    if avg20 <= 0:
        return 3
    ratio = volumes[-1] / avg20
    ret5 = (closes[-1] - closes[-6]) / closes[-6] * 100
    price_up = ret5 > 0.5
    price_down = ret5 < -0.5
    if ratio >= 1.5 and price_up:
        return 5
    if ratio >= 1.2 and price_up:
        return 4
    if 0.8 <= ratio < 1.2:
        return 3
    if ratio >= 1.2 and price_down:
        return 1  # distribution
    if ratio < 0.8:
        return 2
    return 3


def _volatility_stars(closes: Sequence[float]) -> int:
    n = len(closes)
    if n < 21:
        return 3
    rets = []
    for i in range(n - 20, n):
        if closes[i - 1] <= 0:
            continue
        rets.append(math.log(closes[i] / closes[i - 1]))
    if len(rets) < 5:
        return 3
    mean = sum(rets) / len(rets)
    variance = sum((r - mean) ** 2 for r in rets) / (len(rets) - 1)
    annualized = math.sqrt(variance) * math.sqrt(252) * 100  # percent
    if 15 <= annualized <= 25:
        return 5
    if (10 <= annualized < 15) or (25 < annualized <= 35):
        return 4
    if 35 < annualized <= 45:
        return 3
    if 45 < annualized <= 60:
        return 2
    return 1


def _base_range_stars(current: float, ma20: Optional[float]) -> int:
    if ma20 is None or ma20 <= 0:
        return 3
    dist = abs(current - ma20) / ma20 * 100
    if dist <= 2:
        return 5
    if dist <= 4:
        return 4
    if dist <= 7:
        return 3
    if dist <= 12:
        return 2
    return 1
