package models

import "time"

type Post struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"not null;size:200;index"` // Add index for search
	Body      string    `json:"body" gorm:"type:text"`
	UserID    uint      `json:"user_id" gorm:"index;not null"`
	User      User      `json:"author"`
	Likes     []Like    `json:"-"`
	Comments  []Comment `json:"-"`
	CreatedAt time.Time `json:"created_at" gorm:"index"` // Add index for sorting
	UpdatedAt time.Time `json:"updated_at"`
}
