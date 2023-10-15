package dto

type BoardMemberRequest struct {
	BoardID uint `json:"board_id" validate:"required"`
	UserID  uint `json:"user_id" validate:"required"`
}
