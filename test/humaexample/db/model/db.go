package model

import (
	"time"

	"github.com/google/uuid"
)

type Book struct {
	ID          uuid.UUID `gorm:"id;primaryKey"`
	DisplayName string    `gorm:"display_name"`
	CreatedAt   time.Time `gorm:"created_at"`
	UpdatedAt   time.Time `gorm:"updated_at"`
}
