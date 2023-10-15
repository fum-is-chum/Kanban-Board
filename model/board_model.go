package model

import "gorm.io/gorm"

type Board struct {
	gorm.Model
	Name    string  `json:"name,omitempty"`
	Desc    string  `json:"desc,omitempty"`
	OwnerID uint    `json:"owner_id,omitempty"`
	Owner   *User   `json:"owner"`
	Members []*User `json:"members" gorm:"many2many:board_members;constraint:OnDelete:CASCADE;"`
}
