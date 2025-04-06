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