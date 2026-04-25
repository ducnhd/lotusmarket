package market

import "testing"

func TestClassify_Stable(t *testing.T) {
	s := ScoreSignals(12, 0.3, 100, 0, "", 0)
	regime, sub := Classify(s)
	if regime != RegimeStable {
		t.Errorf("got %s, want STABLE (composite=%d)", regime, s.CompositeScore)
	}
	if sub != "" {
		t.Errorf("sub-type should be empty for STABLE, got %q", sub)
	}
}

func TestClassify_Volatile_HighVIX(t *testing.T) {
	s := ScoreSignals(26, -1.5, -400, 0, "", -1.0)
	regime, sub := Classify(s)
	if regime != RegimeVolatile {
		t.Errorf("got %s, want VOLATILE (composite=%d)", regime, s.CompositeScore)
	}
	if sub != "" {
		t.Errorf("sub-type should be empty for VOLATILE, got %q", sub)
	}
}

func TestClassify_Volatile_ModerateVNIndex(t *testing.T) {
	s := ScoreSignals(22, -2.5, 0, 0, "", 0)
	regime, sub := Classify(s)
	if regime != RegimeVolatile {
		t.Errorf("got %s, want VOLATILE (composite=%d)", regime, s.CompositeScore)
	}
	_ = sub
}

func TestClassify_Crisis_GlobalDrop(t *testing.T) {
	s := ScoreSignals(35, -4.0, -3000, 0, "", -6.0)
	regime, sub := Classify(s)
	if regime != RegimeCrisis {
		t.Errorf("got %s, want CRISIS (composite=%d)", regime, s.CompositeScore)
	}
	if sub != CrisisPanic {
		t.Errorf("sub-type = %s, want CRISIS_PANIC (global contagion, no VN keywords)", sub)
	}
}

func TestClassify_CrisisFundamental_NewsKeyword(t *testing.T) {
	// "bắt giữ" in news → fundamental
	s := ScoreSignals(32, -3.5, -2500, 1, "chủ tịch bị bắt giữ vì gian lận", -2.0)
	regime, sub := Classify(s)
	if regime != RegimeCrisis {
		t.Errorf("got %s, want CRISIS", regime)
	}
	if sub != CrisisFundamental {
		t.Errorf("sub = %s, want CRISIS_FUNDAMENTAL for 'bắt giữ' keyword", sub)
	}
}

func TestClassify_Euphoria(t *testing.T) {
	// High composite + upward direction → euphoria
	s := ScoreSignals(35, +4.0, 0, 0, "", 0)
	// VNIndex >= 3.0 → score=80, direction=up
	regime, sub := Classify(s)
	if regime != RegimeEuphoria {
		t.Errorf("got %s, want EUPHORIA (composite=%d, dir=%s)", regime, s.CompositeScore, s.Direction)
	}
	if sub != "" {
		t.Errorf("sub should be empty for EUPHORIA, got %q", sub)
	}
}

func TestClassify_CrisisPanic_HeavyForeignSell(t *testing.T) {
	// Heavy foreign sell, no VN news → panic
	s := ScoreSignals(28, -2.5, -3000, 0, "", -1.0)
	regime, sub := Classify(s)
	if regime != RegimeCrisis {
		t.Errorf("got %s, want CRISIS (composite=%d)", regime, s.CompositeScore)
	}
	if sub != CrisisPanic {
		t.Errorf("sub = %s, want CRISIS_PANIC (heavy sell, no fundamental news)", sub)
	}
}

func TestScoreSignals_Thresholds(t *testing.T) {
	// VIX >= 30 → score 80
	s := ScoreSignals(31, 0, 0, 0, "", 0)
	if s.VIXScore != 80 {
		t.Errorf("VIXScore = %d, want 80 for VIX=31", s.VIXScore)
	}

	// Foreign sell > 2000B → score 80
	s2 := ScoreSignals(12, 0, -2500, 0, "", 0)
	if s2.FlowScore != 80 {
		t.Errorf("FlowScore = %d, want 80 for -2500B", s2.FlowScore)
	}

	// News tier 1 → score 100
	s3 := ScoreSignals(12, 0, 0, 1, "circuit breaker triggered", 0)
	if s3.NewsScore != 100 {
		t.Errorf("NewsScore = %d, want 100 for tier 1", s3.NewsScore)
	}

	// Global drop -5.5% → score 80
	s4 := ScoreSignals(12, 0, 0, 0, "", -5.5)
	if s4.GlobalScore != 80 {
		t.Errorf("GlobalScore = %d, want 80 for -5.5%%", s4.GlobalScore)
	}
}

func TestScoreSignals_Direction(t *testing.T) {
	tests := []struct {
		change float64
		want   string
	}{
		{-2.0, "down"},
		{+2.0, "up"},
		{0.5, "mixed"},
		{-0.5, "mixed"},
	}
	for _, tt := range tests {
		s := ScoreSignals(12, tt.change, 0, 0, "", 0)
		if s.Direction != tt.want {
			t.Errorf("change=%.1f → direction=%s, want %s", tt.change, s.Direction, tt.want)
		}
	}
}

func TestIdentifyTrigger_NewsHighest(t *testing.T) {
	s := ScoreSignals(12, -3.5, 0, 1, "FED tăng lãi suất đột ngột", 0)
	// VNIndex score=80, News score=100 → news wins
	source, detail := IdentifyTrigger(s)
	if source != "news_tier1" {
		t.Errorf("source = %s, want news_tier1", source)
	}
	if detail != "FED tăng lãi suất đột ngột" {
		t.Errorf("detail = %q, want the news detail", detail)
	}
}

func TestIdentifyTrigger_VIX(t *testing.T) {
	// Only VIX elevated
	s := ScoreSignals(26, 0.3, 0, 0, "", 0)
	source, _ := IdentifyTrigger(s)
	if source != "vix" {
		t.Errorf("source = %s, want vix", source)
	}
}
