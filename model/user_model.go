package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
	Boards   []Board `json:"boards" gorm:"foreignkey:OwnerID;constraint:OnDelete:CASCADE;"`
}
