package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pixlcrashr/go-pagetoken/test/humaexample/db/model"
	"gorm.io/gorm"
)

func seed(db *gorm.DB) error {
	now := time.Now()
	books := make([]model.Book, 100)
	for i := range books {
		books[i] = model.Book{
			ID:          uuid.New(),
			DisplayName: fmt.Sprintf("Book %03d", i+1),
			CreatedAt:   now.Add(time.Duration(i) * time.Second),
			UpdatedAt:   now.Add(time.Duration(i) * time.Second),
		}
	}
	return db.Create(&books).Error
}
