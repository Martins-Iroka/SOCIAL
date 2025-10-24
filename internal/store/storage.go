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
	QueryTimeoutDuration      = time.Second * 5
)

// This is repository pattern implementation
type Storage struct {
	Post interface {
		Create(context.Context, *Post) error
		GetByID(context.Context, int64) (*Post, error)
		GetUserFeed(context.Context, int64, int16, int16, string) ([]PostWithMetadata, error)
		Delete(context.Context, int64) error
		Update(context.Context, *Post) error
	}
	User interface {
		CreateUser(context.Context, *User) error
		GetUserByID(context.Context, int64) (*User, error)
		FollowUser(context.Context, int64, int64) error
		UnFollowUser(context.Context, int64, int64) error
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
