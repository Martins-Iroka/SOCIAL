package store

import (
	"context"
	"database/sql"
)

type Comment struct {
	ID       int64  `json:"id"`
	PostID   int64  `json:"post_id"`
	UserID   int64  `json:"user_id"`
	Content  string `json:"content"`
	CreateAt string `json:"created_at"`
	User     User   `json:"user"`
}

type CommentStore struct {
	db *sql.DB
}

func (s *CommentStore) CreateComment(ctx context.Context, comment *Comment) error {
	query := `
		INSERT INTO comments (post_id, user_id, content)
		VALUES ($1, $2, $3) RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		comment.PostID,
		comment.UserID,
		comment.Content,
	).Scan(
		&comment.ID,
		&comment.CreateAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *CommentStore) GetByPostID(ctx context.Context, postID int64) ([]Comment, error) {
	query := `
		SELECT comments.id, comments.post_id, comments.user_id, comments.content, comments.created_at, 
		users.username, users.id, users.email FROM comments JOIN users on users.id = comments.user_id
		WHERE comments.post_id = $1 ORDER BY comments.created_at DESC
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, postID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []Comment{}
	for rows.Next() {
		var c Comment
		c.User = User{}
		err := rows.Scan(
			&c.ID,
			&c.PostID,
			&c.UserID,
			&c.Content,
			&c.CreateAt,
			&c.User.Username,
			&c.User.ID,
			&c.User.Email,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}
