package posts

import (
	"comu/internal/modules/post/domain"
	"context"

	"github.com/google/uuid"
)

type UpdatePostInput struct {
	PostID uuid.UUID
	Title 	string
	Content string
}

type updatePostUC struct {
	repo domain.PostRepository
}

func NewUpdatePostUseCase(repository domain.PostRepository) *updatePostUC {
	return &updatePostUC{
		repo: repository,
	}
}

func (useCase *updatePostUC) Execute(ctx context.Context, input UpdatePostInput) (slug string, err error) {
	post, err := useCase.repo.FindByID(ctx, input.PostID)

	if err != nil {
		return
	}

	if input.Title != post.Title {
		slug := domain.MakePostSlug(input.Title)
		post.Title = input.Title
		post.Slug = slug
	}
	post.Content = input.Content
	err = useCase.repo.Update(ctx, post)

	if err != nil {
		return
	}
	
	return post.Slug, nil
}
