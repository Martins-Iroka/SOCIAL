package store

import (
	"context"
	"database/sql"
	"time"
)

func NewMockStore() Storage {
	return Storage{
		User: &MockUserStore{},
	}
}

type MockUserStore struct {
}

func (s *MockUserStore) CreateUser(ctx context.Context, tx *sql.Tx, user *User) error {

	return nil
}

func (s *MockUserStore) CreateAndInviteUser(ctx context.Context, user *User, token string, invitationExp time.Duration) error {
	return nil
}

func (s *MockUserStore) GetUserByID(ctx context.Context, userID int64) (*User, error) {
	return nil, nil
}

func (s *MockUserStore) FollowUser(ctx context.Context, followerID int64, userID int64) error {
	return nil
}

func (s *MockUserStore) UnFollowUser(ctx context.Context, followerID int64, userID int64) error {
	return nil
}

func (s *MockUserStore) ActivateUser(ctx context.Context, token string) error {
	return nil
}

func (s *MockUserStore) DeleteUser(ctx context.Context, userID int64) error {
	return nil
}

func (s *MockUserStore) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return nil, nil
}
