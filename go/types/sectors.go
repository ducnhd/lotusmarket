package types

var VN30 = []string{
	"ACB", "BID", "CTG", "DGC", "FPT", "GAS", "GVR", "HDB",
	"HPG", "LPB", "MBB", "MSN", "MWG", "PLX", "SAB", "SHB",
	"SSB", "SSI", "STB", "TCB", "TPB", "VCB", "VHM", "VIB",
	"VIC", "VJC", "VNM", "VPB", "VPL", "VRE",
}

var Sectors = map[string][]string{
	"Ngân hàng":    {"ACB", "BID", "CTG", "HDB", "LPB", "MBB", "SHB", "SSB", "STB", "TCB", "TPB", "VCB", "VIB", "VPB"},
	"Tài chính":    {"SSI"},
	"Công nghệ":    {"FPT"},
	"Bán lẻ":       {"MWG"},
	"Bất động sản": {"VHM", "VIC", "VRE"},
	"Năng lượng":   {"GAS", "PLX"},
	"Tiêu dùng":    {"MSN", "SAB", "VNM"},
	"Vật liệu":     {"HPG"},
	"Hóa chất":     {"DGC"},
	"Nông nghiệp":  {"GVR"},
	"Vận tải":      {"VJC"},
	"Du lịch":      {"VPL"},
}

var sectorIndex map[string]string

func init() {
	sectorIndex = make(map[string]string)
	for sector, tickers := range Sectors {
		for _, t := range tickers {
			sectorIndex[t] = sector
		}
	}
}

func GetSector(ticker string) string {
	if s, ok := sectorIndex[ticker]; ok {
		return s
	}
	return "Others"
}
