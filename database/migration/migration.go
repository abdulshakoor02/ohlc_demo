package migration

import (
	"github.com/abdulshakoor02/ohlc_exinity/database/dbAdapter"
	"github.com/abdulshakoor02/ohlc_exinity/models/ohlcRecord"
)

func MigrateDb() {
	dbAdapter.DB.AutoMigrate(&ohlcRecord.OhlcRecord{})
}
