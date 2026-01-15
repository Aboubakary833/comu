package comments

import (
	"comu/internal/modules/post/domain"
	"comu/internal/modules/post/infra/memory"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateCommentUseCase(t *testing.T) {
	repo := memory.NewInMemoryCommentsRepository(nil)
	useCase := NewCreateCommentUseCase(repo)

	input := CreateCommentInput{
		PostID:   uuid.New(),
		AuthorID: uuid.New(),
		Content:  "Test comment",
	}

	_, err := useCase.Execute(context.Background(), input)
	assert.NoError(t, err)
}

func TestUpdateCommentUseCase(t *testing.T) {

	t.Run("it should fail and return ErrCommentNotFound", func(t *testing.T) {
		repo := memory.NewInMemoryCommentsRepository(nil)
		useCase := NewUpdateCommentUseCase(repo)

		err := useCase.Execute(context.Background(), uuid.New(), uuid.New(), "Test comment text")
		assert.ErrorIs(t, err, domain.ErrCommentNotFound)
	})

	t.Run("it should fail and return ErrUnauthorized", func(t *testing.T) {
		repo := memory.NewInMemoryCommentsRepository(nil)
		ctx := context.Background()

		comment := domain.NewComment(uuid.New(), uuid.New(), "Comment content")
		repo.Store(ctx, comment)

		useCase := NewUpdateCommentUseCase(repo)
		err := useCase.Execute(ctx, comment.ID, uuid.New(), "Updated comment content")
		assert.ErrorIs(t, err, domain.ErrUnauthorized)
	})

	t.Run("it should successfully update the comment", func(t *testing.T) {
		repo := memory.NewInMemoryCommentsRepository(nil)
		ctx := context.Background()
		_assert := assert.New(t)

		comment := domain.NewComment(uuid.New(), uuid.New(), "Comment content")
		repo.Store(ctx, comment)

		useCase := NewUpdateCommentUseCase(repo)
		err := useCase.Execute(ctx, comment.ID, comment.UserID, "Updated comment content")

		if _assert.NoError(err) {
			retrievedComment, _ := repo.Find(ctx, comment.ID)
			_assert.NotEqual(comment.Content, retrievedComment.Content)
			_assert.NotEqual(comment.UpdatedAt, retrievedComment.UpdatedAt)
			_assert.Equal(comment.CreatedAt, retrievedComment.CreatedAt)
		}
	})
}
