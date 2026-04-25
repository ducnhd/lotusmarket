"""CafeF insider trading fetcher.

Fetches insider transactions from the CafeF AJAX JSON API.
Only returns transactions with start_date within the past 90 days.

Usage::

    from lotusmarket.fetchers.cafef import cafef_insider
    rows = cafef_insider("HPG")
"""

from __future__ import annotations

import math
from datetime import datetime, timezone, timedelta
from typing import Any, Dict, List, Optional

from lotusmarket.types import InsiderTransaction

_CAFEF_URL = (
    "https://cafef.vn/du-lieu/Ajax/PageNew/DataHistory/GDCoDong.ashx"
    "?Symbol={ticker}&PageIndex=1&PageSize=50"
)
_VN_TZ = timezone(timedelta(hours=7))


def _parse_cafef_date(s: str) -> Optional[datetime]:
    """Parse CafeF's '/Date(ms)/' format to a datetime in VN timezone (UTC+7)."""
    if not s:
        return None
    s = s.removeprefix("/Date(").removesuffix(")/")
    if not s:
        return None
    try:
        ms = int(s)
    except ValueError:
        return None
    return datetime.fromtimestamp(ms / 1000, tz=_VN_TZ)


def _format_date(dt: Optional[datetime]) -> str:
    """Format datetime to YYYY-MM-DD, or '' if None."""
    if dt is None:
        return ""
    return dt.strftime("%Y-%m-%d")


def _clean_name(s: str) -> str:
    """Strip HTML anchor tags sometimes wrapped around insider names."""
    if ">" in s:
        s = s[s.index(">") + 1 :]
    if "<" in s:
        s = s[: s.index("<")]
    return s.strip()


def _parse_response(data: Dict[str, Any], ticker: str) -> List[InsiderTransaction]:
    """Parse a decoded CafeF JSON response dict into InsiderTransaction list."""
    if not data.get("Success"):
        raise ValueError(f"CafeF API returned success=false for {ticker}")

    cutoff = datetime.now(tz=_VN_TZ) - timedelta(days=90)
    rows: List[InsiderTransaction] = []

    for item in data.get("Data", {}).get("Data", []):
        start_dt = _parse_cafef_date(item.get("PlanBeginDate", ""))
        if start_dt is None or start_dt < cutoff:
            continue

        end_dt = _parse_cafef_date(item.get("PlanEndDate", ""))
        real_end_raw = item.get("RealEndDate")
        completion_dt = _parse_cafef_date(real_end_raw) if real_end_raw else None

        related_man = item.get("RelatedMan", "")
        related_pos = item.get("RelatedManPosition", "")
        related_party = f"{related_man} ({related_pos})" if related_man else ""

        ownership = item.get("TyLeSoHuu", 0.0)

        rows.append(
            InsiderTransaction(
                ticker=ticker,
                insider_name=_clean_name(item.get("TransactionMan", "")),
                position=item.get("TransactionManPosition", ""),
                related_party=related_party,
                buy_registered=item.get("PlanBuyVolume", 0),
                sell_registered=item.get("PlanSellVolume", 0),
                buy_result=item.get("RealBuyVolume", 0),
                sell_result=item.get("RealSellVolume", 0),
                start_date=_format_date(start_dt),
                end_date=_format_date(end_dt),
                completion_date=_format_date(completion_dt),
                shares_after=item.get("VolumeAfterTransaction", 0),
                ownership_pct=round(ownership * 100) / 100,
            )
        )
    return rows


def cafef_insider(ticker: str) -> List[InsiderTransaction]:
    """Fetch insider transactions for *ticker* from CafeF AJAX API.

    Requires the ``fetchers`` extra (``pip install lotusmarket[fetchers]``).
    Only returns transactions with start_date within the past 90 days.

    Args:
        ticker: VN stock ticker, e.g. "HPG".

    Returns:
        List of InsiderTransaction sorted by the API's default order (most recent first).

    Raises:
        ImportError: if httpx is not installed.
        httpx.HTTPStatusError: on non-200 response.
        ValueError: if the API returns success=false.
    """
    import httpx

    ticker = ticker.upper()
    url = _CAFEF_URL.format(ticker=ticker)
    resp = httpx.get(
        url,
        headers={"User-Agent": "Mozilla/5.0 (compatible; lotusmarket/1.0)"},
        timeout=30,
    )
    resp.raise_for_status()
    return _parse_response(resp.json(), ticker)
