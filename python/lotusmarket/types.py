"""Shared types and constants for lotusmarket."""

from dataclasses import dataclass


@dataclass
class StockData:
    ticker: str = ""
    date: str = ""
    open: float = 0.0
    high: float = 0.0
    low: float = 0.0
    close: float = 0.0
    volume: int = 0
    change_percent: float = 0.0
    change_value: float = 0.0
    ref_price: float = 0.0
    ceiling: float = 0.0
    floor: float = 0.0
    foreign_buy_vol: int = 0
    foreign_sell_vol: int = 0
    foreign_net_vol: int = 0
    bid: float = 0.0
    ask: float = 0.0
    bid_vol: int = 0
    ask_vol: int = 0

    def to_vnd(self):
        for attr in (
            "open",
            "high",
            "low",
            "close",
            "change_value",
            "ref_price",
            "ceiling",
            "floor",
            "bid",
            "ask",
        ):
            setattr(self, attr, getattr(self, attr) * 1000)


@dataclass
class KBSQuote:
    date: str = ""
    open: float = 0.0
    high: float = 0.0
    low: float = 0.0
    close: float = 0.0
    volume: float = 0.0
    market_cap: float = 0.0
    pe: float = 0.0
    pb: float = 0.0
    eps: float = 0.0
    bvps: float = 0.0
    beta: float = 0.0
    dividend_yield: float = 0.0


@dataclass
class InsiderTransaction:
    """Single insider trading record from CafeF."""

    ticker: str = ""
    insider_name: str = ""
    position: str = ""
    related_party: str = ""
    buy_registered: int = 0
    sell_registered: int = 0
    buy_result: int = 0
    sell_result: int = 0
    start_date: str = ""  # YYYY-MM-DD
    end_date: str = ""  # YYYY-MM-DD
    completion_date: str = ""  # YYYY-MM-DD or ""
    shares_after: int = 0
    ownership_pct: float = 0.0


TRADING_COMMISSION_RATE = 0.0015
SELLING_TAX_RATE = 0.001

VN30 = [
    "ACB",
    "BID",
    "CTG",
    "DGC",
    "FPT",
    "GAS",
    "GVR",
    "HDB",
    "HPG",
    "LPB",
    "MBB",
    "MSN",
    "MWG",
    "PLX",
    "SAB",
    "SHB",
    "SSB",
    "SSI",
    "STB",
    "TCB",
    "TPB",
    "VCB",
    "VHM",
    "VIB",
    "VIC",
    "VJC",
    "VNM",
    "VPB",
    "VPL",
    "VRE",
]

SECTORS = {
    "Ngân hàng": [
        "ACB",
        "BID",
        "CTG",
        "HDB",
        "LPB",
        "MBB",
        "SHB",
        "SSB",
        "STB",
        "TCB",
        "TPB",
        "VCB",
        "VIB",
        "VPB",
    ],
    "Tài chính": ["SSI"],
    "Công nghệ": ["FPT"],
    "Bán lẻ": ["MWG"],
    "Bất động sản": ["VHM", "VIC", "VRE"],
    "Năng lượng": ["GAS", "PLX"],
    "Tiêu dùng": ["MSN", "SAB", "VNM"],
    "Vật liệu": ["HPG"],
    "Hóa chất": ["DGC"],
    "Nông nghiệp": ["GVR"],
    "Vận tải": ["VJC"],
    "Du lịch": ["VPL"],
}

_SECTOR_INDEX = {t: s for s, tickers in SECTORS.items() for t in tickers}


def get_sector(ticker: str) -> str:
    return _SECTOR_INDEX.get(ticker, "Others")
