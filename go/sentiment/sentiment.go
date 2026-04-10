package sentiment

import (
	"math"
	"strings"
)

type Result struct {
	Label string  // "POSITIVE", "NEGATIVE", "NEUTRAL"
	Score float64 // -1.0 to 1.0
}

var positiveWords = map[string]float64{
	"tăng": 1.0, "tăng trưởng": 1.5, "tăng mạnh": 1.5, "tăng vọt": 2.0,
	"bùng nổ": 2.0, "đột phá": 1.5, "kỷ lục": 1.5, "bứt phá": 2.0,
	"phục hồi": 1.0, "hồi phục": 1.0, "khởi sắc": 1.5, "lạc quan": 1.5,
	"tích cực": 1.5, "thuận lợi": 1.0, "triển vọng": 1.0, "cơ hội": 1.0,
	"thăng hoa": 2.0, "lãi": 1.0, "lãi lớn": 1.5, "lợi nhuận": 1.0,
	"cổ tức": 1.0, "chia cổ tức": 1.5, "vượt kỳ vọng": 1.5,
	"mua ròng": 1.0, "dòng tiền": 0.5, "hấp dẫn": 1.0,
	"xanh": 1.0, "xanh rực": 1.5, "rực rỡ": 1.5, "lan tỏa": 1.0,
	"sôi động": 1.0, "nâng hạng": 1.5, "vượt đỉnh": 1.5, "tăng trần": 1.5,
}

var negativeWords = map[string]float64{
	"giảm": 1.0, "giảm mạnh": 1.5, "giảm sâu": 2.0, "giảm sốc": 2.0,
	"lao dốc": 2.0, "rớt": 1.5, "rớt mạnh": 2.0, "sụt": 1.5,
	"sụt giảm": 1.5, "sụp đổ": 2.0, "bán tháo": 2.0, "tháo chạy": 2.0,
	"hoảng loạn": 2.0, "lo ngại": 1.0, "bi quan": 1.5, "tiêu cực": 1.5,
	"rủi ro": 1.0, "lỗ": 1.0, "thua lỗ": 1.5, "phá sản": 2.0,
	"nợ xấu": 1.5, "đóng cửa": 1.0, "đình chỉ": 1.5, "cảnh báo": 1.0,
	"giảm sàn": 2.0, "đỏ sàn": 1.5, "bán ròng": 1.0, "rút vốn": 1.5,
}

func Analyze(text string) Result {
	if text == "" {
		return Result{Label: "NEUTRAL", Score: 0}
	}
	lower := strings.ToLower(text)

	posScore := 0.0
	negScore := 0.0

	for phrase, weight := range positiveWords {
		if strings.Contains(lower, phrase) {
			posScore += weight
		}
	}
	for phrase, weight := range negativeWords {
		if strings.Contains(lower, phrase) {
			negScore += weight
		}
	}

	total := posScore + negScore
	if total == 0 {
		return Result{Label: "NEUTRAL", Score: 0}
	}

	score := (posScore - negScore) / total
	score = math.Max(-1, math.Min(1, score))

	label := "NEUTRAL"
	if score > 0.1 {
		label = "POSITIVE"
	} else if score < -0.1 {
		label = "NEGATIVE"
	}

	return Result{Label: label, Score: score}
}
