package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string  `json:"name,omitempty"`
	Email    string  `json:"email,omitempty" gorm:"unique"`
	Password string  `json:"password,omitempty"`
	Boards   []Board `json:"boards" gorm:"foreignkey:OwnerID;constraint:OnDelete:CASCADE;"`
	MemberOf []*Board `json:"member_of" gorm:"many2many:board_members;"`
}
