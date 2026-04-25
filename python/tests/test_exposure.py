import os

os.environ["LOTUSMARKET_QUIET"] = "1"

from lotusmarket.exposure import (
    DriverClose,
    analyze,
    bounce_signals,
    format_for_prompt,
    get_mappings_for_ticker,
    regime_change_flags,
    unique_driver_symbols,
)


def _rising_provider(n=80):
    """Stub returning n closes that rise monotonically."""
    closes = [DriverClose(date="2025-01-01", close=100.0 + i * 0.5) for i in range(n)]

    def provider(symbol: str, days: int) -> list:
        return closes

    return provider


def _flat_provider(n=35):
    closes = [DriverClose(date="2025-01-01", close=100.0) for _ in range(n)]

    def provider(symbol: str, days: int) -> list:
        return closes

    return provider


def _empty_provider(symbol: str, days: int) -> list:
    return []


# ── get_mappings_for_ticker ───────────────────────────────────────────────────


def test_mappings_hpg():
    m = get_mappings_for_ticker("HPG")
    assert len(m) > 0
    assert m[0].driver == "000300.SS"


def test_mappings_unknown():
    assert get_mappings_for_ticker("UNKNOWN") == []


def test_mappings_banks():
    banks = [
        "VCB",
        "TCB",
        "BID",
        "CTG",
        "ACB",
        "MBB",
        "HDB",
        "LPB",
        "SHB",
        "SSB",
        "STB",
        "TPB",
        "VIB",
        "VPB",
    ]
    for b in banks:
        assert get_mappings_for_ticker(b), f"bank {b} has no mappings"


# ── unique_driver_symbols ─────────────────────────────────────────────────────


def test_unique_driver_symbols_count():
    syms = unique_driver_symbols()
    assert len(syms) >= 18, f"expected >=18 symbols, got {len(syms)}"


def test_unique_driver_symbols_sorted():
    syms = unique_driver_symbols()
    assert syms == sorted(syms)


def test_unique_driver_symbols_contains():
    syms = set(unique_driver_symbols())
    for expected in ["^VIX", "VNM", "EEM", "BHP", "TIO=F", "QQQ", "CL=F"]:
        assert expected in syms, f"{expected} not in unique_driver_symbols"


# ── analyze ───────────────────────────────────────────────────────────────────


def test_analyze_unmapped():
    a = analyze("UNKNOWN", _empty_provider)
    assert not a.has_mapping
    assert a.ticker == "UNKNOWN"


def test_analyze_hpg_positive_net_impact():
    a = analyze("HPG", _rising_provider())
    assert a.has_mapping
    assert a.net_impact > 0, f"expected positive net_impact, got {a.net_impact}"
    assert a.summary != ""


def test_analyze_hpg_sides():
    a = analyze("HPG", _rising_provider())
    assert len(a.input) > 0
    assert len(a.peer) > 0
    assert len(a.output) > 0


def test_analyze_no_history():
    a = analyze("HPG", _empty_provider)
    assert a.has_mapping
    # All drivers have no data → net_impact should be 0
    assert a.net_impact == 0.0


def test_analyze_summary_contains_ticker():
    a = analyze("ACB", _rising_provider())
    assert "ACB" in a.summary


# ── format_for_prompt ─────────────────────────────────────────────────────────


def test_format_for_prompt_empty():
    assert format_for_prompt([], _empty_provider) == ""


def test_format_for_prompt_no_data():
    out = format_for_prompt(["HPG"], _empty_provider)
    assert out == ""


def test_format_for_prompt_with_data():
    out = format_for_prompt(["HPG", "ACB"], _rising_provider())
    assert out != ""
    assert len(out) > 10


# ── regime_change_flags ───────────────────────────────────────────────────────


def test_regime_change_flags_flat():
    out = regime_change_flags(_flat_provider())
    assert out == ""


def test_regime_change_flags_big_move():
    # 60% rise over 35 days → >10% in 20d window
    closes = [DriverClose(date="2025-01-01", close=100.0 + i * 2.0) for i in range(35)]

    def hp(symbol, days):
        return closes

    out = regime_change_flags(hp)
    assert out != ""


# ── bounce_signals ────────────────────────────────────────────────────────────


def test_bounce_no_oversold():
    out = bounce_signals({"HPG": 55, "ACB": 60}, _rising_provider())
    assert out == ""


def test_bounce_oversold_positive_drivers():
    out = bounce_signals({"HPG": 32}, _rising_provider())
    assert out != ""


def test_bounce_empty_input():
    assert bounce_signals({}, _rising_provider()) == ""
