package postgres

import (
	"context"
	"fmt"
	dom "github.com/ee-crocush/go-news/go-comments/internal/domain/comment"
	"github.com/ee-crocush/go-news/go-comments/internal/infrastructure/repo/postgres/mapper"
	"github.com/jackc/pgx/v4/pgxpool"
)

var _ dom.Repository = (*CommentRepository)(nil)

// CommentRepository представляет собой репозиторий для работы с комментариями.
type CommentRepository struct {
	pool *pgxpool.Pool
}

// NewCommentRepository создаёт новый PostgreSQL-репозиторий с комментариями.
func NewCommentRepository(pool *pgxpool.Pool) *CommentRepository {
	return &CommentRepository{pool: pool}
}

// Create сохраняет комментарий.
func (r *CommentRepository) Create(ctx context.Context, comment *dom.Comment) error {
	const query = `
		INSERT INTO comments (parent_id, user_name, content, pub_time)
		VALUES ($1, $2, $3, $4)
	`

	var parentID *int64

	_, err := r.pool.Exec(
		ctx, query, parentID, comment.Username().Value(), comment.Content().Value(),
		comment.PubTime().Time().UTC().Unix(),
	)
	if err != nil {
		return fmt.Errorf("NewCommentRepository.Create: %w", err)
	}

	return nil
}

// FindByID находит комментарий по его ID.
func (r *CommentRepository) FindByID(ctx context.Context, id dom.ID) (*dom.Comment, error) {
	var row mapper.CommentRow

	const query = `SELECT id, parent_id, user_name, content, pub_time FROM comments WHERE id=$1 LIMIT 1`

	err := r.pool.QueryRow(ctx, query, id.Value()).Scan(
		&row.ID, &row.ParentID, &row.Username, &row.Content, &row.PubTime,
	)
	if err != nil {
		return nil, fmt.Errorf("CommentRepository.FindByID: %w", err)
	}

	return mapper.MapRowToComment(row)
}

// FindAll получает все комментарии.
func (r *CommentRepository) FindAll(ctx context.Context) ([]*dom.Comment, error) {
	const query = `SELECT id, parent_id, user_name, content, pub_time FROM comments`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("CommentRepository.FindAll: %w", err)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("CommentRepository.FindAll: %w", err)
	}
	defer rows.Close()

	var comments []*dom.Comment

	for rows.Next() {
		var row mapper.CommentRow

		if err = rows.Scan(&row.ID, &row.ParentID, &row.Username, &row.Content, &row.PubTime); err != nil {
			return nil, fmt.Errorf("CommentRepository.FindAll: %w", err)
		}

		comment, err := mapper.MapRowToComment(row)
		if err != nil {
			return nil, fmt.Errorf("CommentRepository.FindAll: %w", err)
		}

		comments = append(comments, comment)
	}

	return comments, nil
}
