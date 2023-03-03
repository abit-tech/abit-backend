package models

import "gorm.io/gorm"

type OldUser struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string
}
