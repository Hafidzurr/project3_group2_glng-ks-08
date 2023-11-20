package models

import "time"

type Category struct {
	ID        uint `gorm:"primaryKey"`
	Type      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Tasks     []Task `gorm:"foreignKey:CategoryID"` // Relasi ke Task
}
