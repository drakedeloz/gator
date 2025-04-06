-- name: CreateFeed :one
INSERT INTO feeds (id, user_id, name, url)
VALUES (
  gen_random_uuid(),
  $1,
  $2,
  $3
)
RETURNING *;