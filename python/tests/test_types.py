import os

os.environ["LOTUSMARKET_QUIET"] = "1"

from lotusmarket.types import StockData, get_sector


def test_to_vnd():
    s = StockData(
        open=26.9,
        high=27.5,
        low=26.0,
        close=27.2,
        bid=27.1,
        ask=27.3,
        ref_price=26.9,
        ceiling=28.7,
        floor=25.1,
        change_value=0.3,
    )
    s.to_vnd()
    assert s.close == 27200
    assert s.open == 26900
    assert s.bid == 27100


def test_get_sector():
    assert get_sector("ACB") == "Ngân hàng"
    assert get_sector("FPT") == "Công nghệ"
    assert get_sector("UNKNOWN") == "Others"
