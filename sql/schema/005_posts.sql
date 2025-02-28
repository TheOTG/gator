-- +goose Up
CREATE TABLE posts (
    id UUID,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE,
    description TEXT,
    published_at TIMESTAMP,
    feed_id UUID NOT NULL REFERENCES feeds ON DELETE CASCADE,
    PRIMARY KEY(id)
);

-- +goose Down
DROP TABLE posts;