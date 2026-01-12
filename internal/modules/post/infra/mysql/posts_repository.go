package mysql

import (
	"comu/internal/modules/post/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type postsRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *postsRepository {
	return &postsRepository{
		db: db,
	}
}

func (repo *postsRepository) FindByID(ctx context.Context, ID uuid.UUID) (*domain.Post, error) {
	return repo.findQuery(ctx, "id", ID.String())
}

func (repo *postsRepository) FindBySlug(ctx context.Context, slug string) (*domain.Post, error) {
	return repo.findQuery(ctx, "slug", slug)
}

func (repo *postsRepository) ListAll(ctx context.Context) ([]domain.Post, error) {
	query := "SELECT * FROM posts ORDER BY created_at DESC;"
	rows, err := repo.db.QueryContext(ctx, query)

	if err != nil {
		return []domain.Post{}, err
	}

	return repo.getPostFromRows(rows)
}

func (repo *postsRepository) List(ctx context.Context, paginator domain.Paginator) ([]domain.Post, error) {

	if paginator.After == nil {
		query := `
			SELECT *
			FROM posts
			ORDER BY created_at DESC, id DESC
			LIMIT ?;
		`
		rows, err := repo.db.QueryContext(ctx, query, paginator.Limit)

		if err != nil {
			return []domain.Post{}, err
		}

		return repo.getPostFromRows(rows)
	}

	query := `
			SELECT *
			FROM posts
			WHERE
				(created_at < ?)
				OR (created_at = ? AND id < ?)
			ORDER BY created_at DESC, id DESC
			LIMIT ?;
		`
	createdAt := paginator.After.CreatedAt
	rows, err := repo.db.QueryContext(
		ctx, query, createdAt, createdAt,
		paginator.After.ID.String(), paginator.Limit,
	)

	if err != nil {
		return []domain.Post{}, err
	}

	return repo.getPostFromRows(rows)
}

func (repo *postsRepository) Store(ctx context.Context, post *domain.Post) error {
	query := `
		INSERT INTO posts (
			id, user_id, title, slug,
			content, created_at, updated_at
		) VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), ?, ?, ?, ?, ?);
	`

	_, err := repo.db.ExecContext(
		ctx, query, post.ID.String(), post.UserID.String(), post.Title,
		post.Slug, post.Content, post.CreatedAt, post.UpdatedAt,
	)

	return err
}

func (repo *postsRepository) Update(ctx context.Context, post *domain.Post) error {
	query := `
		UPDATE posts SET title = ?, slug = ?, content = ?, updated_at = ?
		WHERE id = UUID_TO_BIN(?);
	`
	post.UpdatedAt = time.Now()

	_, err := repo.db.ExecContext(
		ctx, query, post.Title, post.Slug,
		post.Content, post.UpdatedAt, post.ID.String(),
	)

	return err
}

func (repo *postsRepository) Delete(ctx context.Context, post *domain.Post) error {
	query := "DELETE FROM posts WHERE id = UUID_TO_BIN(?);"
	_, err := repo.db.ExecContext(ctx, query, post.ID.String())

	return err
}


func (repo *postsRepository) findQuery(ctx context.Context, column, value string) (*domain.Post, error) {
	queryVal := "?"

	if column == "id" {
		queryVal = "UUID_TO_BIN(?)"
	}
	query := fmt.Sprintf("SELECT * FROM posts WHERE %s = %s;", column, queryVal)
	post := &domain.Post{}

	err := repo.db.QueryRowContext(ctx, query, value).Scan(
		&post.ID, &post.UserID, &post.Title, &post.Slug,
		&post.Content, &post.CreatedAt, &post.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrPostNotFound
		}

		return nil, err
	}

	return post, nil
}

func (repo *postsRepository) getPostFromRows(rows *sql.Rows) ([]domain.Post, error) {
	posts := []domain.Post{}

	for rows.Next() {
		post := domain.Post{}
		err := rows.Scan(
			&post.ID, &post.UserID, &post.Title, &post.Slug,
			&post.Content, &post.CreatedAt, &post.UpdatedAt,
		)

		if err != nil {
			return []domain.Post{}, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}
