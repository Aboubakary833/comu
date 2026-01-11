package comments

import (
	"comu/internal/modules/post/domain"
	"context"

	"github.com/google/uuid"
)

type listCommentsUC struct {
	repo domain.CommentRepository
}

type deleteCommentUC struct {
	repo domain.CommentRepository
}

func NewListCommentsUseCase(repository domain.CommentRepository) *listCommentsUC {
	return &listCommentsUC{
		repo: repository,
	}
}

func NewDeleteCommentUseCase(repository domain.CommentRepository) *deleteCommentUC {
	return &deleteCommentUC{
		repo: repository,
	}
}


func (useCase *listCommentsUC) Execute(ctx context.Context, postID uuid.UUID, paginator domain.Paginator) ([]domain.Comment, *domain.Cursor, error) {
	if paginator.Limit <= 0 {
		paginator.Limit = domain.DefaultPaginatorLimit
	}

	comments, cursor, err := useCase.repo.List(ctx, postID, paginator)

	if err != nil {
		return []domain.Comment{}, nil, err
	}

	return comments, cursor, nil
}

func (useCase *deleteCommentUC) Execute(ctx context.Context, commentID, authorID uuid.UUID) error {
	comment, err := useCase.repo.Find(ctx, commentID)

	if err != nil {
		return err
	}

	if comment.UserID != authorID {
		return domain.ErrUnauthorized
	}

	return useCase.repo.Delete(ctx, comment)
}
