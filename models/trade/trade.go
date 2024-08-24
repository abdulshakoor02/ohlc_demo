package trade

type Trade struct {
	Stream string `json:"stream"`
	Data   Data   `json:"data"`
}

type Data struct {
	EventType   string  `json:"e"`
	EventTime   int64   `json:"E"`
	Price       float64 `json:"p,string"`
	Quantity    float64 `json:"q,string"`
	TradingPair string  `json:"s"`
}
