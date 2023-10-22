package model

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Title         string `json:"title"`
	Desc          string `json:"desc"`
	BoardID       uint   `json:"board_id"`
	BoardColumnID uint   `json:"board_column_id"`
}
