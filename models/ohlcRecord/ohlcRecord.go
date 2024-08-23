package ohlcRecord

import (
	"time"

	"gorm.io/gorm"
)

type OhlcRecord struct {
	ID          string         `json:"id"           gorm:"type:string;size:100;primary_key;default:gen_random_uuid()"`
	TradingPair string         `json:"trading_pair" gorm:"size:100"`
	OpenTime    time.Time      `json:"open_time"    gorm:"size:100"`
	CloseTime   time.Time      `json:"close_time"   gorm:"size:100"`
	Open        float64        `json:"open"         gorm:"size:100"`
	High        float64        `json:"high"         gorm:"size:100"`
	Low         float64        `json:"low"          gorm:"size:100"`
	Close       float64        `json:"close"        gorm:"size:100"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `                    gorm:"index"`
}
