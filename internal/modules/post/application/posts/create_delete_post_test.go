package posts

import (
	"comu/internal/modules/post/domain"
	"comu/internal/modules/post/infra/memory"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreatePostUseCase(t *testing.T) {
	repo := memory.NewInMemoryPostsRepository(nil)
	useCase := NewCreatePostUseCase(repo)

	input := CreatePostInput{
		UserID: uuid.New(),
		Title: "Test post",
		Content: "This is a test post",
	}

	post, err := useCase.Execute(context.Background(), input)
	_assert := assert.New(t)

	if _assert.NoError(err) {
		_assert.Equal(input.UserID, post.UserID)
		_assert.Equal(input.Title, post.Title)
		_assert.Equal(input.Content, post.Content)
	}
}

func TestDeletePostUseCase(t *testing.T) {
	repo := memory.NewInMemoryPostsRepository(nil)
	ctx := context.Background()

	post := domain.NewPost(uuid.New(), "Test post", "Weird test post content")
	repo.Store(ctx, post)

	useCase := NewDeletePostUseCase(repo)
	err := useCase.Execute(ctx, post.ID)

	if assert.NoError(t, err) {
		_, err = repo.FindByID(ctx, post.ID)
		assert.ErrorIs(t, err, domain.ErrPostNotFound)
	}
}
