package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Martins-Iroka/social/internal/store"
	"github.com/go-redis/redis/v8"
)

type UserStore struct {
	rdb *redis.Client
}

// Redis is a key:value storage
func (s *UserStore) Get(ctx context.Context, userID int64) (*store.User, error) {
	cacheKey := fmt.Sprintf("user-%v", userID)
	data, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var user store.User
	if data != "" {
		if err := json.Unmarshal([]byte(data), &user); err != nil {
			return nil, err
		}
	}
	return &user, nil
}

// TTL (Time To Leave)
func (s *UserStore) Set(ctx context.Context, user *store.User) error {
	if user.ID <= 0 {
		return store.ErrorNotFound
	}
	cacheKey := fmt.Sprintf("user-%v", user.ID)

	json, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return s.rdb.SetEX(ctx, cacheKey, json, 5*time.Minute).Err()
}
