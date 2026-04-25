package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	GoogleID     string `json:"google_id" gorm:"size:255;uniqueIndex;not null"`
	Email        string `json:"email" gorm:"size:255;uniqueIndex;not null"`
	FullName     string `json:"full_name" gorm:"size:255"`
	ProfileImage string `json:"profile_image" gorm:"size:255"`
}
