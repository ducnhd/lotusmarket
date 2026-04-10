"""Vietnamese intent classification and ticker extraction."""

import re
from dataclasses import dataclass, field
from typing import List, Optional
from lotusmarket.types import VN30


@dataclass
class NLUResult:
    intent: str = ""
    confidence: float = 0.0
    tickers: List[str] = field(default_factory=list)
    timeframe: str = "1d"
    query: str = ""


_VN_STOP_WORDS = {
    "CAC",
    "CON",
    "CUA",
    "CHO",
    "DAU",
    "DAN",
    "GIA",
    "HAI",
    "HOA",
    "HOI",
    "LAM",
    "MOI",
    "NAM",
    "NEN",
    "SAU",
    "TAI",
    "THE",
    "TIN",
    "VAN",
    "VOI",
}

_TICKER_PATTERN = re.compile(r"\b([A-Z]{3})\b")

_INTENT_KEYWORDS = {
    "GET_STOCK_INFO": [
        "giá",
        "bao nhiêu",
        "hiện tại",
        "price",
        "current",
        "info",
        "thông tin",
    ],
    "ANALYZE_TREND": [
        "xu hướng",
        "tăng",
        "giảm",
        "chart",
        "trend",
        "biểu đồ",
        "phân tích",
    ],
    "GET_SENTIMENT": [
        "tin tức",
        "ý kiến",
        "news",
        "sentiment",
        "dư luận",
        "thị trường",
    ],
    "GET_RECOMMENDATION": [
        "nên",
        "mua",
        "bán",
        "buy",
        "sell",
        "signal",
        "khuyến nghị",
        "đầu tư",
    ],
    "COMPARE_STOCKS": ["so sánh", "so với", "compare", "vs", "và"],
}


class Parser:
    def __init__(self, tickers: Optional[List[str]] = None):
        self.ticker_set = set(tickers or VN30)

    def parse(self, query: str) -> NLUResult:
        lower = query.lower()
        result = NLUResult(query=query, tickers=self._extract_tickers(query))
        best_intent = "GET_STOCK_INFO"
        best_score = 0.0
        for intent, keywords in _INTENT_KEYWORDS.items():
            score = sum(1.0 for kw in keywords if kw in lower)
            if score > best_score:
                best_score = score
                best_intent = intent
        result.intent = best_intent
        result.confidence = min(best_score / 3.0, 1.0)
        result.timeframe = _extract_timeframe(lower)
        return result

    def _extract_tickers(self, query: str) -> List[str]:
        matches = _TICKER_PATTERN.findall(query.upper())
        seen = set()
        tickers = []
        for m in matches:
            if m in _VN_STOP_WORDS:
                continue
            if m in self.ticker_set and m not in seen:
                tickers.append(m)
                seen.add(m)
        return tickers


def _extract_timeframe(lower: str) -> str:
    if "tuần" in lower or "week" in lower:
        return "1w"
    if "tháng" in lower or "month" in lower:
        return "1m"
    if "năm" in lower or "year" in lower:
        return "1y"
    return "1d"
