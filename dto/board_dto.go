package dto

import "time"

type BoardRequest struct {
	Name    string `json:"name,omitifempty" validate:"required"`
	Desc    string `json:"desc,omitifempty" validate:"required"`
	OwnerID uint   `json:"owner_id,omitifempty" validate:"required"`
}

type BoardResponse struct {
	Id      uint              `json:"id"`
	Name    string            `json:"name"`
	Desc    string            `json:"desc"`
	Owner   *MemberResponse   `json:"owner,omitempty"`
	Members []*MemberResponse `json:"members,omitempty"`
}

type MemberResponse struct {
	Id        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
