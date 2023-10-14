package model

import (
	"time"
)

type BoardMember struct {
	BoardID uint `json:"board_id" gorm:"primaryKey;autoIncrement:false"`
	UserID  uint `json:"user_id" gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}