package models

import (
	"time"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null;size:100"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null;size:255"`
	Password  string    `json:"-" gorm:"not null"` // bcrypt hash
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Posts     []Post    `json:"-"`
}
