import os

os.environ["LOTUSMARKET_QUIET"] = "1"

import pytest
from lotusmarket.ai import (
    AIConfig,
    MissingAPIKeyError,
    DEFAULT_MODEL,
    DEFAULT_MAX_TOKENS,
)


def test_missing_api_key():
    """Should raise MissingAPIKeyError with clear instructions."""
    # Make sure env var is not set
    os.environ.pop("CLAUDE_API_KEY", None)
    with pytest.raises(MissingAPIKeyError) as exc_info:
        from lotusmarket.ai import AIClient

        AIClient()
    msg = str(exc_info.value)
    assert "CLAUDE_API_KEY" in msg
    assert "console.anthropic.com" in msg


def test_missing_api_key_explicit_empty():
    """Explicit empty key should also raise."""
    os.environ.pop("CLAUDE_API_KEY", None)
    with pytest.raises(MissingAPIKeyError):
        from lotusmarket.ai import AIClient

        AIClient(AIConfig(api_key=""))


def test_config_defaults():
    """Default config should use sonnet model."""
    cfg = AIConfig()
    assert cfg.model == DEFAULT_MODEL
    assert cfg.max_tokens == DEFAULT_MAX_TOKENS
    assert cfg.api_key == ""


def test_config_custom():
    """Custom config should override defaults."""
    cfg = AIConfig(api_key="sk-test", model="claude-opus-4-6", max_tokens=8192)
    assert cfg.api_key == "sk-test"
    assert cfg.model == "claude-opus-4-6"
    assert cfg.max_tokens == 8192


def test_error_message_helpful():
    """Error message should include setup instructions."""
    err = MissingAPIKeyError()
    msg = str(err)
    assert "pip install" in msg or "CLAUDE_API_KEY" in msg
    assert "console.anthropic.com" in msg
