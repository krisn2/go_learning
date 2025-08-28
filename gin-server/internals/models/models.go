package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm: "uniqueIndex"`
	Password string
	Name     string
}

type Item struct {
	gorm.Model
	Title  string
	UserID uint
}
