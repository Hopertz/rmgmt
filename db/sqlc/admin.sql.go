// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: admin.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createAdmin = `-- name: CreateAdmin :one
INSERT INTO admin (email, password_hash, activated)
VALUES ($1, $2, $3 )
RETURNING id, created_at, version
`

type CreateAdminParams struct {
	Email        string `json:"email"`
	PasswordHash []byte `json:"password_hash"`
	Activated    bool   `json:"activated"`
}

type CreateAdminRow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Version   uuid.UUID `json:"version"`
}

func (q *Queries) CreateAdmin(ctx context.Context, arg CreateAdminParams) (CreateAdminRow, error) {
	row := q.db.QueryRowContext(ctx, createAdmin, arg.Email, arg.PasswordHash, arg.Activated)
	var i CreateAdminRow
	err := row.Scan(&i.ID, &i.CreatedAt, &i.Version)
	return i, err
}

const getAdminByEmail = `-- name: GetAdminByEmail :one
SELECT id, created_at, email, password_hash, activated, version
FROM admin
WHERE email = $1
`

func (q *Queries) GetAdminByEmail(ctx context.Context, email string) (Admin, error) {
	row := q.db.QueryRowContext(ctx, getAdminByEmail, email)
	var i Admin
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Email,
		&i.PasswordHash,
		&i.Activated,
		&i.Version,
	)
	return i, err
}

const getHashTokenForAdmin = `-- name: GetHashTokenForAdmin :one
SELECT admin.id, admin.created_at,admin.email, admin.password_hash,admin.version, admin.activated
FROM admin
INNER JOIN token
ON admin.id = tokens.id
WHERE token.hash = $1
AND token.scope = $2
AND token.expiry > $3
`

type GetHashTokenForAdminParams struct {
	Hash   []byte    `json:"hash"`
	Scope  string    `json:"scope"`
	Expiry time.Time `json:"expiry"`
}

type GetHashTokenForAdminRow struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	Email        string    `json:"email"`
	PasswordHash []byte    `json:"password_hash"`
	Version      uuid.UUID `json:"version"`
	Activated    bool      `json:"activated"`
}

func (q *Queries) GetHashTokenForAdmin(ctx context.Context, arg GetHashTokenForAdminParams) (GetHashTokenForAdminRow, error) {
	row := q.db.QueryRowContext(ctx, getHashTokenForAdmin, arg.Hash, arg.Scope, arg.Expiry)
	var i GetHashTokenForAdminRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Email,
		&i.PasswordHash,
		&i.Version,
		&i.Activated,
	)
	return i, err
}

const updateAdmin = `-- name: UpdateAdmin :one
UPDATE admin
SET email = $1, password_hash = $2, activated = $3, version = uuid_generate_v4()
WHERE id = $4 AND version = $5
RETURNING version
`

type UpdateAdminParams struct {
	Email        string    `json:"email"`
	PasswordHash []byte    `json:"password_hash"`
	Activated    bool      `json:"activated"`
	ID           uuid.UUID `json:"id"`
	Version      uuid.UUID `json:"version"`
}

func (q *Queries) UpdateAdmin(ctx context.Context, arg UpdateAdminParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, updateAdmin,
		arg.Email,
		arg.PasswordHash,
		arg.Activated,
		arg.ID,
		arg.Version,
	)
	var version uuid.UUID
	err := row.Scan(&version)
	return version, err
}
