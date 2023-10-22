package model

import "gorm.io/gorm"

type Board struct {
	gorm.Model
	Name    string         `json:"name"`
	Desc    string         `json:"desc"`
	Owner   *User          `json:"owner"`
	OwnerID uint           `json:"owner_id,omitempty"`
	Members []*User        `json:"members" gorm:"many2many:board_members;constraint:OnDelete:CASCADE;"`
	Columns []*BoardColumn `json:"columns" gorm:"foreignkey:BoardID;constraint:OnDelete:CASCADE;"`
}
