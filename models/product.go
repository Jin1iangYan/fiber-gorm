package models

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID           uuid.UUID `json:"id" gorm:"primaryKey"`
	CreateAt     time.Time
	Name         string `json:"name"`
	SerialNumber string `json:"serial_number"`
}
