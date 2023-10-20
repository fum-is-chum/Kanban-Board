package dto

import "time"

type BoardRequest struct {
	Name    string `json:"name,omitempty" validate:"required"`
	Desc    string `json:"desc,omitempty" validate:"required"`
	OwnerID uint   `json:"owner_id,omitempty" validate:"required"`
}

type BoardResponse struct {
	ID      uint                   `json:"id"`
	Name    string                 `json:"name"`
	Desc    string                 `json:"desc"`
	Owner   *MemberResponse        `json:"owner,omitempty"`
	Members []*MemberResponse      `json:"members,omitempty"`
	Columns []*BoardColumnResponse `json:"columns,omitempty"`
}

type MemberResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
