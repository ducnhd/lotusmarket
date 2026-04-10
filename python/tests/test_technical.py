import os

os.environ["LOTUSMARKET_QUIET"] = "1"

import pandas as pd
import numpy as np
from lotusmarket.technical import (
    ma,
    ma_value,
    rsi_value,
    momentum,
    signal,
    score,
    dashboard,
    aggregate_weekly,
)


def test_ma_basic():
    prices = pd.Series(range(1, 21), dtype=float)
    val = ma_value(prices, 20)
    assert abs(val - 10.5) < 0.01


def test_ma_insufficient():
    val = ma_value(pd.Series([1.0, 2.0, 3.0]), 20)
    assert val is None


def test_rsi_uptrend():
    prices = pd.Series([100.0 + i for i in range(30)])
    r = rsi_value(prices, 14)
    assert r > 50


def test_rsi_insufficient():
    r = rsi_value(pd.Series([1.0, 2.0, 3.0]), 14)
    assert r == 50.0


def test_momentum():
    prices = pd.Series([100.0] + [100.0 + i for i in range(1, 11)])
    m = momentum(prices, 10)
    assert abs(m - 10.0) < 0.01


def test_signal_buy():
    # Create uptrend with oscillation so RSI stays below 70:
    # base uptrend + zigzag noise keeps price > MA20 > MA50 while RSI < 70
    np.random.seed(42)
    base = [50.0 + i * 0.5 for i in range(60)]
    noise = [2.0 * ((-1) ** i) for i in range(60)]
    prices = pd.Series([b + n for b, n in zip(base, noise)])
    sig = signal(prices)
    assert sig == "BUY"


def test_signal_hold_insufficient():
    sig = signal(pd.Series([1.0, 2.0, 3.0]))
    assert sig == "HOLD"


def test_score_range():
    prices = pd.Series([float(i) for i in range(50, 110)])
    s = score(prices)
    assert 0 <= s <= 100


def test_dashboard():
    prices = pd.Series([20000.0 + i * 10 for i in range(200)])
    d = dashboard(prices)
    assert d.rsi > 0
    assert d.ma20 is not None
    assert d.signal in ("BUY", "SELL", "HOLD")


def test_aggregate_weekly_empty():
    result = aggregate_weekly(pd.DataFrame())
    assert result.empty


def test_aggregate_weekly_basic():
    daily = pd.DataFrame(
        {
            "date": ["2025-03-03", "2025-03-04", "2025-03-05"],
            "open": [100.0, 105.0, 108.0],
            "high": [110.0, 115.0, 120.0],
            "low": [90.0, 95.0, 100.0],
            "close": [105.0, 108.0, 112.0],
            "volume": [1000, 1500, 2000],
        }
    )
    result = aggregate_weekly(daily)
    assert len(result) == 1
    assert result.iloc[0]["open"] == 100.0
    assert result.iloc[0]["high"] == 120.0
    assert result.iloc[0]["low"] == 90.0
    assert result.iloc[0]["close"] == 112.0
    assert result.iloc[0]["volume"] == 4500
