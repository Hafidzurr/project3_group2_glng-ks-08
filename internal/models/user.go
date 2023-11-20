package models

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	FullName  string `json:"full_name" validate:"required"`
	Email     string `gorm:"unique;not null" json:"email" validate:"required,email"`
	Password  string `gorm:"not null" json:"password" validate:"required,min=6"`
	Role      string `gorm:"not null" json:"role" validate:"required,oneof=admin member"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
