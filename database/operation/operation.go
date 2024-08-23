package operation

import (
	"encoding/json"

	"github.com/abdulshakoor02/ohlc_exinity/database/dbAdapter"
	"github.com/abdulshakoor02/ohlc_exinity/logger"
)

func Create[T any](data *T) {
	log := logger.Logger
	logData, err := json.Marshal(data)
	if err != nil {
		log.Err(err).Msg("failed to Marshal data %v")
	}
	log.Info().Msgf("inserting data %v\n", string(logData))
	if err := dbAdapter.DB.Create(&data).Error; err != nil {
		log.Fatal().Err(err).Msgf("failed to insert data %v\n", string(logData))
	}
	log.Info().Msgf("data inserted %v\n", string(logData))
}
