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

type createPostUC struct {
	repo domain.PostRepository
}

type deletePostUC struct {
	repo domain.PostRepository
}

func NewCreatePostUseCase(repository domain.PostRepository) *createPostUC {
	return &createPostUC{
		repo: repository,
	}
}

func NewDeletePostUseCase(repository domain.PostRepository) *deletePostUC {
	return &deletePostUC{
		repo: repository,
	}
}

func (useCase *createPostUC) Execute(ctx context.Context, input CreatePostInput) (*domain.Post, error) {
	post := domain.NewPost(input.UserID, input.Title, input.Content)

	err := useCase.repo.Store(ctx, post)

	if err != nil {
		return nil, err
	}

	return post, nil
}

func (useCase *deletePostUC) Execute(ctx context.Context, postID, authorID uuid.UUID) error {
	post, err := useCase.repo.FindByID(ctx, postID)

	if err != nil {
		return err
	}

	if post.UserID != authorID {
		return domain.ErrUnauthorized
	}

	return useCase.repo.Delete(ctx, post)
}
