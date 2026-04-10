package fetchers

import "testing"

func TestParseVPSResponse_Valid(t *testing.T) {
	raw := `[{"sym":"ACB","lastPrice":26.9,"r":26.5,"openPrice":26.6,"highPrice":27.0,"lowPrice":26.3,"lot":1234567,"fBVol":50000,"fSVolume":30000,"c":28.3,"f":24.7,"g1":"26.8|100570|d","g2":"27.0|80000|d"}]`
	stocks, err := parseVPSResponse([]byte(raw))
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(stocks) != 1 {
		t.Fatalf("got %d, want 1", len(stocks))
	}
	s := stocks[0]
	if s.Ticker != "ACB" {
		t.Errorf("Ticker = %q", s.Ticker)
	}
	if s.Close != 26900 {
		t.Errorf("Close = %v, want 26900", s.Close)
	}
	if s.ForeignNetVol != 20000 {
		t.Errorf("ForeignNetVol = %d, want 20000", s.ForeignNetVol)
	}
}

func TestParseVPSResponse_Empty(t *testing.T) {
	stocks, err := parseVPSResponse([]byte(`[]`))
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(stocks) != 0 {
		t.Errorf("got %d, want 0", len(stocks))
	}
}

func TestParseEntradeResponse_Valid(t *testing.T) {
	raw := `{"t":[1709510400,1709596800],"o":[26.5,27.0],"h":[27.0,27.5],"l":[26.0,26.5],"c":[26.8,27.2],"v":[1000000,1500000]}`
	stocks, err := parseEntradeResponse([]byte(raw), "ACB")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(stocks) != 2 {
		t.Fatalf("got %d, want 2", len(stocks))
	}
	if stocks[0].Close != 26800 {
		t.Errorf("Close = %v, want 26800", stocks[0].Close)
	}
}

func TestParsePipeVol(t *testing.T) {
	if v := parsePipeVol("22.5|100570|d"); v != 100570 {
		t.Errorf("got %d, want 100570", v)
	}
	if v := parsePipeVol(""); v != 0 {
		t.Errorf("got %d, want 0", v)
	}
}
