-- name: CreateFeed :one
INSERT INTO feeds (id, user_id, name, url)
VALUES (
  gen_random_uuid(),
  $1,
  $2,
  $3
)
RETURNING *;

-- name: GetAllFeeds :many
SELECT
  feeds.name AS feed_name,
  feeds.url,
  users.name as user_name
FROM
  feeds
JOIN
  users ON feeds.user_id = users.id;

-- name: GetFeedByUrl :one
SELECT * FROM feeds
WHERE url = $1;

-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
  INSERT INTO feed_follows (id, user_id, feed_id)
  VALUES (
    gen_random_uuid(),
    $1,
    $2
  )
  RETURNING *
)
SELECT
  inserted_feed_follow.*,
  feeds.name AS feed_name,
  users.name AS user_name
FROM inserted_feed_follow
INNER JOIN users
  ON inserted_feed_follow.user_id = users.id
INNER JOIN feeds
  ON inserted_feed_follow.feed_id = feeds.id;

-- name: GetFeedFollowsByUser :many
SELECT
  feeds.name AS feed_name,
  users.name AS user_name
FROM feed_follows
INNER JOIN users
  ON feed_follows.user_id = users.id
INNER JOIN feeds
  ON feed_follows.feed_id = feeds.id
WHERE feed_follows.user_id = $1;

-- name: UnfollowFeed :exec
DELETE FROM feed_follows
WHERE user_id = $1 AND feed_id = $2;