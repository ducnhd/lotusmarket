import os

os.environ["LOTUSMARKET_QUIET"] = "1"

import pandas as pd
from lotusmarket.signals import classify_volume_surge, volume_surge_label, domestic_flow


def test_no_surge():
    s = classify_volume_surge(1000, 1000, 0, 1.0)
    assert s.signal == ""


def test_accumulation():
    s = classify_volume_surge(3000, 1000, 0, 2.0)
    assert s.signal == "accumulation"
    assert s.intensity == "strong_surge"


def test_distribution():
    s = classify_volume_surge(2500, 1000, 0, -1.5)
    assert s.signal == "distribution"


def test_zscore():
    s = classify_volume_surge(5000, 1000, 500, 1.0)
    assert s.intensity == "strong_surge"


def test_labels():
    assert volume_surge_label("accumulation", "strong_surge") == "Tiền đổ vào mạnh"
    assert volume_surge_label("accumulation", "surge") == "Tiền đang đổ vào"
    assert volume_surge_label("distribution", "strong_surge") == "Bán tháo mạnh"
    assert volume_surge_label("distribution", "surge") == "Đang bán tháo"
    assert volume_surge_label("high_activity", "surge") == "Giao dịch đột biến"
    assert volume_surge_label("", "") == ""


def test_domestic_flow_buy():
    df = pd.DataFrame(
        {
            "ticker": ["ACB", "VNM"],
            "volume": [1000, 2000],
            "bid_vol": [700, 1400],
            "ask_vol": [300, 600],
            "foreign_buy_vol": [50, 100],
            "foreign_sell_vol": [50, 100],
        }
    )
    flow = domestic_flow(df)
    assert flow.signal == "mua mạnh"
    assert flow.total_volume == 3000


def test_domestic_flow_sell():
    df = pd.DataFrame(
        {
            "ticker": ["ACB"],
            "volume": [1000],
            "bid_vol": [200],
            "ask_vol": [800],
            "foreign_buy_vol": [50],
            "foreign_sell_vol": [50],
        }
    )
    flow = domestic_flow(df)
    assert flow.signal == "bán mạnh"
