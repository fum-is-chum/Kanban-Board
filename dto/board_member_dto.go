package dto

import "time"

type BoardMemberRequest struct {
	BoardID uint `json:"board_id,omitempty" validate:"required"`
	UserID  uint `json:"user_id,omitempty" validate:"required"`
}

type BoardMemberResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

