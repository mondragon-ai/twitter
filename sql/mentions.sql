
CREATE TABLE public.mentions (
    id VARCHAR(255) PRIMARY KEY,
    content TEXT NOT NULL,
    author VARCHAR(255) NOT NULL, 
    created TEXT NOT NULL
);