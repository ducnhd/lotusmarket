package types

// StockData represents a single stock quote with OHLCV + flow data.
type StockData struct {
	Ticker         string
	Date           string
	Open           float64
	High           float64
	Low            float64
	Close          float64
	Volume         int64
	ChangePercent  float64
	ChangeValue    float64
	RefPrice       float64
	Ceiling        float64
	Floor          float64
	ForeignBuyVol  int64
	ForeignSellVol int64
	ForeignNetVol  int64
	Bid            float64
	Ask            float64
	BidVol         int64
	AskVol         int64
}

func (s *StockData) ToVND() {
	s.Open *= 1000
	s.High *= 1000
	s.Low *= 1000
	s.Close *= 1000
	s.ChangeValue *= 1000
	s.RefPrice *= 1000
	s.Ceiling *= 1000
	s.Floor *= 1000
	s.Bid *= 1000
	s.Ask *= 1000
}

type KBSQuote struct {
	Date           string  `json:"TradingDate"`
	Open           float64 `json:"OpenPrice"`
	High           float64 `json:"HighestPrice"`
	Low            float64 `json:"LowestPrice"`
	Close          float64 `json:"ClosePrice"`
	Volume         float64 `json:"TotalVol"`
	MarketCap      float64 `json:"MarketCapital"`
	PE             float64 `json:"PE"`
	PB             float64 `json:"PB"`
	EPS            float64 `json:"EPS"`
	BVPS           float64 `json:"BVPS"`
	Beta           float64 `json:"Beta"`
	DividendYield  float64 `json:"Yield"`
	ForeignBuyVol  float64 `json:"ForeignBuyVol"`
	ForeignSellVol float64 `json:"ForeignSellVol"`
}

const (
	TradingCommissionRate = 0.0015
	SellingTaxRate        = 0.001
)
