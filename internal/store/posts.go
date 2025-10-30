package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

// CommentList implements sql.Scanner for JSONB array data
type CommentList []Comment

type Post struct {
	ID        int64       `json:"id"` //Json unmarshal
	Content   string      `json:"content"`
	Title     string      `json:"title"`
	UserID    int64       `json:"user_id"`
	Tags      []string    `json:"tags"`
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at"`
	Version   string      `json:"version"`
	Comments  CommentList `json:"comments"`
	User      User        `json:"user"`
}

type PostWithMetadata struct {
	Post
	CommentCount int `json:"comments_count"`
}

type PaginatedFeedQuery struct {
	Limit  int      `json:"limit" validate:"gte=1,lte=20"`
	Offset int      `json:"offset" validate:"gte=0"`
	Sort   string   `json:"sort" validate:"oneof=asc desc"`
	Tags   []string `json:"tags" validate:"max=5"`
	Search string   `json:"search" validate:"max=100"`
	Since  string   `json:"since"`
	Until  string   `json:"until"`
}

type PostStore struct {
	db *sql.DB
}

// check out packages that can help with sql. E.g, go-gorm, sqlx
func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts (content, title, user_id, tags)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	if err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	); err != nil {
		return err
	}

	return nil
}

func (s *PostStore) GetByID(ctx context.Context, postID int64) (*Post, error) {
	query := `
		SELECT id, user_id, title, content, created_at, updated_at, tags, version FROM posts WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	// come back and check the comment slices
	var post Post
	err := s.db.QueryRowContext(ctx, query, postID).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		pq.Array(&post.Tags),
		&post.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrorNotFound
		default:
			return nil, err
		}
	}

	return &post, nil
}

func (s *PostStore) GetUserFeed(ctx context.Context, userID int64, feedQuery *PaginatedFeedQuery) ([]PostWithMetadata, error) {
	// query := `
	// 	SELECT p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags, u.username,
	// 	COUNT(c.id) AS comments_count FROM posts p LEFT JOIN comments c ON c.post_id = p.id
	// 	LEFT JOIN users u ON p.user_id = u.id JOIN followers f ON f.follower_id = p.user_id OR
	// 	p.user_id = $1 WHERE f.user_id = $1 OR p.user_id = $1 GROUP BY p.id, u.username
	// 	ORDER BY p.created_at DESC
	// OR (p.title ILIKE '%' || $4 || '%' OR p.content ILIKE '%' || $4 || '%')
	//OR (p.tags @> $5 OR $5 = '{}')
	// `
	query := `
        SELECT
            p.id, p.user_id, p.title, p.content, p.created_at, p.version,
            p.tags,
            u.username,
            COALESCE(
                c_agg.comments_json,
                '[]'::jsonb
            ) AS comments, -- Aggregated comments array
            (SELECT COUNT(c.id) FROM comments c WHERE c.post_id = p.id) AS comments_count -- Simpler way to get count
        FROM
            posts p
        LEFT JOIN
            users u ON p.user_id = u.id
        -- Aggregation of comments into a JSON array for each post
        LEFT JOIN LATERAL (
            SELECT
                JSONB_AGG(
                    JSONB_BUILD_OBJECT(
                        'id', c.id,
                        'content', c.content,
                        'created_at', c.created_at,
                        'user', JSONB_BUILD_OBJECT(
                            'id', cu.id,
                            'username', cu.username
                        )
                    )
                ) AS comments_json
            FROM
                comments c
            LEFT JOIN
                users cu ON c.user_id = cu.id
            WHERE
                c.post_id = p.id
        ) c_agg ON TRUE
        -- Feed Logic
        JOIN followers f ON f.follower_id = p.user_id OR p.user_id = $1
        WHERE f.user_id = $1 OR p.user_id = $1 -- come back and filter
        GROUP BY p.id, u.username, c_agg.comments_json
        ORDER BY p.created_at ` + feedQuery.Sort + `
		LIMIT $2 OFFSET $3
    `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(
		ctx, query, userID, feedQuery.Limit, feedQuery.Offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var feed []PostWithMetadata

	for rows.Next() {
		var p PostWithMetadata
		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Title,
			&p.Content,
			&p.CreatedAt,
			&p.Version,
			pq.Array(&p.Tags),
			&p.User.Username,
			&p.Comments,
			&p.CommentCount,
		)
		if err != nil {
			return nil, err
		}

		feed = append(feed, p)
	}
	return feed, nil
}

func (s *PostStore) Delete(ctx context.Context, postID int64) error {
	query := `DELETE FROM posts WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, postID)

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

func (s *PostStore) Update(ctx context.Context, post *Post) error {
	query := `
		UPDATE posts SET title = $1, content = $2, version = version + 1 
		WHERE id = $3 AND version = $4 RETURNING version
		`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Content,
		post.ID,
		post.Version,
	).Scan(&post.Version)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrorConflict
		default:
			return err
		}
	}

	return nil
}

func (cl *CommentList) Scan(src interface{}) error {
	var source []byte
	switch v := src.(type) {
	case string:
		source = []byte(v)
	case []byte:
		source = v
	case nil:
		*cl = CommentList{}
		return nil
	default:
		return fmt.Errorf("unsupported type for CommentList scan: %T", src)
	}
	return json.Unmarshal(source, cl)
}
