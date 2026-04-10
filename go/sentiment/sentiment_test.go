package sentiment

import (
	"math"
	"testing"
)

func TestAnalyze_Positive(t *testing.T) {
	r := Analyze("Ngân hàng tăng mạnh, khối ngoại mua ròng")
	if r.Label != "POSITIVE" {
		t.Errorf("Label = %q, want POSITIVE", r.Label)
	}
	if r.Score <= 0 {
		t.Errorf("Score = %v, want > 0", r.Score)
	}
}

func TestAnalyze_Negative(t *testing.T) {
	r := Analyze("Thị trường giảm sâu, bán tháo hoảng loạn")
	if r.Label != "NEGATIVE" {
		t.Errorf("Label = %q, want NEGATIVE", r.Label)
	}
	if r.Score >= 0 {
		t.Errorf("Score = %v, want < 0", r.Score)
	}
}

func TestAnalyze_Neutral(t *testing.T) {
	r := Analyze("Hôm nay thời tiết đẹp")
	if r.Label != "NEUTRAL" {
		t.Errorf("Label = %q, want NEUTRAL", r.Label)
	}
	if math.Abs(r.Score) > 0.01 {
		t.Errorf("Score = %v, want ~0", r.Score)
	}
}

func TestAnalyze_Empty(t *testing.T) {
	r := Analyze("")
	if r.Label != "NEUTRAL" {
		t.Errorf("Label = %q, want NEUTRAL", r.Label)
	}
}

func TestAnalyze_ScoreBounded(t *testing.T) {
	r := Analyze("tăng vọt bùng nổ đột phá kỷ lục bứt phá thăng hoa tăng trưởng lãi lớn")
	if r.Score > 1.0 {
		t.Errorf("Score = %v, should be <= 1.0", r.Score)
	}
	if r.Score < -1.0 {
		t.Errorf("Score = %v, should be >= -1.0", r.Score)
	}
}
