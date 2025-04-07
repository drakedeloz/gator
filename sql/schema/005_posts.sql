-- +goose Up
CREATE TABLE posts (
  id UUID NOT NULL PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  title VARCHAR(255) NOT NULL,
  url VARCHAR(2048) NOT NULL UNIQUE,
  description TEXT NOT NULL,
  published_at TIMESTAMP NOT NULL,
  feed_id UUID NOT NULl,
  FOREIGN KEY (feed_id) REFERENCES feeds(id)
);

-- +goose Down
DROP TABLE posts;