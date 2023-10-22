package dto

import "time"

type BoardRequest struct {
	Name string `json:"name,omitempty" validate:"required"`
	Desc string `json:"desc,omitempty" validate:"required"`
}

type BoardResponse struct {
	ID      uint                   `json:"id"`
	Name    string                 `json:"name"`
	Desc    string                 `json:"desc"`
	Members []*MemberResponse      `json:"members"`
	Columns []*BoardColumnResponse `json:"columns"`
}

type MemberResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
