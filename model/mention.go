package model


type Mention struct {
	ID				string			`json:"main_id"`
	ParentID		string			`json:"conversation_id"`
	AuthorID		string      	`json:"author_id"`
	TweetID			string      	`json:"id"`
	Content			string      	`json:"text"`
	ParentContent	string      	`json:"parent_text"`
	AuthorName		string      	`json:"author_name"`
	CreatedAt   	string 			`json:"created_at"` 
}

type ArticleUrl struct {
	ID		int				`json:"id"`
	Url		string			`json:"url"`
	Title	string      	`json:"title"`
}

type TweetClone struct {
	ID			int			`json:"id"`
	AuthorName	string		`json:"author_name,omitempty"`
	Tweet		string      `json:"tweet"`
}

type ThreadIdea struct {
	ID			int		`json:"id"`
	Idea		string		`json:"idea,omitempty"`
	UsedCount	int      	`json:"used_count"`
}

type TweetIdea struct {
	ID			int		`json:"id"`
	Idea		string		`json:"idea,omitempty"`
	UsedCount	int      	`json:"used_count"`
}