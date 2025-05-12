package entity

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:100"`
	Email     string `gorm:"uniqueIndex"`
	Password  string `gorm:"size:255"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
