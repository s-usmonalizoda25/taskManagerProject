package models

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Username  string `gorm:"type:varchar(50);unique;not null;" json:"username"`
	Email     string `gorm:"type:varchar(100);unique;not null" json:"email"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   *time.Time `gorm:"autoUpdateTime:false;default:null" json:"updated_at,omitempty"`
	DeletedAt   *time.Time `gorm:"index;default:null" json:"deleted_at,omitempty"`
	Tasks     []Task  `gorm:"foreignKey:UserID" json:"tasks,omitempty"`
}