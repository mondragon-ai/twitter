CREATE TABLE mention (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    author VARCHAR(255),
    created TEXT NOT NULL
);