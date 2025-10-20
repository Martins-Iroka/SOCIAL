package store

import (
	"context"
	"database/sql"
)

type User struct {
	ID        int64  `json:"id"` //Json unmarshal
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  int64  `json:"-"` // this indicates that password won't be returned to the user.
	CreatedAt string `json:"created_at"`
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (username, password, email) VALUES ($1, $2, $3)
		RETURNING id, created_at 
	`

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
