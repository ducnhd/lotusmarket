from datetime import datetime
from lotusmarket.historical import (
    FeatureRow,
    analyze,
    to_markdown,
    rsi_bucket,
    ma_trend,
)


def _row(rsi, close, ma50, ma200, fwd60):
    return FeatureRow(
        ticker="ACB",
        date=datetime(2026, 1, 1),
        close=close,
        rsi14=rsi,
        ma50=ma50,
        ma200=ma200,
        fwd60=fwd60,
    )


def test_rsi_bucket_boundaries():
    assert rsi_bucket(20) == "RSI<30 (oversold)"
    assert rsi_bucket(45) == "RSI 30-50"
    assert rsi_bucket(60) == "RSI 50-70"
    assert rsi_bucket(80) == "RSI>70 (overbought)"


def test_ma_trend_logic():
    assert ma_trend(110, 105, 100) == "uptrend"
    assert ma_trend(90, 95, 100) == "downtrend"
    assert ma_trend(110, 95, 100) == "mixed"
    assert ma_trend(100, None, 95) == ""


def test_analyze_simple():
    rows = [
        _row(60, 110, 105, 100, fwd60=5.0),
        _row(60, 112, 106, 100, fwd60=6.0),
        _row(75, 120, 110, 100, fwd60=10.0),
        _row(20, 90, 95, 100, fwd60=-3.0),
    ]
    report = analyze(rows, "test")
    assert report.clean_rows == 4
    assert report.baseline.n == 4
    assert "rsi" in report.groups
    md = to_markdown(report)
    assert "Baseline" in md
    assert "RSI" in md


def test_outlier_clip():
    rows = [
        _row(60, 110, 105, 100, fwd60=5.0),
        _row(60, 110, 105, 100, fwd60=200.0),  # clipped — split anomaly
    ]
    report = analyze(rows, "test")
    # Only 1 valid baseline row after clipping
    assert report.baseline.n == 1
    assert report.baseline.mean_60d == 5.0
