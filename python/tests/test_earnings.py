import os

os.environ["LOTUSMARKET_QUIET"] = "1"

from lotusmarket.earnings import (
    bctc_deadlines,
    detect_period,
    extract_annual_target,
    extract_profit,
    extract_revenue,
    extract_tickers,
    parse_vietnamese_number,
)


# ── parse_vietnamese_number ───────────────────────────────────────────────────


def test_parse_number_with_dots():
    assert parse_vietnamese_number("35.000") == 35000
    assert parse_vietnamese_number("4.500") == 4500
    assert parse_vietnamese_number("6.200") == 6200


def test_parse_number_plain():
    assert parse_vietnamese_number("100") == 100


def test_parse_number_empty():
    assert parse_vietnamese_number("") == 0


def test_parse_number_invalid():
    assert parse_vietnamese_number("abc") == 0


# ── extract_tickers ───────────────────────────────────────────────────────────


def test_extract_tickers_basic():
    assert extract_tickers("HPG báo lãi quý 3 đạt 4.500 tỷ") == ["HPG"]


def test_extract_tickers_skip_ceo():
    tickers = extract_tickers("CEO của VCB phát biểu")
    assert "CEO" not in tickers
    assert "VCB" in tickers


def test_extract_tickers_skip_gdp():
    tickers = extract_tickers("GDP tăng 7%, FPT tăng trưởng tốt")
    assert "GDP" not in tickers
    assert "FPT" in tickers


def test_extract_tickers_no_duplicates():
    tickers = extract_tickers("HPG tăng, HPG tiếp tục tăng")
    assert tickers.count("HPG") == 1


def test_extract_tickers_empty():
    assert extract_tickers("thị trường ổn định") == []


# ── detect_period ─────────────────────────────────────────────────────────────


def test_detect_period_q3_with_year():
    assert detect_period("HPG lãi ròng 6.200 tỷ Q3/2025") == "Q3/2025"


def test_detect_period_q1_viet():
    assert detect_period("Lợi nhuận quý 1/2026 đạt 2.000 tỷ") == "Q1/2026"


def test_detect_period_q4_year_in_title():
    assert detect_period("quý 4 năm 2025 đạt kỷ lục") == "Q4/2025"


def test_detect_period_9m():
    assert detect_period("9 tháng đầu năm 2025 lãi 20.000 tỷ") == "9M/2025"


def test_detect_period_6m():
    assert detect_period("6 tháng năm 2026 doanh thu 50.000 tỷ") == "6M/2026"


def test_detect_period_fy():
    assert detect_period("Kế hoạch lãi năm 2026 đạt 15.000 tỷ") == "FY/2026"


def test_detect_period_not_found():
    assert detect_period("không có kỳ nào") == ""


# ── extract_profit ────────────────────────────────────────────────────────────


def test_extract_profit_basic():
    profit, period = extract_profit("HPG lãi ròng 6.200 tỷ Q3/2025")
    assert profit == 6200
    assert period == "Q3/2025"


def test_extract_profit_not_found():
    profit, period = extract_profit("HPG thông báo kế hoạch sản xuất")
    assert profit == 0
    assert period == ""


def test_extract_profit_loi_nhuan():
    profit, _ = extract_profit("Lợi nhuận quý 3 đạt 4.500 tỷ đồng")
    assert profit == 4500


# ── extract_revenue ───────────────────────────────────────────────────────────


def test_extract_revenue_basic():
    assert extract_revenue("doanh thu đạt 35.000 tỷ trong quý 2") == 35000


def test_extract_revenue_not_found():
    assert extract_revenue("HPG tăng trưởng tốt") == 0


def test_extract_revenue_with_qualifier():
    assert extract_revenue("doanh thu vượt 10.000 tỷ") == 10000


# ── extract_annual_target ─────────────────────────────────────────────────────


def test_extract_annual_target_basic():
    assert extract_annual_target("Kế hoạch lãi 12.000 tỷ năm 2026") == 12000


def test_extract_annual_target_muc_tieu():
    assert extract_annual_target("Mục tiêu lợi nhuận đạt 8.000 tỷ năm nay") == 8000


def test_extract_annual_target_not_found():
    # Quarterly profit should not match target pattern
    assert extract_annual_target("HPG lãi ròng 6.200 tỷ Q3/2025") == 0


# ── bctc_deadlines ────────────────────────────────────────────────────────────


def test_bctc_deadlines_count():
    assert len(bctc_deadlines(2025)) == 6


def test_bctc_deadlines_dates_2025():
    dl = bctc_deadlines(2025)
    expected = [
        "2025-04-20",
        "2025-07-20",
        "2025-08-14",
        "2025-10-20",
        "2026-01-20",
        "2026-03-31",
    ]
    for i, e in enumerate(expected):
        assert dl[i].date == e, f"deadline[{i}].date = {dl[i].date}, want {e}"


def test_bctc_deadlines_descriptions_have_bctc():
    for d in bctc_deadlines(2025):
        assert d.description != ""
        assert "BCTC" in d.description


def test_bctc_deadlines_year_2026():
    dl = bctc_deadlines(2026)
    assert dl[0].date == "2026-04-20"
    assert dl[4].date == "2027-01-20"
    assert dl[5].date == "2027-03-31"
