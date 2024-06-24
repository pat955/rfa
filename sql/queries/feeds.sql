-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: RetrieveFeeds :many
SELECT * FROM feeds;

-- name: DeleteFeed :exec
DELETE FROM feeds
WHERE id = $1;

-- name: GetFeed :one
SELECT * FROM feeds
WHERE id = $1;

-- name: GetNextFeedsToFetch :one
SELECT * FROM feeds
ORDER BY last_fetched_at
LIMIT $1;