package dto

type BoardColumnRequest struct {
	Label   string `json:"label,omitempty" validate:"required"`
	Desc    string `json:"desc,omitempty" validate:"required"`
	BoardID uint   `json:"board_id,omitempty" validate:"required"`
}

type BoardColumnResponse struct {
	ID      uint   `json:"id"`
	Label   string `json:"label"`
	Desc    string `json:"desc"`
	BoardID uint   `json:"board_id"`
}
