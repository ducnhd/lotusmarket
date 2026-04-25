"""Earnings extraction — pure-logic parser for Vietnamese news headlines.

Extracts profit, revenue, annual targets, detects reporting periods, and
generates BCTC regulatory deadlines. No database coupling.

Usage::

    from lotusmarket.earnings import extract_profit, extract_revenue, bctc_deadlines

    profit, period = extract_profit("HPG lãi ròng 6.200 tỷ Q3/2025")
    # profit=6200, period="Q3/2025"

    deadlines = bctc_deadlines(2025)  # 6 Deadline objects
"""

from __future__ import annotations

import re
from dataclasses import dataclass
from typing import List, Optional, Tuple

# ---------------------------------------------------------------------------
# Types
# ---------------------------------------------------------------------------


@dataclass
class Deadline:
    """BCTC regulatory submission deadline."""

    date: str  # YYYY-MM-DD
    description: str  # Vietnamese description


# ---------------------------------------------------------------------------
# Constants
# ---------------------------------------------------------------------------

_SKIP_TICKERS = {
    "CEO",
    "CFO",
    "IPO",
    "GDP",
    "VND",
    "USD",
    "ETF",
    "EPS",
    "FDI",
    "FED",
    "THE",
    "AND",
}

_NUMBER_PAT = r"([\d.,]+)"
_QUAL_OPT = r"(?:hơn|gần|khoảng|đạt|vượt)?\s*"

_PROFIT_RE = re.compile(
    r"(?i)(?:lãi(?:\s+(?:ròng|sau\s+thuế|trước\s+thuế))?|lợi\s+nhuận|LNST)"
    r"[^.]{0,40}?" + _QUAL_OPT + _NUMBER_PAT + r"\s*tỷ"
)

_REVENUE_RE = re.compile(
    r"(?i)doanh\s+thu[^.]{0,40}?" + _QUAL_OPT + _NUMBER_PAT + r"\s*tỷ"
)

_TARGET_RE = re.compile(
    r"(?i)(?:kế\s+hoạch|mục\s+tiêu|chỉ\s+tiêu)[^.]{0,80}?"
    r"(?:lãi(?:\s+(?:ròng|sau\s+thuế|trước\s+thuế))?|lợi\s+nhuận)"
    r"[^.]{0,30}?" + _QUAL_OPT + _NUMBER_PAT + r"\s*tỷ"
)

_TICKER_RE = re.compile(r"\b([A-Z]{2,5})\b")
_YEAR_RE = re.compile(r"\b(20\d{2})\b")

_PERIOD_PATTERNS = {
    "Q1": re.compile(r"(?i)(?:quý\s*(?:1|I|một)|Q1)(?:\s*[/\-]\s*(\d{4}))?"),
    "Q2": re.compile(r"(?i)(?:quý\s*(?:2|II|hai)|Q2)(?:\s*[/\-]\s*(\d{4}))?"),
    "Q3": re.compile(r"(?i)(?:quý\s*(?:3|III|ba)|Q3)(?:\s*[/\-]\s*(\d{4}))?"),
    "Q4": re.compile(r"(?i)(?:quý\s*(?:4|IV|bốn)|Q4)(?:\s*[/\-]\s*(\d{4}))?"),
    "9M": re.compile(r"(?i)(?:9\s*tháng)(?:\s+(?:đầu\s+năm|năm)\s*(\d{4}))?"),
    "6M": re.compile(
        r"(?i)(?:6\s*tháng|nửa\s+đầu\s+năm|bán\s+niên)(?:\s*(?:năm)?\s*(\d{4}))?"
    ),
    "FY": re.compile(r"(?i)(?:cả\s+năm|cả\s+năm\s+(\d{4})|năm\s+(\d{4}))"),
}
_PERIOD_ORDER = ["Q1", "Q2", "Q3", "Q4", "9M", "6M", "FY"]


# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------


def parse_vietnamese_number(s: str) -> int:
    """Parse a Vietnamese-formatted number string.

    "35.000" → 35000  (dots are thousands separators in Vietnamese)
    """
    s = s.replace(".", "").replace(",", "").strip()
    try:
        return int(s)
    except ValueError:
        return 0


def detect_period(title: str) -> str:
    """Identify the reporting period from a Vietnamese headline.

    Returns strings like "Q3/2025", "9M/2026", "FY/2025", "Q4" (no year if absent).
    """
    for key in _PERIOD_ORDER:
        pat = _PERIOD_PATTERNS[key]
        m = pat.search(title)
        if m is None:
            continue
        year = next((g for g in m.groups() if g), None)
        if year is None:
            ym = _YEAR_RE.search(title)
            if ym:
                year = ym.group(1)
        if year:
            return f"{key}/{year}"
        return key
    return ""


def extract_tickers(title: str) -> List[str]:
    """Find 3-letter uppercase stock tickers in a Vietnamese headline.

    Filters out common non-ticker acronyms (CEO, GDP, etc.).
    """
    seen: set = set()
    result: List[str] = []
    for m in _TICKER_RE.finditer(title):
        t = m.group(1)
        if len(t) == 3 and t not in _SKIP_TICKERS and t not in seen and t.isupper():
            result.append(t)
            seen.add(t)
    return result


# ---------------------------------------------------------------------------
# Public extraction functions
# ---------------------------------------------------------------------------


def extract_profit(title: str) -> Tuple[int, str]:
    """Extract net profit (tỷ VND) and reporting period from a Vietnamese headline.

    Returns:
        (profit, period) — profit is 0 if not found, period is "" if not found.
    """
    m = _PROFIT_RE.search(title)
    if m is None:
        return 0, ""
    profit = parse_vietnamese_number(m.group(m.lastindex or 1))
    if profit == 0:
        return 0, ""
    return profit, detect_period(title)


def extract_revenue(title: str) -> int:
    """Extract revenue (tỷ VND) from a Vietnamese headline. Returns 0 if not found."""
    m = _REVENUE_RE.search(title)
    if m is None:
        return 0
    return parse_vietnamese_number(m.group(1))


def extract_annual_target(title: str) -> int:
    """Extract annual profit target (tỷ VND) from ĐHCĐ-style headlines. Returns 0 if not found."""
    m = _TARGET_RE.search(title)
    if m is None:
        return 0
    return parse_vietnamese_number(m.group(m.lastindex or 1))


def bctc_deadlines(year: int) -> List[Deadline]:
    """Return the 6 BCTC regulatory submission deadlines for *year*.

    Deadlines (Vietnamese regulations):
        - Q1:        April 20
        - Q2:        July 20
        - H1 audit:  August 14
        - Q3:        October 20
        - Q4 prelim: January 20 of year+1
        - FY audit:  March 31 of year+1
    """
    return [
        Deadline(f"{year}-04-20", f"Hạn nộp BCTC Q1/{year}"),
        Deadline(f"{year}-07-20", f"Hạn nộp BCTC Q2/{year}"),
        Deadline(f"{year}-08-14", f"Hạn nộp BCTC bán niên kiểm toán {year}"),
        Deadline(f"{year}-10-20", f"Hạn nộp BCTC Q3/{year}"),
        Deadline(f"{year + 1}-01-20", f"Hạn nộp BCTC Q4/{year}"),
        Deadline(f"{year + 1}-03-31", f"Hạn nộp BCTC năm {year} kiểm toán"),
    ]
