package main

import (
	"encoding/json"
	"time"

	"github.com/abdulshakoor02/ohlc_exinity/config"
	"github.com/abdulshakoor02/ohlc_exinity/database/dbAdapter"
	"github.com/abdulshakoor02/ohlc_exinity/database/migration"
	"github.com/abdulshakoor02/ohlc_exinity/database/operation"
	"github.com/abdulshakoor02/ohlc_exinity/logger"
	"github.com/abdulshakoor02/ohlc_exinity/models/ohlcRecord"
	"github.com/abdulshakoor02/ohlc_exinity/models/trade"
	"github.com/gorilla/websocket"
)

func main() {
	log := logger.Logger
	config.LoadEnv()
	dbAdapter.DbConnect()
	migration.MigrateDb()

	url := "wss://stream.binance.com:9443/ws/btcusdt@aggTrade"

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to load env")
	}
	defer conn.Close()
	var currentOHLC *ohlcRecord.OhlcRecord
	var currentMinute time.Time

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
		tradeTime := time.Unix(0, trade.EventTime*int64(time.Millisecond))
		minuteStart := tradeTime.Truncate(time.Minute)

		if currentOHLC == nil || minuteStart != currentMinute {
			if currentOHLC != nil {
				go operation.Create[ohlcRecord.OhlcRecord](currentOHLC)
			}
			currentOHLC = &ohlcRecord.OhlcRecord{
				TradingPair: trade.TradingPair,
				OpenTime:    minuteStart,
				CloseTime:   minuteStart.Add(time.Minute - time.Nanosecond),
				Open:        trade.Price,
				High:        trade.Price,
				Low:         trade.Price,
				Close:       trade.Price,
			}
			currentMinute = minuteStart
		} else {
			currentOHLC.Close = trade.Price
			if trade.Price > currentOHLC.High {
				currentOHLC.High = trade.Price
			}
			if trade.Price < currentOHLC.Low {
				currentOHLC.Low = trade.Price
			}
		}

		// finalVal, err := json.Marshal(currentOHLC)
		// if err != nil {
		// 	fmt.Println("Error:", err)
		// }
		//
		// fmt.Printf("Received message: %s\n", finalVal)
	}
}
