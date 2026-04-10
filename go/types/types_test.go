package types

import "testing"

func TestStockData_ToVND(t *testing.T) {
	s := StockData{
		Open: 26.9, High: 27.5, Low: 26.0, Close: 27.2,
		ChangeValue: 0.3, RefPrice: 26.9, Ceiling: 28.7,
		Floor: 25.1, Bid: 27.1, Ask: 27.3,
	}
	s.ToVND()
	if s.Close != 27200 {
		t.Errorf("Close: got %v, want 27200", s.Close)
	}
	if s.Open != 26900 {
		t.Errorf("Open: got %v, want 26900", s.Open)
	}
	if s.Bid != 27100 {
		t.Errorf("Bid: got %v, want 27100", s.Bid)
	}
}

func TestGetSector(t *testing.T) {
	tests := []struct {
		ticker string
		want   string
	}{
		{"ACB", "Ngân hàng"},
		{"FPT", "Công nghệ"},
		{"HPG", "Vật liệu"},
		{"UNKNOWN", "Others"},
	}
	for _, tt := range tests {
		got := GetSector(tt.ticker)
		if got != tt.want {
			t.Errorf("GetSector(%q) = %q, want %q", tt.ticker, got, tt.want)
		}
	}
}
