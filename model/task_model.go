package model

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Title         string `json:"title"`
	Desc          string `json:"desc"`
	BoardColumnID uint   `json:"board_column_id"`
	BoardID       uint   `json:"board_id"`
}
