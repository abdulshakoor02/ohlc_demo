package ohlc

import (
	"time"
)

type Ohlc struct {
	TradingPair string
	OpenTime    time.Time
	CloseTime   time.Time
	Open        float64
	High        float64
	Low         float64
	Close       float64
}
