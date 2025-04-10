// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: feeds.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createFeed = `-- name: CreateFeed :one
INSERT INTO feeds (id, user_id, name, url)
VALUES (
  gen_random_uuid(),
  $1,
  $2,
  $3
)
RETURNING id, user_id, name, created_at, updated_at, url, last_fetched_at
`

type CreateFeedParams struct {
	UserID uuid.UUID
	Name   string
	Url    string
}

func (q *Queries) CreateFeed(ctx context.Context, arg CreateFeedParams) (Feed, error) {
	row := q.db.QueryRowContext(ctx, createFeed, arg.UserID, arg.Name, arg.Url)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Url,
		&i.LastFetchedAt,
	)
	return i, err
}

const createFeedFollow = `-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
  INSERT INTO feed_follows (id, user_id, feed_id)
  VALUES (
    gen_random_uuid(),
    $1,
    $2
  )
  RETURNING id, created_at, updated_at, user_id, feed_id
)
SELECT
  inserted_feed_follow.id, inserted_feed_follow.created_at, inserted_feed_follow.updated_at, inserted_feed_follow.user_id, inserted_feed_follow.feed_id,
  feeds.name AS feed_name,
  users.name AS user_name
FROM inserted_feed_follow
INNER JOIN users
  ON inserted_feed_follow.user_id = users.id
INNER JOIN feeds
  ON inserted_feed_follow.feed_id = feeds.id
`

type CreateFeedFollowParams struct {
	UserID uuid.UUID
	FeedID uuid.UUID
}

type CreateFeedFollowRow struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
	FeedName  string
	UserName  string
}

func (q *Queries) CreateFeedFollow(ctx context.Context, arg CreateFeedFollowParams) (CreateFeedFollowRow, error) {
	row := q.db.QueryRowContext(ctx, createFeedFollow, arg.UserID, arg.FeedID)
	var i CreateFeedFollowRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.FeedID,
		&i.FeedName,
		&i.UserName,
	)
	return i, err
}

const getAllFeeds = `-- name: GetAllFeeds :many
SELECT
  feeds.name AS feed_name,
  feeds.url,
  users.name as user_name
FROM
  feeds
JOIN
  users ON feeds.user_id = users.id
`

type GetAllFeedsRow struct {
	FeedName string
	Url      string
	UserName string
}

func (q *Queries) GetAllFeeds(ctx context.Context) ([]GetAllFeedsRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllFeeds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllFeedsRow
	for rows.Next() {
		var i GetAllFeedsRow
		if err := rows.Scan(&i.FeedName, &i.Url, &i.UserName); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFeedByUrl = `-- name: GetFeedByUrl :one
SELECT id, user_id, name, created_at, updated_at, url, last_fetched_at FROM feeds
WHERE url = $1
`

func (q *Queries) GetFeedByUrl(ctx context.Context, url string) (Feed, error) {
	row := q.db.QueryRowContext(ctx, getFeedByUrl, url)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Url,
		&i.LastFetchedAt,
	)
	return i, err
}

const getFeedFollowsByUser = `-- name: GetFeedFollowsByUser :many
SELECT
  feeds.name AS feed_name,
  users.name AS user_name
FROM feed_follows
INNER JOIN users
  ON feed_follows.user_id = users.id
INNER JOIN feeds
  ON feed_follows.feed_id = feeds.id
WHERE feed_follows.user_id = $1
`

type GetFeedFollowsByUserRow struct {
	FeedName string
	UserName string
}

func (q *Queries) GetFeedFollowsByUser(ctx context.Context, userID uuid.UUID) ([]GetFeedFollowsByUserRow, error) {
	rows, err := q.db.QueryContext(ctx, getFeedFollowsByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFeedFollowsByUserRow
	for rows.Next() {
		var i GetFeedFollowsByUserRow
		if err := rows.Scan(&i.FeedName, &i.UserName); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getNextFeedToFetch = `-- name: GetNextFeedToFetch :one
SELECT id, url, last_fetched_at
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1
`

type GetNextFeedToFetchRow struct {
	ID            uuid.UUID
	Url           string
	LastFetchedAt sql.NullTime
}

func (q *Queries) GetNextFeedToFetch(ctx context.Context) (GetNextFeedToFetchRow, error) {
	row := q.db.QueryRowContext(ctx, getNextFeedToFetch)
	var i GetNextFeedToFetchRow
	err := row.Scan(&i.ID, &i.Url, &i.LastFetchedAt)
	return i, err
}

const markFeedFetched = `-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = NOW(), updated_at = NOW()
WHERE id = $1
`

func (q *Queries) MarkFeedFetched(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, markFeedFetched, id)
	return err
}

const unfollowFeed = `-- name: UnfollowFeed :exec
DELETE FROM feed_follows
WHERE user_id = $1 AND feed_id = $2
`

type UnfollowFeedParams struct {
	UserID uuid.UUID
	FeedID uuid.UUID
}

func (q *Queries) UnfollowFeed(ctx context.Context, arg UnfollowFeedParams) error {
	_, err := q.db.ExecContext(ctx, unfollowFeed, arg.UserID, arg.FeedID)
	return err
}
