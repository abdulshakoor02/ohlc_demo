package trade

type Trade struct {
	EventType   string  `json:"e"`
	EventTime   int64   `json:"E"`
	Price       float64 `json:"p,string"`
	Quantity    float64 `json:"q,string"`
	TradingPair string  `json:"s"`
}
