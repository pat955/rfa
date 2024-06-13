-- +goose Up
CREATE TABLE feeds (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    name TEXT,
    url TEXT UNIQUE NOT NULL,
    user_id UUID ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds;
