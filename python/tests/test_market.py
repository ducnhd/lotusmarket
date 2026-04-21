import os

os.environ["LOTUSMARKET_QUIET"] = "1"

import pandas as pd
from lotusmarket.market import (
    score_breadth,
    score_foreign_flow,
    score_vix,
    score_yield_curve,
    sector_flow,
)


def test_score_breadth():
    assert score_breadth(70) == 80
    assert score_breadth(50) == 50
    assert score_breadth(30) == 20


def test_score_foreign_flow():
    assert score_foreign_flow(200e9) == 90
    assert score_foreign_flow(-200e9) == 20
    assert score_foreign_flow(0) == 50


def test_score_vix():
    assert score_vix(12) == 10
    assert score_vix(35) == 90


def test_score_yield_curve():
    assert score_yield_curve(1.5) == 10
    assert score_yield_curve(-0.5) == 90


def test_sector_flow_empty():
    result = sector_flow(pd.DataFrame())
    assert result.empty


def test_sector_flow_basic():
    df = pd.DataFrame(
        {
            "ticker": ["ACB", "VNM"],
            "close": [25000.0, 70000.0],
            "volume": [1000000, 500000],
            "change_percent": [1.5, -0.5],
            "foreign_net_vol": [500, -200],
        }
    )
    result = sector_flow(df)
    assert len(result) > 0
    assert result.iloc[0]["sector"] == "Ngân hàng"
    assert result.iloc[0]["rank"] == 1


# ============================================================
# Price Driver Attribution
# ============================================================

from lotusmarket.market import attribute_drivers, DriverFeatures


def test_attribute_drivers_flat_move():
    a = attribute_drivers(DriverFeatures(price_change_pct=0.02))
    assert a.dominant == "UNKNOWN"


def test_attribute_drivers_news_driven():
    a = attribute_drivers(
        DriverFeatures(
            price_change_pct=2.5,
            news_count=5,
            news_sentiment_avg=0.8,
            market_delta_pct=0.1,
            foreign_net_vnd=1e9,
        )
    )
    assert a.dominant == "NEWS"
    assert a.news_pct > 30


def test_attribute_drivers_foreign_flow():
    a = attribute_drivers(
        DriverFeatures(
            price_change_pct=1.5,
            foreign_net_vnd=50e9,
            market_delta_pct=0.2,
        )
    )
    assert a.dominant == "FOREIGN"


def test_attribute_drivers_misaligned():
    # Price up but news negative — NEWS must not dominate
    a = attribute_drivers(
        DriverFeatures(
            price_change_pct=1.0,
            news_count=3,
            news_sentiment_avg=-0.7,
            market_delta_pct=0.8,
        )
    )
    assert a.dominant != "NEWS"


def test_attribute_drivers_pcts_sum_sane():
    a = attribute_drivers(
        DriverFeatures(
            price_change_pct=1.5,
            foreign_net_vnd=20e9,
            market_delta_pct=0.5,
            sector_delta_pct=1.0,
            news_count=2,
            news_sentiment_avg=0.5,
            rsi=55,
            ma_signal="BUY",
            volume_ratio=1.8,
        )
    )
    total = (
        a.news_pct + a.foreign_pct + a.sector_pct + a.market_pct
        + a.technical_pct + a.volume_pct + a.residual_pct
    )
    assert abs(total - 100) < 1


def test_sector_flow_has_new_columns():
    import pandas as pd
    from lotusmarket.market import sector_flow
    df = pd.DataFrame([
        {"ticker": "ACB", "close": 25000, "volume": 1000000, "change_percent": 1.5, "foreign_net_vol": 500000},
        {"ticker": "VCB", "close": 95000, "volume": 500000, "change_percent": -0.5, "foreign_net_vol": -200000},
        {"ticker": "FPT", "close": 120000, "volume": 800000, "change_percent": 2.0, "foreign_net_vol": 300000},
    ])
    out = sector_flow(df)
    assert not out.empty
    for col in ("avg_change_pct", "total_volume_vnd", "advances_count", "total_count"):
        assert col in out.columns, f"missing column {col}"
