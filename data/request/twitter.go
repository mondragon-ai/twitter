package request

type TweetCreateRequest struct {
	Type  string `json:"type"`
}


// Mention represents a mentioned user in the tweet
type Mention struct {
	Start    int    `json:"start"`
	End      int    `json:"end"`
	Username string `json:"username"`
	ID       string `json:"id"`
}

// Entities represents the entities in the tweet
type Entities struct {
	Mentions []Mention `json:"mentions"`
}


// ReferencedTweet represents a tweet that is referenced (replied to or quoted)
type ReferencedTweet struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// Attachments represents attachments in a tweet
type Attachments struct {
	MediaKeys []string `json:"media_keys"`
}


type Tweet struct {
	CreatedAt           string			    `json:"created_at"`
	AuthorID            string              `json:"author_id"`
	Text                string              `json:"text"`
	Entities            Entities            `json:"entities"`
	EditHistoryTweetIDs []string            `json:"edit_history_tweet_ids"`
	ID                  string              `json:"id"`
	ConversationID      string              `json:"conversation_id"`
	ReferencedTweets    *[]ReferencedTweet  `json:"referenced_tweets,omitempty"`
	InReplyToUserID     *string             `json:"in_reply_to_user_id,omitempty"`
	Attachments         *Attachments        `json:"attachments,omitempty"`
}


type Includes struct {
	Users []struct {
		ID              string `json:"id"`
		Username        string `json:"username"`
		Name            string `json:"name"`
		ProfileImageURL string `json:"profile_image_url"`
	} `json:"users"`
	Tweets []struct {
		ID        string `json:"id"`
		Text      string `json:"text"`
		AuthorID  string `json:"author_id"`
		CreatedAt string `json:"created_at"`
	} `json:"tweets"`
	Media []struct {
		MediaKey string `json:"media_key"`
		Type     string `json:"type"`
		URL      string `json:"url"`
	} `json:"media"`
}

type RootTweetMentions struct {
	Data 		[]Tweet 	`json:"data"`
	Includes	Includes	`json:"includes"`
}

type UrlCreate struct {
	Url		string			`json:"url"`
	Title	string      	`json:"title"`
}

type TweetClone struct {
	AuthorName	string		`json:"author_name,omitempty"`
	Tweet		string      `json:"title"`
}

type ThreadIdea struct {
	Idea		string		`json:"idea,omitempty"`
	UsedCount	int      	`json:"used_count"`
}

type TweetIdea struct {
	Idea		string		`json:"idea,omitempty"`
	UsedCount	int      	`json:"used_count"`
}