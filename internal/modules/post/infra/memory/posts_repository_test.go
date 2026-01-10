package memory

import (
	"comu/internal/modules/post/domain"
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryPostsRepositoryStoreMethod(t *testing.T) {
	repo := NewInMemoryPostsRepository(nil)
	ctx := context.Background()

	for i := range 5 {
		title := fmt.Sprintf("Post #%d", i)
		content := fmt.Sprintf("%s content", title)
		post := domain.NewPost(uuid.New(), title, content)

		err := repo.Store(ctx, post)
		assert.NoError(t, err)
	}
}

func TestInMemoryPostsRepositoryFindByIdAndFindBySlugMethod(t *testing.T) {
	t.Run("FindByID should fail and return ErrPostNotFound", func(t *testing.T) {
		repo := NewInMemoryPostsRepository(nil)
		
		post, err := repo.FindByID(context.Background(), uuid.New())
		assert.ErrorIs(t, err, domain.ErrPostNotFound)
		assert.Nil(t, post)
	})

	t.Run("FindBySlug should fail and return ErrPostNotFound", func(t *testing.T) {
		repo := NewInMemoryPostsRepository(nil)
		slug := "post-#1"

		post, err := repo.FindBySlug(context.Background(), slug)
		assert.ErrorIs(t, err, domain.ErrPostNotFound)
		assert.Nil(t, post)
	})

	t.Run("FindByID should succeed and return the post", func(t *testing.T) {
		repo := NewInMemoryPostsRepository(nil)
		ctx := context.Background()

		post := domain.NewPost(uuid.New(), "Post #34", "Post #34 content")
		repo.Store(ctx, post)
		
		retrievedPost, err := repo.FindByID(ctx, post.ID)
		_assert := assert.New(t)
		
		if _assert.NoError(err) {
			_assert.Equal(post.Title, retrievedPost.Title)
			_assert.Equal(post.Content, retrievedPost.Content)
		}
	})

	t.Run("FindBySlug should succeed and return the post", func(t *testing.T) {
		repo := NewInMemoryPostsRepository(nil)
		ctx := context.Background()

		post := domain.NewPost(uuid.New(), "Post #34", "Post #34 content")
		repo.Store(ctx, post)
		
		retrievedPost, err := repo.FindBySlug(ctx, post.Slug)
		_assert := assert.New(t)
		
		if _assert.NoError(err) {
			_assert.Equal(post.Title, retrievedPost.Title)
			_assert.Equal(post.Content, retrievedPost.Content)
		}
	})
}

func TestInMemoryPostsRepositoryListMethod(t *testing.T) {
	t.Run("it should return empty slice when repo is empty", func(t *testing.T) {
		repo := NewInMemoryPostsRepository(nil)
		
		posts, cursor, err := repo.List(context.Background(), domain.Paginator{})
		_assert := assert.New(t)

		_assert.Empty(posts)
		_assert.Nil(cursor)
		_assert.NoError(err)
	})

	t.Run("it should return 12 posts when no After is provided", func(t *testing.T) {
		repo := NewInMemoryPostsRepository(nil)
		repo.FillWithRandomPosts(uuid.Nil, 18)

		posts, cursor, err := repo.List(context.Background(), domain.Paginator{ Limit: 12 })
		_assert := assert.New(t)

		if _assert.NoError(err) {
			_assert.Equal(12, len(posts))
			_assert.NotNil(cursor)

			last := posts[len(posts)-1]

			_assert.Equal(last.ID, cursor.ID)
			_assert.Equal(last.CreatedAt, cursor.CreatedAt)
		}
	})

	t.Run("it should return the remaining posts after the After element post", func(t *testing.T) {
		repo := NewInMemoryPostsRepository(nil)
		ctx := context.Background()
		_assert := assert.New(t)

		repo.FillWithRandomPosts(uuid.Nil, 9)
		post := domain.NewPost(uuid.New(), "Test post", "Random content")
		repo.Store(ctx, post)
		repo.FillWithRandomPosts(uuid.Nil, 8)

		prevList, prevCursor, err := repo.List(ctx, domain.Paginator{ Limit: 10 })
		
		if _assert.NoError(err) && _assert.Equal(10, len(prevList)) {
			_assert.Equal(post.ID, prevCursor.ID)
			_assert.Equal(post.CreatedAt, prevCursor.CreatedAt)

			currList, currCursor, err := repo.List(ctx, domain.Paginator{ Limit: 10, After: prevCursor })

			if _assert.NoError(err) {
				_assert.Equal(8, len(currList))
				_assert.NotNil(currCursor)
			}
		}
	})

}

func TestInMemoryPostsRepositoryDeleteMethod(t *testing.T) {
	repo := NewInMemoryPostsRepository(nil)
	post := domain.NewPost(uuid.New(), "Random Post", "This is a random content")
	ctx := context.Background()

	repo.Store(ctx, post)
	err := repo.Delete(ctx, post)

	if assert.NoError(t, err) {
		_, err := repo.FindByID(ctx, post.ID)
		assert.ErrorIs(t, err, domain.ErrPostNotFound)
	}
}
