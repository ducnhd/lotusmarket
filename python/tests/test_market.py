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
