package aggregateData

import (
	"encoding/json"
	"time"

	"github.com/abdulshakoor02/ohlc_exinity/database/operation"
	"github.com/abdulshakoor02/ohlc_exinity/logger"
	"github.com/abdulshakoor02/ohlc_exinity/models/ohlcRecord"
	"github.com/abdulshakoor02/ohlc_exinity/models/trade"
	pb "github.com/abdulshakoor02/ohlc_exinity/ohlc"
	"github.com/gorilla/websocket"
)

// func AggregateData(tradePair string, insertCandle bool, ohlcChannel chan<- *pb.OHLC) {
func AggregateData(insertCandle bool, ohlcChannel chan<- *pb.OHLC) {
	url := "wss://stream.binance.com:9443/stream?streams=btcusdt@aggTrade/ethusdt@aggTrade/pepeusdt@aggTrade"
	// url := fmt.Sprintf("wss://stream.binance.com:9443/ws/%v@aggTrade", tradePair)
	log := logger.Logger
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to load env")
	}
	defer conn.Close()
	currentOHLC := make(map[string]*ohlcRecord.OhlcRecord)
	currentMinute := make(map[string]time.Time)

	for {
		var trade trade.Trade
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}

		if err := json.Unmarshal(message, &trade); err != nil {
			log.Println("Error unmarshaling JSON:", err)
			continue
		}
		tradeTime := time.Unix(0, trade.Data.EventTime*int64(time.Millisecond))
		minuteStart := tradeTime.Truncate(time.Minute)

		if currentOHLC[trade.Data.TradingPair] == nil ||
			minuteStart != currentMinute[trade.Data.TradingPair] {
			if currentOHLC[trade.Data.TradingPair] != nil {
				if insertCandle {
					go operation.Create[ohlcRecord.OhlcRecord](currentOHLC[trade.Data.TradingPair])
				}
			}
			currentOHLC[trade.Data.TradingPair] = &ohlcRecord.OhlcRecord{
				TradingPair: trade.Data.TradingPair,
				OpenTime:    minuteStart,
				CloseTime:   minuteStart.Add(time.Minute - time.Nanosecond),
				Open:        trade.Data.Price,
				High:        trade.Data.Price,
				Low:         trade.Data.Price,
				Close:       trade.Data.Price,
			}
			currentMinute[trade.Data.TradingPair] = minuteStart
		} else {
			currentOHLC[trade.Data.TradingPair].Close = trade.Data.Price
			if trade.Data.Price > currentOHLC[trade.Data.TradingPair].High {
				currentOHLC[trade.Data.TradingPair].High = trade.Data.Price
			}
			if trade.Data.Price < currentOHLC[trade.Data.TradingPair].Low {
				currentOHLC[trade.Data.TradingPair].Low = trade.Data.Price
			}
		}
		ohlc := &pb.OHLC{
			TradePair: currentOHLC[trade.Data.TradingPair].TradingPair,
			OpenTime:  currentOHLC[trade.Data.TradingPair].OpenTime.Format(time.RFC3339),
			CloseTime: currentOHLC[trade.Data.TradingPair].CloseTime.Format(time.RFC3339),
			Open:      currentOHLC[trade.Data.TradingPair].Open,
			High:      currentOHLC[trade.Data.TradingPair].High,
			Low:       currentOHLC[trade.Data.TradingPair].Low,
			Close:     currentOHLC[trade.Data.TradingPair].Close,
		}
		ohlcChannel <- ohlc
		// log.Println(ohlc)
		// time.Sleep(time.Second * 4)

		// finalVal, err := json.Marshal(currentOHLC[trade.Data.TradingPair])
		// if err != nil {
		// 	log.Println("Error:", err)
		// }
		//
		// log.Printf("Received message: %s\n", finalVal)
	}
}
