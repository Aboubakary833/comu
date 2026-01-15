package posts

import (
	"comu/internal/modules/post/domain"
	"comu/internal/modules/post/infra/memory"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestListPostsUseCase(t *testing.T) {

	t.Run("it should return 10 posts when negative Limit value is passed to paginator", func(t *testing.T) {
		repo := memory.NewInMemoryPostsRepository(nil)
		_assert := assert.New(t)

		repo.FillWithRandomPosts(uuid.Nil, 15)
		useCase := NewListPostsUseCase(repo)

		posts, cursor, err := useCase.Execute(context.Background(), domain.Paginator{Limit: -3})

		if _assert.NoError(err) {
			_assert.Equal(10, len(posts))
			_assert.NotNil(cursor)
		}
	})

	t.Run("it should return 10 posts when paginator is empty(or has zero as Limit value)", func(t *testing.T) {
		repo := memory.NewInMemoryPostsRepository(nil)
		_assert := assert.New(t)

		repo.FillWithRandomPosts(uuid.Nil, 15)
		useCase := NewListPostsUseCase(repo)

		posts, cursor, err := useCase.Execute(context.Background(), domain.Paginator{})

		if _assert.NoError(err) {
			_assert.Equal(10, len(posts))
			_assert.NotNil(cursor)
		}
	})

	t.Run("it should return the correct length of posts after a given post", func(t *testing.T) {
		repo := memory.NewInMemoryPostsRepository(nil)
		ctx := context.Background()
		_assert := assert.New(t)

		repo.FillWithRandomPosts(uuid.Nil, 2)
		post := domain.NewPost(uuid.New(), "Test post", "Test post content")
		repo.Store(ctx, post)
		repo.FillWithRandomPosts(uuid.Nil, 12)

		useCase := NewListPostsUseCase(repo)

		posts, nextCursor, err := useCase.Execute(ctx, domain.Paginator{
			Limit: 8,
			After: &domain.Cursor{
				ID:        post.ID,
				CreatedAt: post.CreatedAt,
			},
		})

		if _assert.NoError(err) {
			_assert.Equal(8, len(posts))
			_assert.NotNil(nextCursor)

			first := posts[0]
			_assert.True(first.CreatedAt.After(post.CreatedAt))
		}
	})
}

func TestReadPostUseCase(t *testing.T) {
	t.Run("it should fail and return ErrPostNotFound", func(t *testing.T) {
		repo := memory.NewInMemoryPostsRepository(nil)
		useCase := NewReadPostUseCase(repo)

		slug := domain.MakePostSlug("Test post title")

		_, err := useCase.Execute(context.Background(), slug)
		assert.Equal(t, err, domain.ErrPostNotFound)
	})

	t.Run("it should succeed and return a given post", func(t *testing.T) {
		repo := memory.NewInMemoryPostsRepository(nil)
		ctx := context.Background()
		_assert := assert.New(t)

		post := domain.NewPost(uuid.New(), "Test post", "That is test post content")
		repo.Store(ctx, post)

		useCase := NewReadPostUseCase(repo)

		retrievedPost, err := useCase.Execute(ctx, post.Slug)

		if _assert.NoError(err) {
			_assert.Equal(post.Title, retrievedPost.Title)
			_assert.Equal(post.Slug, retrievedPost.Slug)
			_assert.Equal(post.Content, retrievedPost.Content)
		}
	})
}
