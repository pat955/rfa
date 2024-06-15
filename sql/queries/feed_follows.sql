-- name: AddFeedFollow :one
INSERT INTO feed_follows (id, feed_id, user_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteFeedFollow :exec
DELETE FROM feeds
WHERE id = $1;

-- name: GetFeedFollow :one
SELECT * FROM feed_follows
WHERE id = $1;

-- name: GetAllFeedFollows :one
SELECT * FROM feed_follows;