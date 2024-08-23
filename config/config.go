package config

import (
	"os"

	"github.com/abdulshakoor02/ohlc_exinity/logger"
	"github.com/joho/godotenv"
)

var POSTGRES_HOST string
var DB_NAME string
var POSTGRES_USER string
var POSTGRES_PASSWORD string
var POSTGRES_PORT string

func LoadEnv() {
	log := logger.Logger
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to load env")
	}
	POSTGRES_HOST = os.Getenv("POSTGRES_HOST")
	DB_NAME = os.Getenv("DB_NAME")
	POSTGRES_USER = os.Getenv("POSTGRES_USER")
	POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
	POSTGRES_PORT = os.Getenv("POSTGRES_PORT")
}
