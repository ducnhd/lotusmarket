"""KBS API — fundamental data (no API key needed)."""

from datetime import datetime, timedelta
from lotusmarket.types import KBSQuote

KBS_BASE_URL = "https://kbbuddywts.kbsec.com.vn/sas/kbsv-stock-data-store/stock"


def kbs(ticker: str) -> KBSQuote:
    """Fetch latest fundamentals (P/E, P/B, EPS, Beta, Dividend Yield)."""
    import httpx

    to_date = datetime.now().strftime("%Y-%m-%d")
    from_date = (datetime.now() - timedelta(days=7)).strftime("%Y-%m-%d")
    url = f"{KBS_BASE_URL}/{ticker}/historical-quotes?from={from_date}&to={to_date}"
    resp = httpx.get(url, timeout=15)
    resp.raise_for_status()
    quotes = resp.json()
    if not quotes:
        raise ValueError(f"No KBS data for {ticker}")
    q = quotes[0]
    return KBSQuote(
        date=q.get("TradingDate", ""),
        open=q.get("OpenPrice", 0),
        high=q.get("HighestPrice", 0),
        low=q.get("LowestPrice", 0),
        close=q.get("ClosePrice", 0),
        volume=q.get("TotalVol", 0),
        market_cap=q.get("MarketCapital", 0),
        pe=q.get("PE", 0),
        pb=q.get("PB", 0),
        eps=q.get("EPS", 0),
        bvps=q.get("BVPS", 0),
        beta=q.get("Beta", 0),
        dividend_yield=q.get("Yield", 0),
    )
