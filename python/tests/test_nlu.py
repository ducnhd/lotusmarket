import os

os.environ["LOTUSMARKET_QUIET"] = "1"

from lotusmarket.nlu import Parser


def test_stock_info():
    p = Parser()
    r = p.parse("giá ACB hiện tại")
    assert r.intent == "GET_STOCK_INFO"
    assert r.tickers == ["ACB"]


def test_analyze_trend():
    p = Parser()
    r = p.parse("xu hướng FPT tăng hay giảm")
    assert r.intent == "ANALYZE_TREND"


def test_no_ticker():
    p = Parser()
    r = p.parse("thị trường hôm nay thế nào")
    assert r.tickers == []


def test_stop_words():
    p = Parser()
    r = p.parse("giá cổ phiếu")
    assert "GIA" not in r.tickers
