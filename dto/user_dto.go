package dto

import "time"

type UserRequest struct {
	Name     string `json:"name,omitempty" validate:"required"`
	Email    string `json:"email,omitempty" validate:"required,email"`
	Password string `json:"password,omitempty" validate:"required"`
}

type UserResponse struct {
	Id        uint             `json:"id"`
	Name      string           `json:"name"`
	Email     string           `json:"email"`
	MemberOf  []*BoardResponse `json:"member_of"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}
