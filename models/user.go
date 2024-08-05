package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey"`
	CreateAt  time.Time
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
