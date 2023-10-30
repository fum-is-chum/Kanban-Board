package dto

type BoardRequest struct {
	Name string `json:"name,omitempty" validate:"required"`
	Desc string `json:"desc,omitempty" validate:"required"`
}

type BoardResponse struct {
	ID      uint                   `json:"id"`
	Name    string                 `json:"name"`
	Desc    string                 `json:"desc"`
	Owner   *BoardMemberResponse   `json:"owner,omitempty"`
	Members []*BoardMemberResponse `json:"members"`
	Columns []*BoardColumnResponse `json:"columns"`
}
