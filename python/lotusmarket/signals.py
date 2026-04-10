"""Volume surge detection and domestic flow analysis."""

from dataclasses import dataclass, field
from typing import List, Optional
import pandas as pd


@dataclass
class VolumeSurge:
    ticker: str = ""
    current_vol: int = 0
    avg_vol_20: float = 0.0
    stddev_vol_20: float = 0.0
    ratio: float = 0.0
    z_score: float = 0.0
    price_change: float = 0.0
    signal: str = ""  # "accumulation", "distribution", "high_activity", ""
    intensity: str = ""  # "surge", "strong_surge", ""
    label: str = ""  # Vietnamese label


@dataclass
class DomesticFlow:
    total_volume: int = 0
    domestic_volume: int = 0
    foreign_volume: int = 0
    buy_pressure: float = 50.0
    signal: str = ""
    top_buyers: List[str] = field(default_factory=list)
    top_sellers: List[str] = field(default_factory=list)


def volume_surge_label(signal: str, intensity: str) -> str:
    """Vietnamese label for a volume surge."""
    labels = {
        ("accumulation", "strong_surge"): "Tiền đổ vào mạnh",
        ("accumulation", "surge"): "Tiền đang đổ vào",
        ("distribution", "strong_surge"): "Bán tháo mạnh",
        ("distribution", "surge"): "Đang bán tháo",
    }
    if signal == "high_activity":
        return "Giao dịch đột biến"
    return labels.get((signal, intensity), "")


def classify_volume_surge(
    current_vol: int,
    avg_vol_20: float,
    stddev_vol_20: float = 0.0,
    price_change: float = 0.0,
) -> VolumeSurge:
    """Classify a volume anomaly using z-score or ratio fallback."""
    s = VolumeSurge(
        current_vol=current_vol,
        avg_vol_20=avg_vol_20,
        stddev_vol_20=stddev_vol_20,
        price_change=price_change,
    )
    if avg_vol_20 <= 0:
        return s
    s.ratio = current_vol / avg_vol_20
    if stddev_vol_20 > 0:
        s.z_score = (current_vol - avg_vol_20) / stddev_vol_20
        if s.z_score < 2.0:
            return s
        s.intensity = "strong_surge" if s.z_score >= 3.0 else "surge"
    else:
        if s.ratio < 2.0:
            return s
        s.intensity = "strong_surge" if s.ratio >= 3.0 else "surge"
    if price_change >= 0.5:
        s.signal = "accumulation"
    elif price_change <= -0.5:
        s.signal = "distribution"
    else:
        s.signal = "high_activity"
    s.label = volume_surge_label(s.signal, s.intensity)
    return s


def _classify_domestic_signal(buy_pressure: float) -> str:
    if buy_pressure >= 60:
        return "mua mạnh"
    elif buy_pressure >= 53:
        return "nghiêng mua"
    elif buy_pressure >= 47:
        return "cân bằng"
    elif buy_pressure >= 40:
        return "nghiêng bán"
    else:
        return "bán mạnh"


def domestic_flow(df: pd.DataFrame) -> DomesticFlow:
    """Compute domestic trading metrics from a DataFrame with bid_vol, ask_vol, volume, foreign_buy_vol, foreign_sell_vol columns."""
    if df.empty:
        return DomesticFlow()
    total_vol = int(df["volume"].sum())
    foreign_vol = (
        int((df.get("foreign_buy_vol", 0) + df.get("foreign_sell_vol", 0)).sum())
        if "foreign_buy_vol" in df.columns
        else 0
    )
    total_bid = int(df["bid_vol"].sum()) if "bid_vol" in df.columns else 0
    total_ask = int(df["ask_vol"].sum()) if "ask_vol" in df.columns else 0
    buy_pressure = 50.0
    if total_bid + total_ask > 0:
        buy_pressure = total_bid / (total_bid + total_ask) * 100
    result = DomesticFlow(
        total_volume=total_vol,
        domestic_volume=total_vol - foreign_vol,
        foreign_volume=foreign_vol,
        buy_pressure=buy_pressure,
        signal=_classify_domestic_signal(buy_pressure),
    )
    if "bid_vol" in df.columns and "ask_vol" in df.columns and "ticker" in df.columns:
        df2 = df.copy()
        df2["_ratio"] = (
            df2["bid_vol"] / (df2["bid_vol"] + df2["ask_vol"]).replace(0, 1) * 100
        )
        sorted_df = df2.sort_values("_ratio", ascending=False)
        result.top_buyers = sorted_df.head(3)["ticker"].tolist()
        result.top_sellers = sorted_df.tail(3)["ticker"].tolist()
    return result
