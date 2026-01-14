package comments

import (
	"comu/internal/modules/post/domain"
	"context"

	"github.com/google/uuid"
)

type ListCommentsUC struct {
	repo domain.CommentRepository
}

type DeleteCommentUC struct {
	repo domain.CommentRepository
}

func NewListCommentsUseCase(repository domain.CommentRepository) *ListCommentsUC {
	return &ListCommentsUC{
		repo: repository,
	}
}

func NewDeleteCommentUseCase(repository domain.CommentRepository) *DeleteCommentUC {
	return &DeleteCommentUC{
		repo: repository,
	}
}


func (useCase *ListCommentsUC) Execute(ctx context.Context, postID uuid.UUID, paginator domain.Paginator) ([]domain.Comment, *domain.Cursor, error) {
	if paginator.Limit <= 0 {
		paginator.Limit = domain.DefaultPaginatorLimit
	}

	comments, cursor, err := useCase.repo.List(ctx, postID, paginator)

	if err != nil {
		return []domain.Comment{}, nil, err
	}

	return comments, cursor, nil
}

func (useCase *DeleteCommentUC) Execute(ctx context.Context, commentID, authorID uuid.UUID) error {
	comment, err := useCase.repo.Find(ctx, commentID)

	if err != nil {
		return err
	}

	if comment.UserID != authorID {
		return domain.ErrUnauthorized
	}

	return useCase.repo.Delete(ctx, comment)
}
