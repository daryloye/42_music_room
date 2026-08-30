package request

type PostUpdateRequest struct {
	Id          string
	Title       string `validate:"required min=1,max=100" json:"title"`
	Published   bool   `json:"published"`
	Description string `json:"description"`
}
