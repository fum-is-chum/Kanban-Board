package model

import "gorm.io/gorm"

type BoardColumn struct {
	gorm.Model
	Label   string  `json:"label"`
	Desc    string  `json:"desc"`
	BoardID uint    `json:"board_id"`
	Board   *Board  `json:"board"`
	Tasks   []*Task `json:"tasks" gorm:"foreignKey:BoardColumnID;constraint:OnDelete:CASCADE;"`
}
