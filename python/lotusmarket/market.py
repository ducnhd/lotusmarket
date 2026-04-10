"""Market pulse scoring, sector flow ranking, risk indicators."""

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
    Returns: DataFrame with sector, rank, signal, net_foreign, breadth, score columns.
    """
    if df.empty:
        return pd.DataFrame()
    df2 = df.copy()
    df2["sector"] = df2["ticker"].map(get_sector)
    df2 = df2[df2["sector"] != "Others"]
    if df2.empty:
        return pd.DataFrame()
    df2["flow_value"] = df2["foreign_net_vol"] * df2["close"]
    grouped = (
        df2.groupby("sector")
        .agg(
            total_flow=("flow_value", "sum"),
            advances=("change_percent", lambda x: (x > 0).sum()),
            total=("ticker", "count"),
            total_volume=("volume", "sum"),
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
        ["sector", "rank", "signal", "total_flow", "breadth", "score"]
    ].rename(columns={"total_flow": "net_foreign"})
