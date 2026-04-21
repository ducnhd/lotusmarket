package sentiment

import "testing"

func TestTokenizeViTitle(t *testing.T) {
	tokens := TokenizeViTitle("VN-Index tăng 1% phiên cuối tuần")
	if !tokens["index"] || !tokens["tăng"] || !tokens["phiên"] || !tokens["cuối"] {
		t.Errorf("expected main tokens, got %v", tokens)
	}
	// Stopwords should be dropped
	if tokens["tuần"] {
		t.Errorf("'tuần' is stopword, should be dropped")
	}
	// Short tokens (≤2 chars) dropped
	if tokens["1"] || tokens["vn"] {
		t.Errorf("short tokens should be dropped")
	}
}

func TestJaccard(t *testing.T) {
	a := map[string]bool{"x": true, "y": true, "z": true}
	b := map[string]bool{"x": true, "y": true, "w": true}
	got := Jaccard(a, b)
	// intersection = 2 (x, y), union = 4 (x,y,z,w) → 0.5
	if got != 0.5 {
		t.Errorf("Jaccard = %.3f, want 0.5", got)
	}
	if Jaccard(nil, a) != 0 {
		t.Errorf("empty set should give 0")
	}
}

func TestClusterTitles_Basic(t *testing.T) {
	titles := []string{
		"VN-Index tăng 1% phiên cuối tuần",
		"VNIndex đóng cửa tăng 1% phiên cuối",
		"Giá dầu WTI vượt 85 USD",
	}
	groups := ClusterTitles(titles, 0.3)
	if len(groups) != 2 {
		t.Errorf("expected 2 clusters, got %d: %v", len(groups), groups)
	}
	// First two should be grouped
	if len(groups[0]) != 2 {
		t.Errorf("first cluster should have 2 titles, got %d", len(groups[0]))
	}
}

func TestClusterTitles_Empty(t *testing.T) {
	if got := ClusterTitles(nil, 0.3); got != nil {
		t.Errorf("nil input → nil output, got %v", got)
	}
}

func TestClusterTitles_AllDistinct(t *testing.T) {
	titles := []string{
		"Ngân hàng tăng vốn điều lệ",
		"Dầu khí báo lãi quý 3",
		"Bất động sản đón sóng mới",
	}
	groups := ClusterTitles(titles, 0.3)
	if len(groups) != 3 {
		t.Errorf("distinct titles → 3 clusters, got %d", len(groups))
	}
}
