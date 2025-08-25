package models

import "time"

type Like struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"index;not null"`
	PostID    uint      `json:"post_id" gorm:"index;not null"`
	CreatedAt time.Time `json:"created_at"`
}
