package models

import "time"

type TaskStatus string

const (
	StatusNew      TaskStatus = "new"
	StatusDone     TaskStatus = "done"
	StatusCanceled TaskStatus = "canceled"
)

type Task struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	UserID      uint       `gorm:"not null" json:"user_id"`
	Title       string     `gorm:"type:varchar(100);not null" json:"title"`
	Description string     `gorm:"type:text" json:"description"`
	Status      TaskStatus `gorm:"type:varchar(20);default:'new';not null" json:"status"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   *time.Time `gorm:"autoUpdateTime:false;default:null" json:"updated_at,omitempty"`
	DeletedAt   *time.Time `gorm:"index;default:null" json:"deleted_at,omitempty"`
}

type UpdateStatus struct {
	Status TaskStatus `json:"status"`
}
