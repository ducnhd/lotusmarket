import os

os.environ["LOTUSMARKET_QUIET"] = "1"

from lotusmarket.sentiment import analyze


def test_positive():
    r = analyze("Ngân hàng tăng mạnh, khối ngoại mua ròng")
    assert r.label == "POSITIVE"
    assert r.score > 0


def test_negative():
    r = analyze("Thị trường giảm sâu, bán tháo hoảng loạn")
    assert r.label == "NEGATIVE"
    assert r.score < 0


def test_neutral():
    r = analyze("Hôm nay thời tiết đẹp")
    assert r.label == "NEUTRAL"
    assert abs(r.score) < 0.01


def test_empty():
    r = analyze("")
    assert r.label == "NEUTRAL"


def test_bounded():
    r = analyze("tăng vọt bùng nổ đột phá kỷ lục bứt phá thăng hoa tăng trưởng lãi lớn")
    assert -1.0 <= r.score <= 1.0
