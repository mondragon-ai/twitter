package model


type Mention struct {
	ParentID	string			`json:"conversation_id"`
	AuthorID	string      	`json:"author_id"`
	TweetID		string      	`json:"id"`
	Content		string      	`json:"text"`
	AuthorName	string      	`json:"author_name"`
	CreatedAt   string 			`json:"created_at"` 
}
