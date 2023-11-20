package models

import (
	"time"
)

type Task struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description"`
	Status      bool      `json:"status"`
	UserID      uint      `json:"user_id" gorm:"constraint:OnDelete:CASCADE;"`
	CategoryID  uint      `json:"category_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	User        User      `gorm:"foreignKey:UserID" json:"user"`
}
