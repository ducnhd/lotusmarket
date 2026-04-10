package ai

import (
	"testing"
)

func TestNew_NoAPIKey(t *testing.T) {
	_, err := New(Config{})
	if err != ErrNoAPIKey {
		t.Errorf("got %v, want ErrNoAPIKey", err)
	}
}

func TestNew_Defaults(t *testing.T) {
	c, err := New(Config{APIKey: "test-key"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.model != DefaultModel {
		t.Errorf("model = %q, want %q", c.model, DefaultModel)
	}
	if c.maxTokens != DefaultMaxTokens {
		t.Errorf("maxTokens = %d, want %d", c.maxTokens, DefaultMaxTokens)
	}
}

func TestNew_CustomModel(t *testing.T) {
	c, err := New(Config{
		APIKey:    "test-key",
		Model:     "claude-opus-4-6",
		MaxTokens: 8192,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.model != "claude-opus-4-6" {
		t.Errorf("model = %q, want claude-opus-4-6", c.model)
	}
	if c.maxTokens != 8192 {
		t.Errorf("maxTokens = %d, want 8192", c.maxTokens)
	}
}

func TestErrNoAPIKey_Message(t *testing.T) {
	msg := ErrNoAPIKey.Error()
	if msg == "" {
		t.Error("ErrNoAPIKey message should not be empty")
	}
	if !contains(msg, "CLAUDE_API_KEY") {
		t.Errorf("message should mention CLAUDE_API_KEY: %s", msg)
	}
	if !contains(msg, "console.anthropic.com") {
		t.Errorf("message should include setup URL: %s", msg)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
