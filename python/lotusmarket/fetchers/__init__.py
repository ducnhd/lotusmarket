"""Data fetchers for Vietnamese stock market APIs."""

from lotusmarket.fetchers.vps import vps, vps_multiple, parse_vps_response
from lotusmarket.fetchers.entrade import (
    entrade_history,
    entrade_latest,
    parse_entrade_response,
)
from lotusmarket.fetchers.kbs import kbs


def stock_with_fallback(ticker: str):
    """Try VPS first, fall back to Entrade."""
    try:
        return vps(ticker)
    except Exception:
        return entrade_latest(ticker)
