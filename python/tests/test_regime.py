import os

os.environ["LOTUSMARKET_QUIET"] = "1"

from lotusmarket.market import (
    CRISIS_FUNDAMENTAL,
    CRISIS_PANIC,
    REGIME_CRISIS,
    REGIME_EUPHORIA,
    REGIME_STABLE,
    REGIME_VOLATILE,
    classify_regime,
    identify_trigger,
    score_regime_signals,
)


# ── score_regime_signals ──────────────────────────────────────────────────────


def test_scores_vix_thresholds():
    assert score_regime_signals(31, 0, 0, 0, "", 0).vix_score == 80
    assert score_regime_signals(26, 0, 0, 0, "", 0).vix_score == 50
    assert score_regime_signals(21, 0, 0, 0, "", 0).vix_score == 25
    assert score_regime_signals(12, 0, 0, 0, "", 0).vix_score == 0


def test_scores_flow_thresholds():
    assert score_regime_signals(12, 0, -2500, 0, "", 0).flow_score == 80
    assert score_regime_signals(12, 0, -1500, 0, "", 0).flow_score == 50
    assert score_regime_signals(12, 0, -600, 0, "", 0).flow_score == 25
    assert score_regime_signals(12, 0, 0, 0, "", 0).flow_score == 0


def test_scores_news_tier():
    assert score_regime_signals(12, 0, 0, 1, "bắt giữ", 0).news_score == 100
    assert score_regime_signals(12, 0, 0, 2, "Fed tăng lãi suất", 0).news_score == 50
    assert score_regime_signals(12, 0, 0, 0, "", 0).news_score == 0


def test_scores_global_drop():
    assert score_regime_signals(12, 0, 0, 0, "", -5.5).global_score == 80
    assert score_regime_signals(12, 0, 0, 0, "", -3.5).global_score == 50
    assert score_regime_signals(12, 0, 0, 0, "", -1.0).global_score == 0


def test_scores_direction():
    assert score_regime_signals(12, -2.0, 0, 0, "", 0).direction == "down"
    assert score_regime_signals(12, +2.0, 0, 0, "", 0).direction == "up"
    assert score_regime_signals(12, 0.5, 0, 0, "", 0).direction == "mixed"
    assert score_regime_signals(12, -0.5, 0, 0, "", 0).direction == "mixed"


def test_composite_is_max():
    s = score_regime_signals(26, -2.5, 0, 0, "", 0)
    assert s.composite_score == max(
        s.vix_score, s.vnindex_score, s.flow_score, s.news_score, s.global_score
    )


# ── classify_regime ───────────────────────────────────────────────────────────


def test_stable_low_scores():
    s = score_regime_signals(12, 0.3, 100, 0, "", 0)
    regime, sub = classify_regime(s)
    assert regime == REGIME_STABLE
    assert sub == ""


def test_volatile_medium_vix():
    s = score_regime_signals(26, -1.5, -400, 0, "", -1.0)
    regime, sub = classify_regime(s)
    assert regime == REGIME_VOLATILE
    assert sub == ""


def test_volatile_moderate_vnindex():
    s = score_regime_signals(22, -2.5, 0, 0, "", 0)
    regime, sub = classify_regime(s)
    assert regime == REGIME_VOLATILE


def test_crisis_global_drop():
    s = score_regime_signals(35, -4.0, -3000, 0, "", -6.0)
    regime, sub = classify_regime(s)
    assert regime == REGIME_CRISIS
    assert sub == CRISIS_PANIC


def test_crisis_fundamental_bat_giu():
    s = score_regime_signals(
        32, -3.5, -2500, 1, "chủ tịch bị bắt giữ vì gian lận", -2.0
    )
    regime, sub = classify_regime(s)
    assert regime == REGIME_CRISIS
    assert sub == CRISIS_FUNDAMENTAL


def test_crisis_fundamental_vo_no():
    s = score_regime_signals(32, -3.5, -2500, 1, "tập đoàn tuyên bố vỡ nợ", -2.0)
    _, sub = classify_regime(s)
    assert sub == CRISIS_FUNDAMENTAL


def test_euphoria():
    s = score_regime_signals(35, +4.0, 0, 0, "", 0)
    regime, sub = classify_regime(s)
    assert regime == REGIME_EUPHORIA
    assert sub == ""


def test_crisis_panic_heavy_sell_no_news():
    s = score_regime_signals(28, -2.5, -3000, 0, "", -1.0)
    regime, sub = classify_regime(s)
    assert regime == REGIME_CRISIS
    assert sub == CRISIS_PANIC


# ── identify_trigger ──────────────────────────────────────────────────────────


def test_trigger_news_highest():
    # VNIndex score=80 (abs 3.5%), News score=100 (tier 1) → news wins
    s = score_regime_signals(12, -3.5, 0, 1, "FED tăng lãi suất đột ngột", 0)
    source, detail = identify_trigger(s)
    assert source == "news_tier1"
    assert detail == "FED tăng lãi suất đột ngột"


def test_trigger_vix_only():
    s = score_regime_signals(26, 0.3, 0, 0, "", 0)
    source, _ = identify_trigger(s)
    assert source == "vix"


def test_trigger_global():
    s = score_regime_signals(12, -1.5, 0, 0, "", -4.0)
    source, _ = identify_trigger(s)
    assert source == "global_contagion"
