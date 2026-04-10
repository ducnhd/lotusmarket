import os

os.environ["LOTUSMARKET_QUIET"] = "1"

import pandas as pd
from lotusmarket.backtest import run, BacktestConfig


def make_oscillating_prices(n=200, base=25000):
    """Create oscillating prices that trigger RSI signals."""
    return pd.Series(
        [base + (i % 20) * 100 - 1000 for i in range(n)], dtype=float
    ).clip(lower=1000)


def test_insufficient_data():
    try:
        run(BacktestConfig(), pd.Series([100.0]))
        assert False, "should raise"
    except ValueError:
        pass


def test_rsi_strategy():
    prices = make_oscillating_prices(200)
    result = run(BacktestConfig(strategy="rsi"), prices)
    assert result.final_capital > 0
    assert 0 <= result.max_drawdown <= 100


def test_ma_cross_strategy():
    prices = pd.Series([20000.0 + i * 50 + (i % 5) * 100 for i in range(100)])
    result = run(BacktestConfig(strategy="ma_cross"), prices)
    assert result.config.strategy == "ma_cross"


def test_default_capital():
    prices = make_oscillating_prices(100)
    result = run(BacktestConfig(strategy="rsi"), prices)
    assert result.config.capital == 100_000_000


def test_buy_hold_positive():
    prices = pd.Series([20000.0 + i * 100 for i in range(100)])
    result = run(BacktestConfig(strategy="rsi"), prices)
    assert result.buy_hold_return > 0


def test_win_rate_bounded():
    prices = make_oscillating_prices(200)
    result = run(BacktestConfig(strategy="rsi"), prices)
    assert 0 <= result.win_rate <= 100
