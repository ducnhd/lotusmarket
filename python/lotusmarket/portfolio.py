"""Portfolio returns and confidence scoring."""

from dataclasses import dataclass, field
from typing import Dict, List, Optional
import math


@dataclass
class ReturnSegment:
    start_value: float
    end_value: float


@dataclass
class ConfidenceDetail:
    score: int = 0
    components: Dict[str, int] = field(default_factory=dict)


@dataclass
class PositionAnalysis:
    action: str = ""  # BUY_MORE, SELL_PARTIAL, SELL_ALL, CUT_LOSS, HOLD
    timeframe_align: str = ""  # bullish, bearish, mixed
    rsi: float = 0.0
    ma20: Optional[float] = None
    ma50: Optional[float] = None
    ma200: Optional[float] = None
    price_vs_ma20: str = ""  # above, below
    price_vs_ma50: str = ""  # above, below
    pe: float = 0.0
    beta: float = 0.0
    dividend_yield: float = 0.0
    foreign_net_vol: int = 0
    buy_sell_ratio: float = 0.0


def twr(segments: List[ReturnSegment]) -> float:
    """Time-weighted return. TWR = Product(1 + Ri) - 1."""
    if not segments:
        return 0.0
    product = 1.0
    for seg in segments:
        if seg.start_value == 0:
            return 0.0
        r = (seg.end_value - seg.start_value) / seg.start_value
        product *= 1 + r
    return product - 1


def position_confidence(
    pa: PositionAnalysis, news_sentiment: str = ""
) -> ConfidenceDetail:
    """6-factor weighted confidence score (0-100).
    Weights: timeframe 25%, RSI 20%, MA trend 20%, fundamentals 15%, flow 10%, news 10%.
    """
    c = {}
    is_buy = pa.action == "BUY_MORE"
    is_sell = pa.action in ("SELL_PARTIAL", "SELL_ALL", "CUT_LOSS")

    # 1. Timeframe alignment (25%)
    if pa.timeframe_align == "bullish":
        c["timeframe"] = 100 if is_buy else (30 if is_sell else 60)
    elif pa.timeframe_align == "bearish":
        c["timeframe"] = 100 if is_sell else (30 if is_buy else 60)
    else:
        c["timeframe"] = 40

    # 2. RSI zone (20%)
    if is_buy:
        c["rsi"] = 90 if pa.rsi < 30 else (60 if pa.rsi < 50 else 30)
    elif is_sell:
        c["rsi"] = 90 if pa.rsi > 70 else (60 if pa.rsi > 50 else 30)
    else:
        c["rsi"] = 70 if 30 <= pa.rsi <= 70 else 40

    # 3. MA trend (20%)
    ma_score = 40
    if pa.ma20 is not None and pa.ma50 is not None and pa.ma200 is not None:
        if pa.price_vs_ma20 == "above" and pa.price_vs_ma50 == "above":
            ma_score = 100 if (is_buy or pa.action == "HOLD") else 40
        elif pa.price_vs_ma20 == "below" and pa.price_vs_ma50 == "below":
            if is_sell:
                ma_score = 100
            elif is_buy:
                ma_score = 30
            else:
                ma_score = 50
        else:
            ma_score = 50
    c["ma_trend"] = ma_score

    # 4. Fundamentals (15%)
    f = 0
    if 0 < pa.pe < 15:
        f += 50
    elif pa.pe > 0:
        f += 20
    if 0.5 <= pa.beta <= 1.5:
        f += 25
    if pa.dividend_yield > 0:
        f += 25
    c["fundamental"] = min(f, 100)

    # 5. Flow (10%)
    fl = 50
    if is_buy:
        if pa.foreign_net_vol > 0:
            fl += 25
        elif pa.foreign_net_vol < 0:
            fl -= 25
        if pa.buy_sell_ratio > 1:
            fl += 25
        elif 0 < pa.buy_sell_ratio < 1:
            fl -= 25
    elif is_sell:
        if pa.foreign_net_vol < 0:
            fl += 25
        elif pa.foreign_net_vol > 0:
            fl -= 25
        if 0 < pa.buy_sell_ratio < 1:
            fl += 25
        elif pa.buy_sell_ratio > 1:
            fl -= 25
    c["flow"] = max(0, min(100, fl))

    # 6. News (10%)
    if is_buy:
        c["news"] = {"positive": 90, "mixed": 50}.get(news_sentiment, 30)
    elif is_sell:
        c["news"] = {"negative": 90, "mixed": 50}.get(news_sentiment, 30)
    else:
        c["news"] = 50

    weights = {
        "timeframe": 25,
        "rsi": 20,
        "ma_trend": 20,
        "fundamental": 15,
        "flow": 10,
        "news": 10,
    }
    total_w = sum(weights[k] for k in c if k in weights)
    weighted_sum = sum(c[k] * weights[k] for k in c if k in weights)
    final = round(weighted_sum / total_w) if total_w > 0 else 0
    return ConfidenceDetail(score=final, components=c)
