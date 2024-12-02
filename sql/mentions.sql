CREATE TABLE mention (
    id SERIAL PRIMARY KEY,
    parent_id TEXT,                 -- Conversation ID (Parent Tweet ID)
    author_id TEXT NOT NULL,        -- Author ID
    tweet_id TEXT NOT NULL UNIQUE,  -- Tweet ID (Unique identifier for the tweet)
    content TEXT NOT NULL,          -- Tweet content
    author_name TEXT NOT NULL,      -- Author's name
    created_at TEXT NOT NULL        -- Timestamp for when the mention was recorded
);


-- Table for ArticleUrl
CREATE TABLE article_url (
    id SERIAL PRIMARY KEY,          -- Unique identifier for each article URL
    url TEXT NOT NULL,              -- URL of the article
    title TEXT NOT NULL             -- Title of the article
);

-- Table for TweetClone
CREATE TABLE tweet_clone (
    id SERIAL PRIMARY KEY,          -- Unique identifier for each tweet clone
    author_name TEXT,               -- Author's name (optional)
    tweet TEXT NOT NULL             -- Tweet content
);

-- Table for ThreadIdea
CREATE TABLE thread_idea (
    id SERIAL PRIMARY KEY,          -- Unique identifier for each thread idea
    idea TEXT,                      -- Idea content (optional)
    used_count INT NOT NULL         -- Number of times the idea has been used
);

-- Table for TweetIdea
CREATE TABLE tweet_idea (
    id SERIAL PRIMARY KEY,          -- Unique identifier for each tweet idea
    idea TEXT,                      -- Idea content (optional)
    used_count INT NOT NULL         -- Number of times the idea has been used
);