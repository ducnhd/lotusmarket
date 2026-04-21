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


# ============================================================
# News clustering tests
# ============================================================

from lotusmarket.sentiment import tokenize_vi_title, jaccard, cluster_titles


def test_tokenize_vi_title():
    tokens = tokenize_vi_title("VN-Index tăng 1% phiên cuối tuần")
    assert "index" in tokens
    assert "tăng" in tokens
    assert "phiên" in tokens
    assert "cuối" in tokens
    # Stopwords dropped
    assert "tuần" not in tokens
    # Short tokens dropped
    assert "vn" not in tokens
    assert "1" not in tokens


def test_jaccard():
    a = {"x", "y", "z"}
    b = {"x", "y", "w"}
    # intersection=2, union=4 → 0.5
    assert jaccard(a, b) == 0.5
    assert jaccard(set(), a) == 0


def test_cluster_titles_basic():
    titles = [
        "VN-Index tăng 1% phiên cuối tuần",
        "VNIndex đóng cửa tăng 1% phiên cuối",
        "Giá dầu WTI vượt 85 USD",
    ]
    groups = cluster_titles(titles, 0.3)
    assert len(groups) == 2
    assert len(groups[0]) == 2   # first two grouped
    assert groups[1] == [2]      # third alone


def test_cluster_titles_empty():
    assert cluster_titles([], 0.3) == []


def test_cluster_titles_all_distinct():
    titles = [
        "Ngân hàng tăng vốn điều lệ",
        "Dầu khí báo lãi quý 3",
        "Bất động sản đón sóng mới",
    ]
    groups = cluster_titles(titles, 0.3)
    assert len(groups) == 3
