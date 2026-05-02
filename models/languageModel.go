package models

import "gorm.io/gorm"

type Language struct {
	gorm.Model
	Name string `json:"name" gorm:"size:100;not null"`
	Code string `json:"code" gorm:"size:10;unique;not null"`
}
