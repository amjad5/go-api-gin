// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: query.sql

package user

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
  name, phone_number
) VALUES (
  $1, $2
)
ON CONFLICT (phone_number) DO NOTHING
RETURNING id, name, phone_number, otp, otp_expiration_time
`

type CreateUserParams struct {
	Name        string
	PhoneNumber string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser, arg.Name, arg.PhoneNumber)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.PhoneNumber,
		&i.Otp,
		&i.OtpExpirationTime,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteUser, id)
	return err
}

const getPhonenumber = `-- name: GetPhonenumber :one
SELECT phone_number from users
where phone_number = $1 LIMIT 1
`

func (q *Queries) GetPhonenumber(ctx context.Context, phoneNumber string) (string, error) {
	row := q.db.QueryRow(ctx, getPhonenumber, phoneNumber)
	var phone_number string
	err := row.Scan(&phone_number)
	return phone_number, err
}

const getUser = `-- name: GetUser :one
SELECT id, name, phone_number, otp, otp_expiration_time FROM users
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, id int32) (User, error) {
	row := q.db.QueryRow(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.PhoneNumber,
		&i.Otp,
		&i.OtpExpirationTime,
	)
	return i, err
}

const listUser = `-- name: ListUser :many
SELECT id, name, phone_number, otp, otp_expiration_time FROM users
ORDER BY name
`

func (q *Queries) ListUser(ctx context.Context) ([]User, error) {
	rows, err := q.db.Query(ctx, listUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.PhoneNumber,
			&i.Otp,
			&i.OtpExpirationTime,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateOtp = `-- name: UpdateOtp :exec
UPDATE users
  set otp = $2,
  otp_expiration_time = $3
WHERE phone_number = $1
`

type UpdateOtpParams struct {
	PhoneNumber       string
	Otp               pgtype.Text
	OtpExpirationTime pgtype.Timestamp
}

func (q *Queries) UpdateOtp(ctx context.Context, arg UpdateOtpParams) error {
	_, err := q.db.Exec(ctx, updateOtp, arg.PhoneNumber, arg.Otp, arg.OtpExpirationTime)
	return err
}

const updateUser = `-- name: UpdateUser :exec
UPDATE users
  set name = $2,
  phone_number = $3
WHERE id = $1
`

type UpdateUserParams struct {
	ID          int32
	Name        string
	PhoneNumber string
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.Exec(ctx, updateUser, arg.ID, arg.Name, arg.PhoneNumber)
	return err
}

const verifyOtp = `-- name: VerifyOtp :one
SELECT phone_number, otp_expiration_time from users
where phone_number = $1 and otp = $2 LIMIT 1
`

type VerifyOtpParams struct {
	PhoneNumber string
	Otp         pgtype.Text
}

type VerifyOtpRow struct {
	PhoneNumber       string
	OtpExpirationTime pgtype.Timestamp
}

func (q *Queries) VerifyOtp(ctx context.Context, arg VerifyOtpParams) (VerifyOtpRow, error) {
	row := q.db.QueryRow(ctx, verifyOtp, arg.PhoneNumber, arg.Otp)
	var i VerifyOtpRow
	err := row.Scan(&i.PhoneNumber, &i.OtpExpirationTime)
	return i, err
}
