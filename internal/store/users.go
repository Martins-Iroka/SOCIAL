package store

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/lib/pq"
)

type User struct {
	ID        int64  `json:"id"` //Json unmarshal
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"` // - indicates that password won't be returned to the user upon calling the endpoint.
	CreatedAt string `json:"created_at"`
}

type Follower struct {
	UserID     int64 `json:"user_id"`
	FollowerID int64 `json:"follower_id"`
	CreatedAt  int64 `json:"created_at"`
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

func (s *UserStore) FollowUser(ctx context.Context, followerID int64, userID int64) error {
	query := `INSERT INTO followers (user_id, follower_id) VALUES ($1, $2)`

	c, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(c, query, userID, followerID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrorUserFollowConflict
		}
		return err
	}
	return err
}

func (s *UserStore) UnFollowUser(ctx context.Context, followerID int64, userID int64) error {
	query := `DELETE FROM followers WHERE user_id = $1 AND follower_id = $2`

	c, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := s.db.ExecContext(c, query, userID, followerID)
	log.Println(err)
	if err != nil {
		return err
	}

	row, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if row == 0 {
		return ErrorNotFound
	}

	return nil
}
