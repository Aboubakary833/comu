package mysql

import (
	"comu/internal/modules/post/domain"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type commentsRepository struct {
	db *sql.DB
}

func NewCommentsRepository(db *sql.DB) *commentsRepository {
	return &commentsRepository{
		db: db,
	}
}

func (repo *commentsRepository) Find(ctx context.Context, ID uuid.UUID) (*domain.Comment, error) {
	query := "SELECT * FROM comments WHERE id = UUID_TO_BIN(?);"
	comment := &domain.Comment{}

	err := repo.db.QueryRowContext(ctx, query, ID.String()).Scan(
		&comment.ID, &comment.PostID, &comment.UserID,
		&comment.Content, &comment.CreatedAt, &comment.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrCommentNotFound
		}

		return nil, err
	}

	return comment, nil
}

func (repo *commentsRepository) ListAll(ctx context.Context, postID uuid.UUID) ([]domain.Comment, error) {
	query := "SELECT * FROM comments WHERE post_id = UUID_TO_BIN(?);"
	rows, err := repo.db.QueryContext(ctx, query, postID.String())

	if err != nil {
		return []domain.Comment{}, err
	}

	return repo.getCommentsFromRows(rows)
}

func (repo *commentsRepository) List(ctx context.Context, postID uuid.UUID, paginator domain.Paginator) ([]domain.Comment, error) {

	if paginator.After == nil {
		query := `
			SELECT *
			FROM comments
			WHERE post_id = UUID_TO_BIN(?)
			ORDER BY created_at DESC, id DESC
			LIMIT ?;
		`
		rows, err := repo.db.QueryContext(ctx, query, postID.String(), paginator.Limit)

		if err != nil {
			return []domain.Comment{}, err
		}

		return repo.getCommentsFromRows(rows)
	}

	query := `
		SELECT *
		FROM comments
		WHERE
			post_id = UUID_TO_BIN(?)
			(created_at < ?)
			OR (created_at = ? AND id < ?)
		ORDER BY created_at DESC, id DESC
		LIMIT ?;
	`
	createdAt := paginator.After.CreatedAt

	rows, err := repo.db.QueryContext(
		ctx, query, postID.String(), createdAt, createdAt,
		paginator.After.ID, paginator.Limit,
	)

	if err != nil {
		return []domain.Comment{}, err
	}

	return repo.getCommentsFromRows(rows)
}

func (repo *commentsRepository) Store(ctx context.Context, comment *domain.Comment) error {
	query := `
		INSERT INTO comments (
			id, post_id, user_id, content, created_at, updated_at
		) VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), UUID_TO_BIN(?), ?, ?, ?);
	`

	_, err := repo.db.ExecContext(
		ctx, query, comment.ID.String(), comment.PostID.String(), comment.UserID.String(),
		comment.Content, comment.CreatedAt, comment.UpdatedAt,
	)

	return err
}

func (repo *commentsRepository) Update(ctx context.Context, comment *domain.Comment) error {
	query := "UPDATE comments SET content = ?, updated_at = ? WHERE id = UUID_TO_BIN(?);"
	comment.UpdatedAt = time.Now()

	_, err := repo.db.ExecContext(ctx, query, comment.Content, comment.UpdatedAt, comment.ID.String())

	return err
}

func (repo *commentsRepository) Delete(ctx context.Context, comment *domain.Comment) error {
	query := "DELETE FROM comments WHERE id = UUID_TO_BIN(?);"
	_, err := repo.db.ExecContext(ctx, query, comment.ID.String())

	return err
}

func (repo *commentsRepository) getCommentsFromRows(rows *sql.Rows) ([]domain.Comment, error) {
	comments := []domain.Comment{}

	for rows.Next() {
		comment := domain.Comment{}

		err := rows.Scan(
			&comment.ID, &comment.PostID, &comment.UserID,
			&comment.Content, &comment.CreatedAt, &comment.UpdatedAt,
		)

		if err != nil {
			return []domain.Comment{}, err
		}

		comments = append(comments, comment)
	}

	return comments, nil
}
