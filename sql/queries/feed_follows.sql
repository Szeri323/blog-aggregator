-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (INSERT INTO feed_follows(id, created_at, updated_at, user_id, feed_id) 
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
) RETURNING *)
SELECT 
    inserted_feed_follow.*,
    users.name AS user_name, 
    feeds.name AS feed_name
FROM inserted_feed_follow
    INNER JOIN users ON inserted_feed_follow.user_id = users.id 
    INNER JOIN feeds ON inserted_feed_follow.feed_id = feeds.id;

-- name: TruncateFeedFollowsTable :exec
TRUNCATE TABLE feed_follows CASCADE;

-- name: GetFeedFollowsForUser :many
SELECT feed_follows.*, 
    feeds.name AS feed_name 
FROM feed_follows 
INNER JOIN feeds ON feed_follows.feed_id = feeds.id
WHERE feed_follows.user_id = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows WHERE feed_id = (SELECT id FROM feeds WHERE url = $1)
AND feed_follows.user_id = $2;