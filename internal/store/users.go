package store

import (
	"context"
	"database/sql"
	"errors"
)

type User struct {
	ID        int64  `json:"id"` //Json unmarshal
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"` // - indicates that password won't be returned to the user upon calling the endpoint.
	CreatedAt string `json:"created_at"`
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) CreateUser(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (username, password, email) VALUES ($1, $2, $3)
		RETURNING id, created_at 
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	if err := s.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Password,
		user.Email,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	); err != nil {
		return err
	}

	return nil
}

func (s *UserStore) GetUserByID(ctx context.Context, userID int64) (*User, error) {
	query := `
		SELECT id, username, email FROM users WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var user User

	err := s.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrorNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
