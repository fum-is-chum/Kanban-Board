package model

import "gorm.io/gorm"

type Board struct {
	gorm.Model
	Name    string         `json:"name,omitempty"`
	Desc    string         `json:"desc,omitempty"`
	OwnerID uint           `json:"owner_id,omitempty"`
	Members []*User        `json:"members" gorm:"many2many:board_members;constraint:OnDelete:CASCADE;"`
	Columns []*BoardColumn `json:"columns" gorm:"foreignkey:BoardID;constraint:OnDelete:CASCADE;"`
}
