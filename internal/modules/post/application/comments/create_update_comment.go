package comments

import (
	"comu/internal/modules/post/domain"
	"context"

	"github.com/google/uuid"
)

type CreateCommentInput struct {
	PostID   uuid.UUID
	AuthorID uuid.UUID
	Content  string
}

type CreateCommentUC struct {
	repo domain.CommentRepository
}

type UpdateCommentUC struct {
	repo domain.CommentRepository
}

func NewCreateCommentUseCase(repository domain.CommentRepository) *CreateCommentUC {
	return &CreateCommentUC{
		repo: repository,
	}
}

func NewUpdateCommentUseCase(repository domain.CommentRepository) *UpdateCommentUC {
	return &UpdateCommentUC{
		repo: repository,
	}
}

func (useCase *CreateCommentUC) Execute(ctx context.Context, input CreateCommentInput) (*domain.Comment, error) {
	comment := domain.NewComment(input.PostID, input.AuthorID, input.Content)

	err := useCase.repo.Store(ctx, comment)

	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (useCase *UpdateCommentUC) Execute(ctx context.Context, commentID, authorID uuid.UUID, content string) error {
	comment, err := useCase.repo.Find(ctx, commentID)

	if err != nil {
		return err
	}

	if comment.UserID != authorID {
		return domain.ErrUnauthorized
	}

	comment.Content = content

	return useCase.repo.Update(ctx, comment)
}
