package posts

import (
	"comu/internal/modules/post/domain"
	"context"
)

type ListPostsUC struct {
	repo domain.PostRepository
}

type ReadPostUC struct {
	repo domain.PostRepository
}

func NewListPostsUseCase(repository domain.PostRepository) *ListPostsUC {
	return &ListPostsUC{
		repo: repository,
	}
}

func NewReadPostUseCase(repository domain.PostRepository) *ReadPostUC {
	return &ReadPostUC{
		repo: repository,
	}
}

func (useCase *ListPostsUC) Execute(ctx context.Context, paginator domain.Paginator) ([]domain.Post, *domain.Cursor, error) {
	if paginator.Limit <= 0 {
		paginator.Limit = domain.DefaultPaginatorLimit
	}

	post, cursor, err := useCase.repo.List(ctx, paginator)

	if err != nil {
		return []domain.Post{}, nil, err
	}

	return post, cursor, nil
}

func (useCase *ReadPostUC) Execute(ctx context.Context, slug string) (*domain.Post, error) {
	post, err := useCase.repo.FindBySlug(ctx, slug)

	if err != nil {
		return nil, err
	}

	return post, nil
}
