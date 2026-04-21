// Package sentiment — Vietnamese news title clustering.
//
// Groups similar news titles together (same story from different sources)
// using deterministic Jaccard similarity. No external calls, no dependencies.
// Useful for deduplicating news feeds that republish wire stories.
//
// Example:
//
//	titles := []string{
//	    "VN-Index tăng 1% phiên cuối tuần",
//	    "Vn-Index đóng cửa tăng 1% chiều thứ Sáu",
//	    "Giá dầu WTI vượt 85 USD",
//	}
//	groups := sentiment.ClusterTitles(titles, 0.3)
//	// groups = [[0, 1], [2]]
package sentiment

import (
	"strings"
	"unicode"
)

// ClusterTitles groups title indices by Jaccard similarity. Titles with Jaccard
// ≥ threshold share a group. Greedy: first unclaimed title seeds a group, then
// consumes every later title above threshold.
//
// threshold=0.3 is a sensible default (share ≥30% of significant tokens).
func ClusterTitles(titles []string, threshold float64) [][]int {
	if len(titles) == 0 {
		return nil
	}
	tokenSets := make([]map[string]bool, len(titles))
	for i, t := range titles {
		tokenSets[i] = TokenizeViTitle(t)
	}

	clusterOf := make([]int, len(titles))
	for i := range clusterOf {
		clusterOf[i] = -1
	}
	numClusters := 0
	for i := range titles {
		if clusterOf[i] != -1 {
			continue
		}
		clusterOf[i] = numClusters
		numClusters++
		for j := i + 1; j < len(titles); j++ {
			if clusterOf[j] != -1 {
				continue
			}
			if Jaccard(tokenSets[i], tokenSets[j]) >= threshold {
				clusterOf[j] = clusterOf[i]
			}
		}
	}

	buckets := make(map[int][]int, numClusters)
	for i, c := range clusterOf {
		buckets[c] = append(buckets[c], i)
	}
	out := make([][]int, 0, numClusters)
	// Preserve original cluster order (by first-seen index)
	for c := 0; c < numClusters; c++ {
		if b, ok := buckets[c]; ok {
			out = append(out, b)
		}
	}
	return out
}

// TokenizeViTitle produces a set of significant tokens from a Vietnamese title.
// Lowercased, split on non-letter/digit, stopwords and short tokens (≤2 chars) dropped.
// Exposed so callers can pre-compute or inspect tokens.
func TokenizeViTitle(title string) map[string]bool {
	title = strings.ToLower(strings.TrimSpace(title))
	var buf strings.Builder
	for _, r := range title {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			buf.WriteRune(r)
		} else {
			buf.WriteRune(' ')
		}
	}
	out := make(map[string]bool)
	for _, t := range strings.Fields(buf.String()) {
		if len(t) <= 2 {
			continue
		}
		if viStopwords[t] {
			continue
		}
		out[t] = true
	}
	return out
}

// Jaccard similarity of two token sets: |A∩B| / |A∪B|.
// Returns 0 for empty sets.
func Jaccard(a, b map[string]bool) float64 {
	if len(a) == 0 || len(b) == 0 {
		return 0
	}
	inter := 0
	for k := range a {
		if b[k] {
			inter++
		}
	}
	union := len(a) + len(b) - inter
	if union == 0 {
		return 0
	}
	return float64(inter) / float64(union)
}

// viStopwords: common Vietnamese news fillers that don't help distinguish stories.
var viStopwords = map[string]bool{
	"đang": true, "được": true, "này": true, "đó": true, "các": true,
	"những": true, "một": true, "hai": true, "khi": true, "nếu": true,
	"nhưng": true, "hơn": true, "như": true, "cùng": true, "vẫn": true,
	"sau": true, "trước": true, "trong": true, "ngoài": true, "trên": true,
	"dưới": true, "giữa": true, "bởi": true, "đến": true, "theo": true,
	"với": true, "không": true, "chưa": true, "đã": true, "sẽ": true,
	"có": true, "phải": true, "nên": true, "cần": true, "đây": true,
	"tại": true, "lại": true, "lên": true, "xuống": true, "ra": true,
	"vào": true, "thêm": true, "bớt": true, "hôm": true, "nay": true,
	"mai": true, "qua": true, "năm": true, "tháng": true, "tuần": true,
}
