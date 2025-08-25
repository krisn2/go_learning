package models

import "time"

type Comment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Body      string    `json:"body" gorm:"type:text;not null"`
	UserID    uint      `json:"user_id" gorm:"index;not null"`
	User      User      `json:"author"`
	PostID    uint      `json:"post_id" gorm:"index;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
