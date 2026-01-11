package comments

import (
	"comu/internal/modules/post/domain"
	"context"

	"github.com/google/uuid"
)

type CreateCommentInput struct {
	PostID 	 uuid.UUID
	AuthorID uuid.UUID
	Content  string
}

type createCommentUC struct {
	repo domain.CommentRepository
}

type updateCommentUC struct {
	repo domain.CommentRepository
}

func NewCreateCommentUseCase(repository domain.CommentRepository) *createCommentUC {
	return &createCommentUC{
		repo: repository,
	}
}

func NewUpdateCommentUseCase(repository domain.CommentRepository) *updateCommentUC {
	return &updateCommentUC{
		repo: repository,
	}
}

func (useCase *createCommentUC) Execute(ctx context.Context, input CreateCommentInput) (*domain.Comment, error) {
	comment := domain.NewComment(input.PostID, input.AuthorID, input.Content)

	err := useCase.repo.Store(ctx, comment)

	if err != nil {
		return nil, err
	}

	return comment, nil
}


func (useCase *updateCommentUC) Execute(ctx context.Context, commentID, authorID uuid.UUID, content string) error {
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
