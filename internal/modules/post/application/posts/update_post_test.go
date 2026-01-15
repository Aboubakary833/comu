package posts

import (
	"comu/internal/modules/post/domain"
	"comu/internal/modules/post/infra/memory"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpdatePostUseCase(t *testing.T) {

	t.Run("it should update both title and content", func(t *testing.T) {
		repo := memory.NewInMemoryPostsRepository(nil)
		ctx := context.Background()
		_assert := assert.New(t)
		userID := uuid.New()

		post := domain.NewPost(userID, "Test post title", "This is test post title")
		repo.Store(ctx, post)

		useCase := NewUpdatePostUseCase(repo)

		slug, err := useCase.Execute(ctx, UpdatePostInput{
			PostID:   post.ID,
			AuthorID: userID,
			Title:    "Updated test post title",
			Content:  "Test post updated content",
		})

		if _assert.NoError(err) {
			retrievedPost, _ := repo.FindBySlug(ctx, slug)
			_assert.Equal(post.ID, retrievedPost.ID)
			_assert.NotEqual(post.Title, retrievedPost.Title)
			_assert.NotEqual(post.Content, retrievedPost.Content)
			_assert.Equal(post.CreatedAt, retrievedPost.CreatedAt)
			_assert.NotEqual(post.UpdatedAt, retrievedPost.UpdatedAt)
		}
	})

	t.Run("it should update only the content", func(t *testing.T) {
		repo := memory.NewInMemoryPostsRepository(nil)
		ctx := context.Background()
		_assert := assert.New(t)
		userID := uuid.New()

		post := domain.NewPost(userID, "Test post title", "This is test post title")
		repo.Store(ctx, post)

		useCase := NewUpdatePostUseCase(repo)

		slug, err := useCase.Execute(ctx, UpdatePostInput{
			PostID:   post.ID,
			AuthorID: userID,
			Title:    post.Title,
			Content:  "Test post updated content",
		})

		if _assert.NoError(err) {
			retrievedPost, _ := repo.FindBySlug(ctx, slug)
			_assert.Equal(post.ID, retrievedPost.ID)
			_assert.Equal(post.Title, retrievedPost.Title)
			_assert.NotEqual(post.Content, retrievedPost.Content)
			_assert.Equal(post.CreatedAt, retrievedPost.CreatedAt)
			_assert.NotEqual(post.UpdatedAt, retrievedPost.UpdatedAt)
		}
	})

	t.Run("it should fail and return ErrUnauthorized", func(t *testing.T) {
		repo := memory.NewInMemoryPostsRepository(nil)
		ctx := context.Background()

		post := domain.NewPost(uuid.New(), "Test post title", "This is test post title")
		repo.Store(ctx, post)

		useCase := NewUpdatePostUseCase(repo)

		_, err := useCase.Execute(ctx, UpdatePostInput{
			PostID:   post.ID,
			AuthorID: uuid.New(),
			Title:    post.Title,
			Content:  "Test post updated content",
		})

		assert.ErrorIs(t, err, domain.ErrUnauthorized)
	})
}
