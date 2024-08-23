package main

import (
	"github.com/abdulshakoor02/ohlc_exinity/config"
	"github.com/abdulshakoor02/ohlc_exinity/database/dbAdapter"
	"github.com/abdulshakoor02/ohlc_exinity/database/migration"
	"github.com/abdulshakoor02/ohlc_exinity/service/aggregateData"
)

func main() {
	config.LoadEnv()
	dbAdapter.DbConnect()
	migration.MigrateDb()

	go aggregateData.AggregateData("ethusdt", true)
	go aggregateData.AggregateData("pepeusdt", true)
	aggregateData.AggregateData("btcusdt", true)

}
