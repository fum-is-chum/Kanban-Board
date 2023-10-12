package dto

type BoardRequest struct {
	Name    string `json:"name" validate:"required"`
	Desc    string `json:"desc" validate:"required"`
	OwnerID uint   `json:"owner_id" validate:"required"`
}

type BoardResponse struct {
	Id      uint   `json:"id"`
	Name    string `json:"name"`
	Desc    string `json:"desc"`
	OwnerID uint   `json:"owner_id"`
}
