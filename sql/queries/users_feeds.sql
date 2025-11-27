-- name: CreateUsersFeed :one
WITH new_users_feed AS (
    INSERT INTO users_feeds (user_id, feed_id)
    VALUES ($1, $2)
    ON CONFLICT DO NOTHING
    RETURNING *
)
SELECT new_users_feed.id, users.name as user_name, users.id as user_id, feeds.name as feed_name, feeds.id as feed_id, feeds.url as feed_url 
FROM new_users_feed 
JOIN users ON users.id = new_users_feed.user_id
JOIN feeds ON feeds.id = new_users_feed.feed_id;

-- name: GetUsersFeeds :many
SELECT users_feeds.id, users.name as user_name, users.id as user_id, feeds.name as feed_name, feeds.id as feed_id, feeds.url as feed_url 
FROM users_feeds 
JOIN users ON users.id = users_feeds.user_id
JOIN feeds ON feeds.id = users_feeds.feed_id
WHERE users.id = $1;

-- name: RemoveUsersFeed :exec
DELETE FROM users_feeds WHERE user_id = $1 AND feed_id = $2;