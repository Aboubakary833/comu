package posts

import (
	"comu/internal/modules/post/domain"
	"context"
)


type listPostsUC struct {
	repo domain.PostRepository
}

type readPostUC struct {
	repo domain.PostRepository
}

func NewListPostsUseCase(repository domain.PostRepository) *listPostsUC {
	return &listPostsUC{
		repo: repository,
	}
}

func NewReadPostUseCase(repository domain.PostRepository) *readPostUC {
	return &readPostUC{
		repo: repository,
	}
}

func (useCase *listPostsUC) Execute(ctx context.Context, paginator domain.Paginator) ([]domain.Post, *domain.Cursor, error) {
	if paginator.Limit <= 0 {
		paginator.Limit = domain.DefaultPaginatorLimit
	}
	
	post, cursor, err := useCase.repo.List(ctx, paginator)

	if err != nil {
		return []domain.Post{}, nil, err
	}

	return post, cursor, nil
}

func (useCase *readPostUC) Execute(ctx context.Context, slug string) (*domain.Post, error) {
	post, err := useCase.repo.FindBySlug(ctx, slug)

	if err != nil {
		return nil, err
	}

	return post, nil
}
