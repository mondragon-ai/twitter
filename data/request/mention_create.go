package request

type MentionCreateRequest struct {
	Tweet string `json:"tweet"`
	ID int `json:"id"`
}