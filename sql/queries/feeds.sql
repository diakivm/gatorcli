-- name: CreateFeed :one
INSERT INTO feeds (name, url, user_id)
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetFeeds :many
SELECT feeds.*, users.name as user_name FROM feeds JOIN users ON feeds.user_id = users.id;

-- name: GetFeed :one
SELECT * FROM feeds WHERE url = $1;