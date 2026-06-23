package main

// hnx30 — HNX30 index constituents. Deferred for v1: the HNX30 basket is
// reviewed twice a year and a reliably-current list was not available at build
// time, so the event scan runs on VN30 only. To add HNX30 later, populate this
// slice with the official current constituents (https://hnx.vn) — marketUniverse
// already merges and dedups it.
var hnx30 = []string{}

// marketUniverse returns the dedup'd union of VN30 and HNX30 tickers.
func marketUniverse() []string {
	seen := map[string]bool{}
	out := []string{}
	for _, tk := range append(append([]string{}, vn30...), hnx30...) {
		if !seen[tk] {
			seen[tk] = true
			out = append(out, tk)
		}
	}
	return out
}
