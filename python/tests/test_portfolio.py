import os

os.environ["LOTUSMARKET_QUIET"] = "1"

from lotusmarket.portfolio import (
    twr,
    ReturnSegment,
    position_confidence,
    PositionAnalysis,
)


def test_twr():
    segments = [ReturnSegment(100, 110), ReturnSegment(110, 132)]
    assert abs(twr(segments) - 0.32) < 0.001


def test_twr_empty():
    assert twr([]) == 0.0


def test_twr_zero_start():
    assert twr([ReturnSegment(0, 100)]) == 0.0


def test_confidence_buy():
    pa = PositionAnalysis(
        action="BUY_MORE",
        timeframe_align="bullish",
        rsi=35,
        ma20=100.0,
        ma50=95.0,
        ma200=90.0,
        price_vs_ma20="above",
        price_vs_ma50="above",
        pe=10,
        beta=1.0,
        dividend_yield=3.0,
        foreign_net_vol=100000,
        buy_sell_ratio=1.5,
    )
    result = position_confidence(pa, "positive")
    assert result.score >= 60
