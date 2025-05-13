package entity

import (
	"errors"
	"time"
)

var ErrUserNotFound = errors.New("user not found")

type User struct {
	ID        string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name      string `gorm:"type:varchar(100)"`
	Email     string
	Role      string `gorm:"check:role IN ('admin', 'cashier', 'customer')"`
	Password  string `gorm:"type:varchar(255)"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
