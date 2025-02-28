-- +goose Up
CREATE TABLE users (
    id UUID,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL UNIQUE,
    PRIMARY KEY(id)
);

-- +goose Down
DROP TABLE users;