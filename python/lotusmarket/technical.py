"""Technical analysis indicators for Vietnamese stocks."""

import numpy as np
import pandas as pd
from dataclasses import dataclass
from typing import Optional


@dataclass
class DashboardMetrics:
    rsi: float = 0.0
    ma20: Optional[float] = None
    ma50: Optional[float] = None
    ma200: Optional[float] = None
    momentum: float = 0.0
    signal: str = "HOLD"
    score: float = 50.0


def ma(series: pd.Series, period: int) -> pd.Series:
    """Simple moving average."""
    return series.rolling(window=period, min_periods=period).mean()


def ma_value(prices: pd.Series, period: int) -> Optional[float]:
    """Latest MA value, or None if insufficient data."""
    if len(prices) < period:
        return None
    return float(prices.iloc[-period:].mean())


def rsi(series: pd.Series, period: int = 14) -> pd.Series:
    """RSI using Wilder's smoothing. Returns Series."""
    delta = series.diff()
    gain = delta.where(delta > 0, 0.0)
    loss = (-delta).where(delta < 0, 0.0)
    avg_gain = gain.ewm(alpha=1 / period, min_periods=period, adjust=False).mean()
    avg_loss = loss.ewm(alpha=1 / period, min_periods=period, adjust=False).mean()
    rs = avg_gain / avg_loss
    return 100 - (100 / (1 + rs))


def rsi_value(prices: pd.Series, period: int = 14) -> float:
    """Latest RSI value. Returns 50 if insufficient data."""
    if len(prices) < period + 1:
        return 50.0
    r = rsi(prices, period)
    val = r.iloc[-1]
    if pd.isna(val):
        return 50.0
    return float(val)


def momentum(prices: pd.Series, period: int = 20) -> float:
    """Rate of change over N periods as percentage."""
    if len(prices) <= period:
        return 0.0
    current = float(prices.iloc[-1])
    past = float(prices.iloc[-1 - period])
    if past == 0:
        return 0.0
    return ((current - past) / past) * 100


def signal(prices: pd.Series) -> str:
    """BUY/SELL/HOLD based on MA crossover + RSI zones."""
    ma20 = ma_value(prices, 20)
    ma50 = ma_value(prices, 50)
    if ma20 is None or ma50 is None:
        return "HOLD"
    price = float(prices.iloc[-1])
    r = rsi_value(prices, 14)
    if price > ma20 > ma50 and r < 70:
        return "BUY"
    if price < ma20 < ma50 and r > 30:
        return "SELL"
    return "HOLD"


def score(prices: pd.Series) -> float:
    """Composite 0-100 score."""
    sig = signal(prices)
    r = rsi_value(prices, 14)
    s = 50.0
    if sig == "BUY":
        s += 25
    elif sig == "SELL":
        s -= 25
    if r < 30:
        s += 10
    if r > 70:
        s -= 10
    return max(0, min(100, s))


def dashboard(prices: pd.Series) -> DashboardMetrics:
    """All-in-one technical assessment."""
    return DashboardMetrics(
        rsi=rsi_value(prices, 14),
        ma20=ma_value(prices, 20),
        ma50=ma_value(prices, 50),
        ma200=ma_value(prices, 200),
        momentum=momentum(prices, 20),
        signal=signal(prices),
        score=score(prices),
    )


def aggregate_weekly(df: pd.DataFrame) -> pd.DataFrame:
    """Convert daily OHLCV DataFrame to weekly candles.
    Input must have columns: date, open, high, low, close, volume.
    """
    if df.empty:
        return pd.DataFrame()
    d = df.copy()
    d["date"] = pd.to_datetime(d["date"])
    d = d.set_index("date").sort_index()
    weekly = (
        d.resample("W")
        .agg(
            {
                "open": "first",
                "high": "max",
                "low": "min",
                "close": "last",
                "volume": "sum",
            }
        )
        .dropna()
    )
    weekly = weekly.reset_index()
    weekly["date"] = weekly["date"].dt.strftime("%Y-%m-%d")
    return weekly
