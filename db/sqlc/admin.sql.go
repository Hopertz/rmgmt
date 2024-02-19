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

const getAdminByEmail = `-- name: GetAdminByEmail :one
SELECT admin_id, created_at, email, password_hash, activated, version
FROM admins
WHERE email = $1
`

func (q *Queries) GetAdminByEmail(ctx context.Context, email string) (Admin, error) {
	row := q.db.QueryRowContext(ctx, getAdminByEmail, email)
	var i Admin
	err := row.Scan(
		&i.AdminID,
		&i.CreatedAt,
		&i.Email,
		&i.PasswordHash,
		&i.Activated,
		&i.Version,
	)
	return i, err
}

const getForTokenAdmin = `-- name: GetForTokenAdmin :one
SELECT admins.admin_id, admins.created_at,admins.email, admins.password_hash, admins.activated, admins.version
FROM admins
INNER JOIN tokens
ON admins.admin_id = tokens.admin_id
WHERE tokens.hash = $1
AND tokens.scope = $2
AND tokens.expiry > $3
`

type GetForTokenAdminParams struct {
	Hash   []byte    `json:"hash"`
	Scope  string    `json:"scope"`
	Expiry time.Time `json:"expiry"`
}

func (q *Queries) GetForTokenAdmin(ctx context.Context, arg GetForTokenAdminParams) (Admin, error) {
	row := q.db.QueryRowContext(ctx, getForTokenAdmin, arg.Hash, arg.Scope, arg.Expiry)
	var i Admin
	err := row.Scan(
		&i.AdminID,
		&i.CreatedAt,
		&i.Email,
		&i.PasswordHash,
		&i.Activated,
		&i.Version,
	)
	return i, err
}

const insertAdmin = `-- name: InsertAdmin :one
INSERT INTO admins (email, password_hash, activated)
VALUES ($1, $2, $3 )
RETURNING admin_id, created_at, version
`

type InsertAdminParams struct {
	Email        string `json:"email"`
	PasswordHash []byte `json:"password_hash"`
	Activated    bool   `json:"activated"`
}

type InsertAdminRow struct {
	AdminID   uuid.UUID `json:"admin_id"`
	CreatedAt time.Time `json:"created_at"`
	Version   int32     `json:"version"`
}

func (q *Queries) InsertAdmin(ctx context.Context, arg InsertAdminParams) (InsertAdminRow, error) {
	row := q.db.QueryRowContext(ctx, insertAdmin, arg.Email, arg.PasswordHash, arg.Activated)
	var i InsertAdminRow
	err := row.Scan(&i.AdminID, &i.CreatedAt, &i.Version)
	return i, err
}

const updateAdmin = `-- name: UpdateAdmin :one
UPDATE admins
SET email = $1, password_hash = $2, activated = $3, version = version + 1
WHERE admin_id = $4 AND version = $5
RETURNING version
`

type UpdateAdminParams struct {
	Email        string    `json:"email"`
	PasswordHash []byte    `json:"password_hash"`
	Activated    bool      `json:"activated"`
	AdminID      uuid.UUID `json:"admin_id"`
	Version      int32     `json:"version"`
}

func (q *Queries) UpdateAdmin(ctx context.Context, arg UpdateAdminParams) (int32, error) {
	row := q.db.QueryRowContext(ctx, updateAdmin,
		arg.Email,
		arg.PasswordHash,
		arg.Activated,
		arg.AdminID,
		arg.Version,
	)
	var version int32
	err := row.Scan(&version)
	return version, err
}