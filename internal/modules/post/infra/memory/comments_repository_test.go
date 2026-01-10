package memory

import (
	"comu/internal/modules/post/domain"
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)


func TestInMemoryCommentsRepositoryStoreMethod(t *testing.T) {
	repo := NewInMemoryCommentsRepository(nil)
	ctx := context.Background()

	userID := uuid.New()
	postID := uuid.New()

	for i := range 5 {
		content := fmt.Sprintf("Post #%s comment #%d content", postID.String(), i + 1)
		comment := domain.NewComment(postID, userID, content)

		err := repo.Store(ctx, comment)
		assert.NoError(t, err)
	}
}

func TestInMemoryCommentsRepositoryFindMethod(t *testing.T) {

	t.Run("it should fail and return ErrCommentNotFound", func(t *testing.T) {
		repo := NewInMemoryCommentsRepository(nil)
		commentID := uuid.New()

		_, err := repo.Find(context.Background(), commentID)
		assert.ErrorIs(t, err, domain.ErrCommentNotFound)
	})

	t.Run("it should succeed and return the retrieved comment", func(t *testing.T) {
		repo := NewInMemoryCommentsRepository(nil)
		ctx := context.Background()

		comment := domain.NewComment(uuid.New(), uuid.New(), "Random post comment content")
		repo.Store(ctx, comment)

		retrievedComment, err := repo.Find(ctx, comment.ID)

		if assert.NoError(t, err) {
			assert.Equal(t, comment, retrievedComment)
		}
	})
}

func TestInMemoryCommentsRepositoryListAllMethod(t *testing.T) {

	t.Run("it should return empty slice when no comment with given post id was found", func(t *testing.T) {
		repo := NewInMemoryCommentsRepository(nil)
		repo.FillWithRandomComments(uuid.Nil, uuid.Nil, 5)

		postID := uuid.New()

		comments, err := repo.ListAll(context.Background(), postID)
		assert.NoError(t, err)
		assert.Empty(t, comments)
	})

	t.Run("it should all comments for a given post", func(t *testing.T) {
		repo := NewInMemoryCommentsRepository(nil)
		postID := uuid.New()

		repo.FillWithRandomComments(uuid.Nil, uuid.Nil, 7)
		repo.FillWithRandomComments(postID, uuid.Nil, 4)

		comments, err := repo.ListAll(context.Background(), postID)
		_assert := assert.New(t)

		if _assert.NoError(err) {
			for _, comment := range comments {
				_assert.Equal(postID, comment.PostID)
				_assert.NotEqual(uuid.Nil, comment.UserID)
				_assert.NotEmpty(comment.Content)
			}
		}
	})
}

func TestInMemoryCommentsRepositoryListMethod(t *testing.T) {
	t.Run("it should return empty slice when repo is empty", func(t *testing.T) {
		repo := NewInMemoryCommentsRepository(nil)
		
		comments, cursor, err := repo.List(context.Background(), uuid.New(), domain.Paginator{})
		_assert := assert.New(t)

		_assert.Empty(comments)
		_assert.Nil(cursor)
		_assert.NoError(err)
	})

	t.Run("it should return 10 comments when no After is provided", func(t *testing.T) {
		repo := NewInMemoryCommentsRepository(nil)
		postID := uuid.New()
		repo.FillWithRandomComments(postID, uuid.Nil, 15)

		comments, cursor, err := repo.List(context.Background(), postID, domain.Paginator{ Limit: 10 })
		_assert := assert.New(t)

		if _assert.NoError(err) {
			_assert.Equal(10, len(comments))
			_assert.NotNil(cursor)

			last := comments[len(comments)-1]

			_assert.Equal(last.ID, cursor.ID)
			_assert.Equal(last.CreatedAt, cursor.CreatedAt)
		}
	})

	t.Run("it should return the comments after the After element comment", func(t *testing.T) {
		repo := NewInMemoryCommentsRepository(nil)
		ctx := context.Background()
		postID := uuid.New()
		_assert := assert.New(t)

		repo.FillWithRandomComments(postID, uuid.Nil, 9)
		comment := domain.NewComment(postID, uuid.New(), "Random comment content")
		repo.Store(ctx, comment)
		repo.FillWithRandomComments(postID, uuid.Nil, 6)

		prevList, prevCursor, err := repo.List(ctx, postID, domain.Paginator{ Limit: 10 })
		
		if _assert.NoError(err) && _assert.Equal(10, len(prevList)) {
			_assert.Equal(comment.ID, prevCursor.ID)
			_assert.Equal(comment.CreatedAt, prevCursor.CreatedAt)

			currList, currCursor, err := repo.List(ctx, postID, domain.Paginator{ Limit: 10, After: prevCursor })

			if _assert.NoError(err) {
				_assert.Equal(6, len(currList))
				_assert.NotNil(currCursor)
			}
		}
	})
}

func TestInMemoryCommentsRepositoryDeleteMethod(t *testing.T) {
	repo := NewInMemoryCommentsRepository(nil)
	ctx := context.Background()
	comment := domain.NewComment(uuid.New(), uuid.New(), "Random comment content")

	repo.Store(ctx, comment)

	err := repo.Delete(ctx, comment)
	assert.NoError(t, err)
}
