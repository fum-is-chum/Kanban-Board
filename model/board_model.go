package model

import "gorm.io/gorm"

type Board struct {
	gorm.Model
	Name    string  `json:"name"`
	Desc    string  `json:"desc"`
	OwnerID uint    `json:"owner_id"`
	Owner   *User   `json:"owner"`
	Members []*User `json:"members" gorm:"many2many:board_members"`
}
