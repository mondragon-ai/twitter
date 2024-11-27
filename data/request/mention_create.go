package request

type MentionCreateRequest struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Author  string `json:"author"`
	Created string `json:"created"` // Adjust the type if needed
}