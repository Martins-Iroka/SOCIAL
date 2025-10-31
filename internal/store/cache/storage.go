package cache

import (
	"context"

	"github.com/Martins-Iroka/social/internal/store"
	"github.com/go-redis/redis/v8"
)

type Storage struct {
	User interface {
		Get(context.Context, int64) (*store.User, error)
		Set(context.Context, *store.User) error
	}
}

func NewRedisStore(rdb *redis.Client) Storage {
	return Storage{
		User: &UserStore{rdb},
	}
}
