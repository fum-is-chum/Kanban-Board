package dto

type BoardMemberRequest struct {
	BoardId uint `json:"board_id" validate:"required"`
	UserId  uint `json:"user_id" validate:"required"`
}
