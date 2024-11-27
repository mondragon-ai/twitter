package response

type MentionResponse struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Author  string `json:"author"`
	Created string `json:"created"` 
}