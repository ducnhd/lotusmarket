from lotusmarket.ratings import compute


def test_bullish_path():
    closes = [100 + i * 1.5 for i in range(30)]
    volumes = [1000 + i * 50 for i in range(30)]
    ma20 = closes[-10]
    ma50 = closes[5]
    ma200 = closes[0]
    r = compute(closes, volumes, 75, 58, ma20, ma50, ma200)
    assert r.overall_verdict == "Outperform"
    assert r.trend_strength >= 4
    assert r.overall_gauge >= 60


def test_bearish_path():
    closes = [200 - i * 1.5 for i in range(30)]
    volumes = [1000] * 30
    ma20 = closes[0]
    ma50 = closes[0] + 5
    ma200 = closes[0] + 10
    r = compute(closes, volumes, 15, 28, ma20, ma50, ma200)
    assert r.overall_verdict != "Outperform"


def test_short_history_defaults():
    r = compute([100, 101, 102], None, 50, 50, None, None, None)
    assert r.price_strength == 3
    assert r.trend_strength == 3
    assert r.overall_verdict == "Neutral"
