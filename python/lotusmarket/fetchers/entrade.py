"""Entrade API — historical OHLCV data (no API key needed)."""

import json
import time as _time
from datetime import datetime, timedelta
from typing import List
from lotusmarket.types import StockData

ENTRADE_BASE_URL = "https://services.entrade.com.vn/chart-api/v2/ohlcs/stock"


def parse_entrade_response(data: bytes, ticker: str) -> List[StockData]:
    """Parse raw Entrade JSON response into StockData list (prices in VND)."""
    resp = json.loads(data)
    timestamps = resp.get("t", [])
    opens = resp.get("o", [])
    highs = resp.get("h", [])
    lows = resp.get("l", [])
    closes = resp.get("c", [])
    volumes = resp.get("v", [])
    n = len(timestamps)
    prices = []
    for i in range(n):
        dt = datetime.fromtimestamp(timestamps[i])
        sd = StockData(
            ticker=ticker,
            date=dt.strftime("%Y-%m-%d"),
            open=opens[i],
            high=highs[i],
            low=lows[i],
            close=closes[i],
            volume=int(volumes[i]),
        )
        sd.to_vnd()
        prices.append(sd)
    return prices


def entrade_history(ticker: str, days: int = 200) -> List[StockData]:
    """Fetch historical daily OHLCV data."""
    import httpx

    now = int(_time.time())
    from_ts = int((datetime.now() - timedelta(days=days)).timestamp())
    url = f"{ENTRADE_BASE_URL}?from={from_ts}&to={now}&symbol={ticker}&resolution=1D"
    resp = httpx.get(url, timeout=15)
    resp.raise_for_status()
    return parse_entrade_response(resp.content, ticker)


def entrade_latest(ticker: str) -> StockData:
    """Fetch most recent trading day's data."""
    prices = entrade_history(ticker, 5)
    if not prices:
        raise ValueError(f"No entrade data for {ticker}")
    return prices[-1]
