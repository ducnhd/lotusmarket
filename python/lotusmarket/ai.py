"""AI-powered stock analysis using Claude API.

Requires CLAUDE_API_KEY. If not configured, functions raise MissingAPIKeyError
with setup instructions.

Default model: claude-sonnet-4-6. Configurable via AIConfig.
"""

from dataclasses import dataclass, field
from typing import Optional, List
import os

import pandas as pd

from lotusmarket.technical import dashboard as compute_dashboard


class MissingAPIKeyError(Exception):
    """Raised when CLAUDE_API_KEY is not configured."""

    def __init__(self):
        super().__init__(
            "lotusmarket: CLAUDE_API_KEY not configured. "
            "Get your key at https://console.anthropic.com/settings/keys\n\n"
            "Usage:\n"
            "  from lotusmarket.ai import AIClient, AIConfig\n"
            "  client = AIClient(AIConfig(api_key='sk-ant-...'))\n"
            "  # or set CLAUDE_API_KEY environment variable"
        )


DEFAULT_MODEL = "claude-sonnet-4-6"
DEFAULT_MAX_TOKENS = 4096


@dataclass
class AIConfig:
    """Configuration for Claude API client.

    Args:
        api_key: Anthropic API key. If empty, reads from CLAUDE_API_KEY env var.
        model: Model ID. Default: claude-sonnet-4-6
        max_tokens: Max output tokens. Default: 4096
    """

    api_key: str = ""
    model: str = DEFAULT_MODEL
    max_tokens: int = DEFAULT_MAX_TOKENS


@dataclass
class AnalysisResult:
    """Result from AI analysis."""

    text: str = ""
    model: str = ""
    tokens_in: int = 0
    tokens_out: int = 0


class AIClient:
    """Claude API client for Vietnamese stock market analysis.

    All analysis outputs are in Vietnamese, using simple language
    suitable for retail investors.

    Example:
        from lotusmarket.ai import AIClient, AIConfig
        client = AIClient(AIConfig(api_key="sk-ant-..."))
        result = client.ask_question("ACB có nên mua không?")
        print(result.text)
    """

    def __init__(self, config: Optional[AIConfig] = None):
        """Initialize AI client.

        Args:
            config: AIConfig with api_key. If None, reads CLAUDE_API_KEY from env.

        Raises:
            MissingAPIKeyError: If no API key is provided or found in env.
        """
        if config is None:
            config = AIConfig()
        api_key = config.api_key or os.environ.get("CLAUDE_API_KEY", "")
        if not api_key:
            raise MissingAPIKeyError()

        try:
            import anthropic
        except ImportError:
            raise ImportError(
                "anthropic package required for AI features. "
                "Install with: pip install lotusmarket[ai] or pip install anthropic"
            )

        self._client = anthropic.Anthropic(api_key=api_key)
        self._model = config.model or DEFAULT_MODEL
        self._max_tokens = config.max_tokens or DEFAULT_MAX_TOKENS

    def prompt(self, text: str) -> AnalysisResult:
        """Send a custom prompt to Claude."""
        return self._call(text)

    def analyze_trend(self, ticker: str, closes: pd.Series) -> AnalysisResult:
        """Analyze stock price trends with pre-computed technical indicators.

        Args:
            ticker: Stock ticker symbol (e.g. "ACB")
            closes: pandas Series of close prices (at least 5 values)

        Returns:
            AnalysisResult with Vietnamese trend analysis
        """
        if len(closes) < 5:
            raise ValueError(f"Need at least 5 prices, got {len(closes)}")

        # Pre-compute technical indicators locally
        d = compute_dashboard(closes)

        parts = [
            f"Bạn là chuyên gia phân tích kỹ thuật chứng khoán Việt Nam. Phân tích xu hướng giá cổ phiếu {ticker}.\n",
            "CHỈ SỐ KỸ THUẬT (đã tính toán sẵn):",
            f"- RSI(14): {d.rsi:.1f}",
        ]
        if d.ma20 is not None:
            parts.append(f"- MA20: {d.ma20:.0f}")
        if d.ma50 is not None:
            parts.append(f"- MA50: {d.ma50:.0f}")
        if d.ma200 is not None:
            parts.append(f"- MA200: {d.ma200:.0f}")
        parts.extend(
            [
                f"- Momentum(20): {d.momentum:.1f}%",
                f"- Tín hiệu: {d.signal}",
                f"- Điểm: {d.score:.0f}/100\n",
                "DỮ LIỆU GIÁ GẦN ĐÂY (20 phiên):",
            ]
        )

        # Last 20 prices
        recent = closes.tail(20)
        for i, price in enumerate(recent):
            parts.append(f"  {i + 1}. {price:.0f}")

        parts.extend(
            [
                "\nHãy phân tích bằng tiếng Việt đơn giản:",
                "1. Xu hướng giá ngắn hạn và trung hạn",
                "2. Vùng hỗ trợ và kháng cự",
                "3. Phân tích các đường MA",
                "4. Đánh giá đà tăng/giảm dựa trên RSI và momentum",
                "5. Mức độ rủi ro",
                "\nQUY TẮC: Toàn bộ nội dung PHẢI viết bằng tiếng Việt, dùng từ ngữ dễ hiểu cho nhà đầu tư cá nhân.",
            ]
        )

        return self._call("\n".join(parts))

    def ask_question(self, question: str) -> AnalysisResult:
        """Answer a market-related question in Vietnamese.

        Args:
            question: Any market/stock question in Vietnamese or English

        Returns:
            AnalysisResult with Vietnamese answer
        """
        prompt = (
            "Bạn là chuyên gia tư vấn chứng khoán Việt Nam. Trả lời câu hỏi sau:\n\n"
            f"{question}\n\n"
            "Hãy trả lời rõ ràng, dễ hiểu bằng tiếng Việt. "
            "Nếu câu hỏi về cổ phiếu cụ thể, đưa ra các điểm dữ liệu liên quan.\n"
            "QUY TẮC: Toàn bộ câu trả lời PHẢI bằng tiếng Việt, dùng từ ngữ đơn giản cho người mới đầu tư."
        )
        return self._call(prompt)

    def analyze_with_context(
        self, data_context: str, instruction: str
    ) -> AnalysisResult:
        """AI analysis with custom data context.

        Args:
            data_context: Pre-formatted data (prices, indicators, news, etc.)
            instruction: What to analyze

        Returns:
            AnalysisResult with Vietnamese analysis
        """
        prompt = (
            f"{data_context}\n\n{instruction}\n\n"
            "QUY TẮC: Toàn bộ nội dung PHẢI viết bằng tiếng Việt, dùng từ ngữ dễ hiểu."
        )
        return self._call(prompt)

    def _call(self, prompt: str) -> AnalysisResult:
        """Send prompt to Claude with retry logic."""
        import time as _time

        last_err = None
        for attempt in range(3):
            if attempt > 0:
                _time.sleep(2**attempt)
            try:
                message = self._client.messages.create(
                    model=self._model,
                    max_tokens=self._max_tokens,
                    messages=[{"role": "user", "content": prompt}],
                )
                text = ""
                for block in message.content:
                    if block.type == "text":
                        text += block.text
                return AnalysisResult(
                    text=text,
                    model=self._model,
                    tokens_in=message.usage.input_tokens,
                    tokens_out=message.usage.output_tokens,
                )
            except Exception as e:
                last_err = e

        raise ConnectionError(f"Claude API failed after 3 attempts: {last_err}")
