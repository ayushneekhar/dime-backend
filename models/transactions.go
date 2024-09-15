package models

import "time"

type Transaction struct {
	ID          uint `gorm:"primaryKey"`
	Timestamp   time.Time
	Description string
	Amount      float64
	CategoryID  uint
	Category    Category
	Recurring   bool
}
