"""VPS DataFeed API — real-time quotes (no API key needed)."""

import json
from datetime import date
from typing import List
from lotusmarket.types import StockData

VPS_BASE_URL = "https://bgapidatafeed.vps.com.vn/getliststockdata"


def _parse_pipe_vol(s: str) -> int:
    parts = s.split("|")
    if len(parts) >= 2:
        try:
            return int(parts[1])
        except ValueError:
            return 0
    return 0


def parse_vps_response(data: bytes) -> List[StockData]:
    """Parse raw VPS JSON response into StockData list (prices in VND)."""
    raw = json.loads(data)
    stocks = []
    for item in raw:
        ticker = str(item.get("sym", ""))
        if not ticker:
            continue
        ref_price = float(item.get("r", 0))
        last_price = float(item.get("lastPrice", 0))
        open_price = float(item.get("openPrice", 0))
        if last_price == 0 and ref_price > 0:
            last_price = ref_price
        change_value = 0.0
        change_pct = 0.0
        if ref_price > 0:
            change_value = last_price - ref_price
            change_pct = (change_value / ref_price) * 100
        f_buy = int(float(item.get("fBVol", 0)))
        f_sell = int(float(item.get("fSVolume", 0)))
        g1 = str(item.get("g1", ""))
        g2 = str(item.get("g2", ""))
        s = StockData(
            ticker=ticker,
            date=date.today().isoformat(),
            open=open_price,
            high=float(item.get("highPrice", 0)),
            low=float(item.get("lowPrice", 0)),
            close=last_price,
            volume=int(float(item.get("lot", 0))),
            change_percent=change_pct,
            change_value=change_value,
            ref_price=ref_price,
            ceiling=float(item.get("c", 0)),
            floor=float(item.get("f", 0)),
            foreign_buy_vol=f_buy,
            foreign_sell_vol=f_sell,
            foreign_net_vol=f_buy - f_sell,
            bid=float(g1.split("|")[0]) if "|" in g1 else 0.0,
            ask=float(g2.split("|")[0]) if "|" in g2 else 0.0,
            bid_vol=_parse_pipe_vol(g1),
            ask_vol=_parse_pipe_vol(g2),
        )
        s.to_vnd()
        stocks.append(s)
    return stocks


def vps(ticker: str) -> StockData:
    """Fetch real-time quote for a single ticker."""
    stocks = vps_multiple([ticker])
    if not stocks:
        raise ValueError(f"No data for {ticker}")
    return stocks[0]


def vps_multiple(tickers: List[str]) -> List[StockData]:
    """Fetch real-time quotes for multiple tickers. Retries 3 times."""
    import httpx

    url = f"{VPS_BASE_URL}/{','.join(tickers)}"
    last_err = None
    for attempt in range(3):
        try:
            resp = httpx.get(url, timeout=15)
            resp.raise_for_status()
            return parse_vps_response(resp.content)
        except Exception as e:
            last_err = e
    raise ConnectionError(f"VPS API failed after 3 attempts: {last_err}")
