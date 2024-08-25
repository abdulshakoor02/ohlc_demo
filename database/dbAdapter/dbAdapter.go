package dbAdapter

import (
	"fmt"

	"github.com/abdulshakoor02/ohlc_exinity/config"
	"github.com/abdulshakoor02/ohlc_exinity/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DbConnect() *gorm.DB {
	log := logger.Logger
	config.LoadEnv()

	log.Info().Msgf(config.DB_NAME)
	dsn := fmt.Sprintf(
		"user=%v password=%v host=%v port=%v dbname=%v ",
		config.POSTGRES_USER,
		config.POSTGRES_PASSWORD,
		config.POSTGRES_HOST,
		config.POSTGRES_PORT,
		config.DB_NAME,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Print(err)
		log.Fatal().Err(err).Msgf("failed to load env")
	}
	DB = db
	pgdb, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to initiate a database instance :%v", err)
	}
	pgdb.Ping()
	pgdb.SetMaxIdleConns(10)
	pgdb.SetMaxOpenConns(100)
	log.Info().Msg("Postgres DB connection established")
	return db
}
