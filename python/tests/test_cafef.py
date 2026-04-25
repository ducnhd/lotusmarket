import os

os.environ["LOTUSMARKET_QUIET"] = "1"

import json
import time
from datetime import datetime, timezone, timedelta

import pytest

from lotusmarket.fetchers.cafef import _parse_response, _parse_cafef_date, _clean_name

_VN_TZ = timezone(timedelta(hours=7))


def _now_ms() -> int:
    return int(time.time() * 1000)


def _future_ms(days: int = 30) -> int:
    return int((time.time() + days * 86400) * 1000)


def _ago_ms(days: int) -> int:
    return int((time.time() - days * 86400) * 1000)


def _make_response(begin_ms: int, end_ms: int) -> dict:
    return {
        "Success": True,
        "Data": {
            "TotalCount": 2,
            "Data": [
                {
                    "Stock": "HPG",
                    "TransactionMan": "<a href='#'>Nguyễn Văn A</a>",
                    "TransactionManPosition": "Chủ tịch HĐQT",
                    "RelatedMan": "",
                    "RelatedManPosition": "",
                    "VolumeBeforeTransaction": 1_000_000,
                    "PlanBuyVolume": 500_000,
                    "PlanSellVolume": 0,
                    "PlanBeginDate": f"/Date({begin_ms})/",
                    "PlanEndDate": f"/Date({end_ms})/",
                    "RealBuyVolume": 480_000,
                    "RealSellVolume": 0,
                    "RealEndDate": None,
                    "VolumeAfterTransaction": 1_480_000,
                    "TyLeSoHuu": 2.56789,
                },
                {
                    "Stock": "HPG",
                    "TransactionMan": "Trần Thị B",
                    "TransactionManPosition": "Phó TGĐ",
                    "RelatedMan": "Nguyễn Văn C",
                    "RelatedManPosition": "Anh ruột",
                    "VolumeBeforeTransaction": 200_000,
                    "PlanBuyVolume": 0,
                    "PlanSellVolume": 100_000,
                    "PlanBeginDate": f"/Date({begin_ms})/",
                    "PlanEndDate": f"/Date({end_ms})/",
                    "RealBuyVolume": 0,
                    "RealSellVolume": 95_000,
                    "RealEndDate": f"/Date({end_ms})/",
                    "VolumeAfterTransaction": 105_000,
                    "TyLeSoHuu": 0.12345,
                },
            ],
        },
    }


# ── _parse_cafef_date ─────────────────────────────────────────────────────────


def test_parse_cafef_date_valid():
    ms = 1705276800000  # 2024-01-15 00:00 UTC → 2024-01-15 07:00 VN
    dt = _parse_cafef_date(f"/Date({ms})/")
    assert dt is not None
    assert dt.year == 2024


def test_parse_cafef_date_empty():
    assert _parse_cafef_date("") is None


def test_parse_cafef_date_bad():
    assert _parse_cafef_date("/Date(abc)/") is None


# ── _clean_name ───────────────────────────────────────────────────────────────


def test_clean_name_with_html():
    assert _clean_name("<a href='#'>Nguyễn Văn A</a>") == "Nguyễn Văn A"


def test_clean_name_plain():
    assert _clean_name("  Trần Thị B  ") == "Trần Thị B"


def test_clean_name_no_tag():
    assert _clean_name("Lê Văn C") == "Lê Văn C"


# ── _parse_response ───────────────────────────────────────────────────────────


def test_parse_response_basic():
    data = _make_response(_now_ms(), _future_ms())
    rows = _parse_response(data, "HPG")
    assert len(rows) == 2


def test_parse_response_fields():
    data = _make_response(_now_ms(), _future_ms())
    rows = _parse_response(data, "HPG")
    r0 = rows[0]
    assert r0.ticker == "HPG"
    assert r0.insider_name == "Nguyễn Văn A"
    assert r0.buy_registered == 500_000
    assert r0.buy_result == 480_000
    assert r0.shares_after == 1_480_000
    assert r0.ownership_pct == 2.57
    assert r0.completion_date == ""  # RealEndDate was null


def test_parse_response_related_party():
    data = _make_response(_now_ms(), _future_ms())
    rows = _parse_response(data, "HPG")
    r1 = rows[1]
    assert "Nguyễn Văn C" in r1.related_party


def test_parse_response_completion_date_set():
    data = _make_response(_now_ms(), _future_ms())
    rows = _parse_response(data, "HPG")
    r1 = rows[1]
    assert r1.completion_date != ""  # RealEndDate was provided


def test_parse_response_90day_cutoff():
    # Both transactions started 200 days ago → filtered
    data = _make_response(_ago_ms(200), _ago_ms(180))
    rows = _parse_response(data, "HPG")
    assert len(rows) == 0


def test_parse_response_success_false():
    data = {"Success": False, "Data": {"TotalCount": 0, "Data": []}}
    with pytest.raises(ValueError):
        _parse_response(data, "XYZ")


def test_parse_response_ownership_rounding():
    data = _make_response(_now_ms(), _future_ms())
    rows = _parse_response(data, "HPG")
    # TyLeSoHuu=2.56789 → round(2.56789*100)/100 = 2.57
    assert rows[0].ownership_pct == 2.57
    # TyLeSoHuu=0.12345 → round(0.12345*100)/100 = 0.12
    assert rows[1].ownership_pct == 0.12


# ── live network skipped ──────────────────────────────────────────────────────


def test_cafef_insider_import():
    """cafef_insider function must be importable without hitting the network."""
    from lotusmarket.fetchers.cafef import cafef_insider  # noqa: F401
