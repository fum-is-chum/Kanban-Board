package dto

type TaskCreateRequest struct {
	Title         string `json:"title,omitempty" validate:"required"`
	Desc          string `json:"desc,omitempty" validate:"required"`
	BoardColumnID uint   `json:"board_column_id,omitempty" validate:"required"`
	BoardID       uint   `json:"board_id,omitempty" validate:"required"`
}

type TaskUpdateRequest struct {
	Title         string `json:"title,omitempty"`
	Desc          string `json:"desc,omitempty"`
	BoardColumnID uint   `json:"board_column_id,omitempty"`
}

type TaskResponse struct {
	ID            uint   `json:"id"`
	Title         string `json:"title"`
	Desc          string `json:"desc"`
	BoardColumnID uint   `json:"board_column_id,omitempty"`
}
