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


# ============================================================
# Vietnamese news title clustering (Jaccard similarity)
# ============================================================

_VI_STOPWORDS = {
    "đang", "được", "này", "đó", "các", "những", "một", "hai",
    "khi", "nếu", "nhưng", "hơn", "như", "cùng", "vẫn",
    "sau", "trước", "trong", "ngoài", "trên", "dưới", "giữa",
    "bởi", "đến", "theo", "với", "không", "chưa", "đã", "sẽ",
    "có", "phải", "nên", "cần", "đây", "tại", "lại", "lên",
    "xuống", "ra", "vào", "thêm", "bớt", "hôm", "nay", "mai",
    "qua", "năm", "tháng", "tuần",
}


def tokenize_vi_title(title: str) -> set:
    """Produce a set of significant tokens from a Vietnamese title.

    Lowercased, split on non-letter/digit, stopwords and tokens ≤2 chars dropped.
    """
    title = title.lower().strip()
    # Replace non-alphanumeric with space (Unicode-aware via isalpha/isdigit)
    cleaned = []
    for ch in title:
        if ch.isalpha() or ch.isdigit():
            cleaned.append(ch)
        else:
            cleaned.append(" ")
    tokens = "".join(cleaned).split()
    return {t for t in tokens if len(t) > 2 and t not in _VI_STOPWORDS}


def jaccard(a: set, b: set) -> float:
    """Jaccard similarity: |A∩B| / |A∪B|. Returns 0 for empty sets."""
    if not a or not b:
        return 0.0
    inter = len(a & b)
    union = len(a) + len(b) - inter
    return inter / union if union else 0.0


def cluster_titles(titles: list, threshold: float = 0.3) -> list:
    """Group title indices by Jaccard similarity.

    Titles with Jaccard ≥ threshold share a cluster. Greedy: first unclaimed
    title seeds a group, then consumes every later title above threshold.

    Returns list of index lists, e.g. [[0, 1], [2]] — cluster membership preserved
    in original order.

    threshold=0.3 is a sensible default (share ≥30% of significant tokens).
    """
    if not titles:
        return []
    token_sets = [tokenize_vi_title(t) for t in titles]

    cluster_of = [-1] * len(titles)
    num = 0
    for i in range(len(titles)):
        if cluster_of[i] != -1:
            continue
        cluster_of[i] = num
        num += 1
        for j in range(i + 1, len(titles)):
            if cluster_of[j] != -1:
                continue
            if jaccard(token_sets[i], token_sets[j]) >= threshold:
                cluster_of[j] = cluster_of[i]

    buckets: dict = {}
    for idx, c in enumerate(cluster_of):
        buckets.setdefault(c, []).append(idx)
    # Preserve cluster creation order (0, 1, 2, ...)
    return [buckets[c] for c in range(num) if c in buckets]
