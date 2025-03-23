-- name: CreateFeeds :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES(
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
 ) RETURNING *;

-- name: TruncateFeedsTable :exec
TRUNCATE TABLE feeds CASCADE;

-- name: GetFeeds :many
SELECT feeds.name,url,users.name 
FROM feeds 
INNER JOIN users ON feeds.user_id = users.id;

-- name: GetFeed :one
SELECT * FROM feeds WHERE url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds 
SET updated_at = $1, last_fetched_at = $2 
WHERE id = $3;

-- name: GetNextFeedToFetch :one 
SELECT id, url FROM feeds ORDER BY last_fetched_at NULLS FIRST LIMIT 1;