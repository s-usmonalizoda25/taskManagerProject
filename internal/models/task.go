package models

import "time"

type Task struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"not null" json:"title"`
	Description string    `json:"description"`
	Status      string    `gorm:"default:'new'" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
