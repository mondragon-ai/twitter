CREATE TABLE mention (
    id SERIAL PRIMARY KEY,
    parent_id TEXT,                 -- Conversation ID (Parent Tweet ID)
    author_id TEXT NOT NULL,        -- Author ID
    tweet_id TEXT NOT NULL UNIQUE,  -- Tweet ID (Unique identifier for the tweet)
    content TEXT NOT NULL,          -- Tweet content
    author_name TEXT NOT NULL,      -- Author's name
    created_at TEXT NOT NULL        -- Timestamp for when the mention was recorded
);