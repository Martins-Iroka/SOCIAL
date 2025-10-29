package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrorNotFound             = errors.New("resource not found")
	ErrorConflict             = errors.New("conflict found modifying resource")
	ErrorUserFollowConflict   = errors.New("you're following this user already")
	ErrorUserUnFollowConflict = errors.New("you're unfollowing this user already")
	ErrorDuplicateEmail       = errors.New("a user with that email already exists")
	ErrorDuplicateUsername    = errors.New("a user with that username already exists")
	QueryTimeoutDuration      = time.Second * 5
)

// This is repository pattern implementation
type Storage struct {
	Post interface {
		Create(context.Context, *Post) error
		GetByID(context.Context, int64) (*Post, error)
		GetUserFeed(
			context.Context, int64, *PaginatedFeedQuery) ([]PostWithMetadata, error)
		Delete(context.Context, int64) error
		Update(context.Context, *Post) error
	}
	User interface {
		ActivateUser(ctx context.Context, token string) error
		CreateUser(context.Context, *sql.Tx, *User) error
		CreateAndInviteUser(ctx context.Context, user *User, token string, time time.Duration) error
		GetUserByID(context.Context, int64) (*User, error)
		FollowUser(context.Context, int64, int64) error
		UnFollowUser(context.Context, int64, int64) error
		DeleteUser(context.Context, int64) error
		GetUserByEmail(context.Context, string) (*User, error)
	}
	Comment interface {
		CreateComment(context.Context, *Comment) error
		GetByPostID(ctx context.Context, postID int64) ([]Comment, error)
	}
}

func NewPostgresStorage(db *sql.DB) Storage {
	return Storage{
		Post:    &PostStore{db: db},
		User:    &UserStore{db: db},
		Comment: &CommentStore{db: db},
	}
}

func withTransaction(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()

}
