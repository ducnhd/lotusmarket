"""Vietnamese financial keyword sentiment analysis."""

from dataclasses import dataclass


@dataclass
class SentimentResult:
    label: str  # "POSITIVE", "NEGATIVE", "NEUTRAL"
    score: float  # -1.0 to 1.0


_POSITIVE = {
    "tăng": 1.0,
    "tăng trưởng": 1.5,
    "tăng mạnh": 1.5,
    "tăng vọt": 2.0,
    "bùng nổ": 2.0,
    "đột phá": 1.5,
    "kỷ lục": 1.5,
    "bứt phá": 2.0,
    "phục hồi": 1.0,
    "hồi phục": 1.0,
    "khởi sắc": 1.5,
    "lạc quan": 1.5,
    "tích cực": 1.5,
    "thuận lợi": 1.0,
    "triển vọng": 1.0,
    "cơ hội": 1.0,
    "thăng hoa": 2.0,
    "lãi": 1.0,
    "lãi lớn": 1.5,
    "lợi nhuận": 1.0,
    "cổ tức": 1.0,
    "chia cổ tức": 1.5,
    "vượt kỳ vọng": 1.5,
    "mua ròng": 1.0,
    "dòng tiền": 0.5,
    "hấp dẫn": 1.0,
    "xanh": 1.0,
    "xanh rực": 1.5,
    "rực rỡ": 1.5,
    "lan tỏa": 1.0,
    "sôi động": 1.0,
    "nâng hạng": 1.5,
    "vượt đỉnh": 1.5,
    "tăng trần": 1.5,
}

_NEGATIVE = {
    "giảm": 1.0,
    "giảm mạnh": 1.5,
    "giảm sâu": 2.0,
    "giảm sốc": 2.0,
    "lao dốc": 2.0,
    "rớt": 1.5,
    "rớt mạnh": 2.0,
    "sụt": 1.5,
    "sụt giảm": 1.5,
    "sụp đổ": 2.0,
    "bán tháo": 2.0,
    "tháo chạy": 2.0,
    "hoảng loạn": 2.0,
    "lo ngại": 1.0,
    "bi quan": 1.5,
    "tiêu cực": 1.5,
    "rủi ro": 1.0,
    "lỗ": 1.0,
    "thua lỗ": 1.5,
    "phá sản": 2.0,
    "nợ xấu": 1.5,
    "đóng cửa": 1.0,
    "đình chỉ": 1.5,
    "cảnh báo": 1.0,
    "giảm sàn": 2.0,
    "đỏ sàn": 1.5,
    "bán ròng": 1.0,
    "rút vốn": 1.5,
}


def analyze(text: str) -> SentimentResult:
    """Analyze Vietnamese financial text sentiment.
    Returns SentimentResult with score in [-1.0, 1.0] and label.
    """
    if not text:
        return SentimentResult(label="NEUTRAL", score=0.0)
    lower = text.lower()
    pos = sum(w for phrase, w in _POSITIVE.items() if phrase in lower)
    neg = sum(w for phrase, w in _NEGATIVE.items() if phrase in lower)
    total = pos + neg
    if total == 0:
        return SentimentResult(label="NEUTRAL", score=0.0)
    s = max(-1.0, min(1.0, (pos - neg) / total))
    if s > 0.1:
        label = "POSITIVE"
    elif s < -0.1:
        label = "NEGATIVE"
    else:
        label = "NEUTRAL"
    return SentimentResult(label=label, score=s)
