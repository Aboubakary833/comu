package comments

import (
	"comu/internal/modules/post/domain"
	"comu/internal/modules/post/infra/memory"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestListCommentUseCase(t *testing.T) {

	t.Run("it should return 10 comments when negative Limit value is passed to paginator", func(t *testing.T) {
		repo := memory.NewInMemoryCommentsRepository(nil)
		_assert := assert.New(t)
		postID := uuid.New()

		repo.FillWithRandomComments(postID, uuid.Nil, 15)
		repo.FillWithRandomComments(uuid.Nil, uuid.Nil, 10)

		useCase := NewListCommentsUseCase(repo)

		comments, cursor, err := useCase.Execute(context.Background(), postID, domain.Paginator{Limit: -3})

		if _assert.NoError(err) {
			_assert.Equal(10, len(comments))
			_assert.NotNil(cursor)
		}
	})

	t.Run("it should return 10 comments when paginator is empty(or has zero as Limit value)", func(t *testing.T) {
		repo := memory.NewInMemoryCommentsRepository(nil)
		_assert := assert.New(t)
		postID := uuid.New()

		repo.FillWithRandomComments(postID, uuid.Nil, 15)
		useCase := NewListCommentsUseCase(repo)

		comments, cursor, err := useCase.Execute(context.Background(), postID, domain.Paginator{})

		if _assert.NoError(err) {
			_assert.Equal(10, len(comments))
			_assert.NotNil(cursor)
		}
	})

	t.Run("it should return the correct length of comments after a given comments", func(t *testing.T) {
		repo := memory.NewInMemoryCommentsRepository(nil)
		ctx := context.Background()
		_assert := assert.New(t)
		postID := uuid.New()

		repo.FillWithRandomComments(postID, uuid.Nil, 2)
		comment := domain.NewComment(postID, uuid.New(), "Test post content")
		repo.Store(ctx, comment)
		repo.FillWithRandomComments(postID, uuid.Nil, 12)

		useCase := NewListCommentsUseCase(repo)

		comments, nextCursor, err := useCase.Execute(ctx, postID, domain.Paginator{
			Limit: 8,
			After: &domain.Cursor{
				ID:        comment.ID,
				CreatedAt: comment.CreatedAt,
			},
		})

		if _assert.NoError(err) {
			_assert.Equal(8, len(comments))
			_assert.NotNil(nextCursor)

			first := comments[0]
			_assert.True(first.CreatedAt.After(comment.CreatedAt))
		}
	})

}

func TestDeleteCommentUseCase(t *testing.T) {

	t.Run("it should fail and return ErrCommentNotFound", func(t *testing.T) {
		repo := memory.NewInMemoryCommentsRepository(nil)
		useCase := NewDeleteCommentUseCase(repo)

		err := useCase.Execute(context.Background(), uuid.New(), uuid.New())
		assert.ErrorIs(t, err, domain.ErrCommentNotFound)
	})

	t.Run("it should fail and return ErrUnauthorized", func(t *testing.T) {
		repo := memory.NewInMemoryCommentsRepository(nil)
		ctx := context.Background()

		comment := domain.NewComment(uuid.New(), uuid.New(), "Test comment")
		repo.Store(ctx, comment)

		useCase := NewDeleteCommentUseCase(repo)

		err := useCase.Execute(ctx, comment.ID, uuid.New())
		assert.ErrorIs(t, err, domain.ErrUnauthorized)
	})

	t.Run("it should succeed and delete the given comment", func(t *testing.T) {
		repo := memory.NewInMemoryCommentsRepository(nil)
		ctx := context.Background()

		comment := domain.NewComment(uuid.New(), uuid.New(), "Test comment")
		repo.Store(ctx, comment)

		useCase := NewDeleteCommentUseCase(repo)

		err := useCase.Execute(ctx, comment.ID, comment.UserID)

		if assert.NoError(t, err) {
			_, err := repo.Find(ctx, comment.ID)
			assert.ErrorIs(t, err, domain.ErrCommentNotFound)
		}
	})
}
