package posts

import (
	"comu/internal/modules/post/domain"
	"context"

	"github.com/google/uuid"
)

type CreatePostInput struct {
	UserID	uuid.UUID
	Title	string
	Content string
}

type CreatePostUC struct {
	repo domain.PostRepository
}

type DeletePostUC struct {
	repo domain.PostRepository
}

func NewCreatePostUseCase(repository domain.PostRepository) *CreatePostUC {
	return &CreatePostUC{
		repo: repository,
	}
}

func NewDeletePostUseCase(repository domain.PostRepository) *DeletePostUC {
	return &DeletePostUC{
		repo: repository,
	}
}

func (useCase *CreatePostUC) Execute(ctx context.Context, input CreatePostInput) (*domain.Post, error) {
	post := domain.NewPost(input.UserID, input.Title, input.Content)

	err := useCase.repo.Store(ctx, post)

	if err != nil {
		return nil, err
	}

	return post, nil
}

func (useCase *DeletePostUC) Execute(ctx context.Context, postID, authorID uuid.UUID) error {
	post, err := useCase.repo.FindByID(ctx, postID)

	if err != nil {
		return err
	}

	if post.UserID != authorID {
		return domain.ErrUnauthorized
	}

	return useCase.repo.Delete(ctx, post)
}
