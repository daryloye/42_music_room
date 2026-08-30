package response

type PostResponse struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Published   bool   `json:"published"`
	Description string `json:"description"`
}
