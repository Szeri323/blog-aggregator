// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: feeds.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createFeeds = `-- name: CreateFeeds :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES(
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
 ) RETURNING id, created_at, updated_at, name, url, user_id
`

type CreateFeedsParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Url       string
	UserID    uuid.UUID
}

func (q *Queries) CreateFeeds(ctx context.Context, arg CreateFeedsParams) (Feed, error) {
	row := q.db.QueryRowContext(ctx, createFeeds,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Name,
		arg.Url,
		arg.UserID,
	)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.UserID,
	)
	return i, err
}

const getFeed = `-- name: GetFeed :one
SELECT id, created_at, updated_at, name, url, user_id FROM feeds WHERE url = $1
`

func (q *Queries) GetFeed(ctx context.Context, url string) (Feed, error) {
	row := q.db.QueryRowContext(ctx, getFeed, url)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.UserID,
	)
	return i, err
}

const getFeeds = `-- name: GetFeeds :many
SELECT feeds.name,url,users.name 
FROM feeds 
INNER JOIN users ON feeds.user_id = users.id
`

type GetFeedsRow struct {
	Name   string
	Url    string
	Name_2 string
}

func (q *Queries) GetFeeds(ctx context.Context) ([]GetFeedsRow, error) {
	rows, err := q.db.QueryContext(ctx, getFeeds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFeedsRow
	for rows.Next() {
		var i GetFeedsRow
		if err := rows.Scan(&i.Name, &i.Url, &i.Name_2); err != nil {
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

const truncateFeedsTable = `-- name: TruncateFeedsTable :exec
TRUNCATE TABLE feeds CASCADE
`

func (q *Queries) TruncateFeedsTable(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, truncateFeedsTable)
	return err
}
