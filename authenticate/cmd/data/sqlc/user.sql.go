// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0
// source: user.sql

package db

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    email,
    hashed_password
) VALUES (
    $1, 
    $2
) RETURNING user_id, email, hashed_password, active, created_at, updated_at
`

type CreateUserParams struct {
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Email, arg.HashedPassword)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Email,
		&i.HashedPassword,
		&i.Active,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT user_id, email, hashed_password, active, created_at, updated_at FROM users
WHERE email = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, email)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Email,
		&i.HashedPassword,
		&i.Active,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
