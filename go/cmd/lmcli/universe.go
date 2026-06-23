package main

// hnx30 — HNX30 index constituents. VERIFY against the current official list
// before merge; membership is reviewed periodically by HNX.
var hnx30 = []string{
	"SHS", "PVS", "CEO", "MBS", "IDC", "VCS", "HUT", "TNG", "PVI", "DTD",
	"L14", "NVB", "BVS", "TIG", "LAS", "VGS", "DXP", "NRC", "PVC", "TVC",
	"API", "IDV", "NTP", "PLC", "MST", "S99", "DDG", "AMV", "ART", "VC3",
}

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
