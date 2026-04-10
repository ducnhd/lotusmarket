"""Strategy backtesting engine for Vietnamese stocks."""

import math
from dataclasses import dataclass, field
from typing import List
import pandas as pd
from lotusmarket.technical import rsi_value, ma_value


@dataclass
class BacktestConfig:
    strategy: str = "rsi"  # "rsi", "ma_cross", "combined"
    capital: float = 100_000_000  # 100M VND
    rsi_buy: float = 30.0
    rsi_sell: float = 70.0


@dataclass
class BacktestTrade:
    date: str = ""
    action: str = ""
    price: float = 0.0
    shares: float = 0.0
    value: float = 0.0
    cash: float = 0.0
    pnl: float = 0.0


@dataclass
class BacktestResult:
    config: BacktestConfig = field(default_factory=BacktestConfig)
    trades: List[BacktestTrade] = field(default_factory=list)
    total_return: float = 0.0
    max_drawdown: float = 0.0
    win_rate: float = 0.0
    buy_hold_return: float = 0.0
    final_capital: float = 0.0
    trade_count: int = 0


def run(cfg: BacktestConfig, prices: pd.Series) -> BacktestResult:
    """Run a backtest strategy on price series.

    Args:
        cfg: Backtest configuration
        prices: pandas Series of close prices (at least 20 values)

    Returns:
        BacktestResult with trades, returns, drawdown, win rate
    """
    if len(prices) < 20:
        raise ValueError(f"Insufficient data: need 20+ prices, got {len(prices)}")

    if cfg.capital <= 0:
        cfg.capital = 100_000_000
    if cfg.rsi_buy <= 0:
        cfg.rsi_buy = 30
    if cfg.rsi_sell <= 0:
        cfg.rsi_sell = 70

    result = BacktestResult(config=cfg)
    closes = prices.values.tolist()

    cash = cfg.capital
    shares = 0.0
    trades = []
    peak_value = cfg.capital
    max_dd = 0.0
    wins = 0
    losses = 0

    start_idx = min(50, len(closes) - 1)
    if start_idx < 20:
        start_idx = 20

    for i in range(start_idx, len(closes)):
        price = closes[i]
        sub = pd.Series(closes[: i + 1])

        r = rsi_value(sub, 14)
        ma20 = ma_value(sub, 20)
        ma50 = ma_value(sub, 50)

        buy_signal = False
        sell_signal = False

        if cfg.strategy == "rsi":
            buy_signal = r < cfg.rsi_buy
            sell_signal = r > cfg.rsi_sell
        elif cfg.strategy == "ma_cross":
            if ma20 is not None and ma50 is not None:
                buy_signal = ma20 > ma50 and price > ma20
                sell_signal = ma20 < ma50 and price < ma20
        elif cfg.strategy == "combined":
            ma_buy = ma20 is not None and ma50 is not None and ma20 > ma50
            ma_sell = ma20 is not None and ma50 is not None and ma20 < ma50
            buy_signal = r < cfg.rsi_buy and ma_buy
            sell_signal = r > cfg.rsi_sell or ma_sell
        else:
            buy_signal = r < cfg.rsi_buy
            sell_signal = r > cfg.rsi_sell

        if buy_signal and shares == 0 and cash > 0:
            shares = math.floor(cash / price)
            if shares > 0:
                cost = shares * price
                cash -= cost
                trades.append(
                    BacktestTrade(
                        date=str(i),
                        action="BUY",
                        price=price,
                        shares=shares,
                        value=cost,
                        cash=cash,
                    )
                )
        elif sell_signal and shares > 0:
            value = shares * price
            pnl = value - (shares * trades[-1].price)
            if pnl > 0:
                wins += 1
            else:
                losses += 1
            cash += value
            trades.append(
                BacktestTrade(
                    date=str(i),
                    action="SELL",
                    price=price,
                    shares=shares,
                    value=value,
                    cash=cash,
                    pnl=pnl,
                )
            )
            shares = 0

        total_value = cash + shares * price
        if total_value > peak_value:
            peak_value = total_value
        dd = (peak_value - total_value) / peak_value * 100
        if dd > max_dd:
            max_dd = dd

    last_price = closes[-1]
    result.final_capital = cash + shares * last_price
    result.total_return = (result.final_capital - cfg.capital) / cfg.capital * 100
    result.max_drawdown = max_dd
    result.trades = trades
    result.trade_count = len(trades)

    total_trades = wins + losses
    if total_trades > 0:
        result.win_rate = wins / total_trades * 100

    first_price = closes[start_idx]
    result.buy_hold_return = (last_price - first_price) / first_price * 100

    return result
