package postgres

import (
	"context"
	"database/sql"
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

// NewCommentRepository создаёт новый PostgreSQL-репозиторий CommentRepository с комментариями.
func NewCommentRepository(pool *pgxpool.Pool) *CommentRepository {
	return &CommentRepository{pool: pool}
}

// Create сохраняет комментарий.
func (r *CommentRepository) Create(ctx context.Context, comment *dom.Comment) (dom.ID, error) {
	const query = `
		INSERT INTO comments (news_id, parent_id, user_name, content, created_at, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var id int64
	err := r.pool.QueryRow(
		ctx, query, comment.NewsID().Value(), comment.ParentID().Value(), comment.Username().Value(),
		comment.Content().Value(), comment.CreatedAt().Time().UTC().Unix(), comment.Status().Value(),
	).Scan(&id)
	if err != nil {
		return dom.ID{}, fmt.Errorf("NewCommentRepository.Create: %w", err)
	}

	commentID, err := dom.NewID(id)
	if err != nil {
		return dom.ID{}, fmt.Errorf("NewCommentRepository.Create: %w", err)
	}

	return commentID, nil
}

// UpdateStatus публикует/отклоняет комментарий.
func (r *CommentRepository) UpdateStatus(
	ctx context.Context, id dom.ID, status dom.Status, pubTime *dom.CommentTime,
) error {
	const queryWithTime = `UPDATE comments SET status = $2, pub_time = $3 WHERE id = $1`
	const queryNoTime = `UPDATE comments SET status = $2 WHERE id = $1`

	if pubTime != nil {
		_, err := r.pool.Exec(ctx, queryWithTime, id.Value(), status.Value(), pubTime.Time().UTC().Unix())
		if err != nil {
			return fmt.Errorf("CommentRepository.UpdateStatus: %w", err)
		}

		return nil
	}

	_, err := r.pool.Exec(ctx, queryNoTime, id.Value(), status.Value())
	if err != nil {
		return fmt.Errorf("CommentRepository.UpdateStatus: %w", err)
	}
	return nil
}

// FindByID находит комментарий по его ID.
func (r *CommentRepository) FindByID(ctx context.Context, id dom.ID) (*dom.Comment, error) {
	var row mapper.CommentRow
	var pubTime sql.NullInt64

	const query = `
		SELECT id, news_id, parent_id, user_name, content, pub_time, status
		FROM comments WHERE id=$1 LIMIT 1`

	err := r.pool.QueryRow(ctx, query, id.Value()).Scan(
		&row.ID, &row.NewsID, &row.ParentID, &row.Username, &row.Content, &pubTime, &row.Status,
	)
	if err != nil {
		return nil, fmt.Errorf("CommentRepository.FindByID: %w", err)
	}

	if pubTime.Valid {
		row.PubTime = pubTime.Int64
	} else {
		row.PubTime = 0
	}

	return mapper.MapRowToComment(row)
}

// FindAllByNewsID получает все комментарии конкретной новости.
func (r *CommentRepository) FindAllByNewsID(ctx context.Context, newsID dom.NewsID) ([]*dom.Comment, error) {
	const query = `
		SELECT id, news_id, parent_id, user_name, content, pub_time, status
		FROM comments 
		WHERE news_id=$1 and status=$2
		ORDER BY pub_time DESC`
	rows, err := r.pool.Query(ctx, query, newsID.Value(), dom.Approved)
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

		if err = rows.Scan(
			&row.ID, &row.NewsID, &row.ParentID, &row.Username, &row.Content, &row.PubTime, &row.Status,
		); err != nil {
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
