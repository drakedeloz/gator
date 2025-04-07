-- name: CreatePost :exec
INSERT INTO posts (id, title, url, description, published_at, feed_id)
VALUES (
  gen_random_uuid(),
  $1,
  $2,
  $3,
  $4,
  $5
)
RETURNING *;

-- name: GetPostsForUser :many
SELECT
  posts.*
FROM posts
INNER JOIN feed_follows
  ON feed_follows.feed_id = posts.feed_id
WHERE user_id = $1
ORDER BY posts.published_at DESC
LIMIT $2;