package request

type PostCreateRequest struct {
	Title       string `validate:"required min=1,max=100" json:"title"`
	Published   bool   `json:"published"`
	Description string `json:"description"`
}
