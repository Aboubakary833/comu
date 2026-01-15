package memory

import (
	"comu/internal/modules/post/domain"
	"context"
	"fmt"
	"maps"
	"slices"
	"sync"
	"time"

	"github.com/google/uuid"
)

type postStore map[uuid.UUID]domain.Post

type inMemoryPostsRepository struct {
	store postStore
	sync.Mutex
}

func NewInMemoryPostsRepository(initialStore postStore) *inMemoryPostsRepository {

	if initialStore == nil {
		initialStore = make(postStore)
	}

	return &inMemoryPostsRepository{
		store: initialStore,
	}
}

func (repo *inMemoryPostsRepository) ListAll(ctx context.Context) ([]domain.Post, error) {
	repo.Lock()
	defer repo.Unlock()

	posts := slices.Collect(maps.Values(repo.store))
	sortPosts(posts)

	return posts, nil
}

func (repo *inMemoryPostsRepository) List(ctx context.Context, paginator domain.Paginator) ([]domain.Post, *domain.Cursor, error) {
	allPosts, err := repo.ListAll(ctx)

	if err != nil {
		return []domain.Post{}, nil, err
	}

	if len(allPosts) == 0 {
		return allPosts, nil, nil
	}

	if paginator.After == nil {
		return repo.listReturnValues(allPosts[:paginator.Limit])
	}

	afterIdx := slices.IndexFunc(allPosts, func(post domain.Post) bool {
		return post.ID == paginator.After.ID
	})

	if afterIdx == -1 {
		posts := allPosts[:domain.DefaultPaginatorLimit]
		return repo.listReturnValues(posts)
	}
	posts := allPosts[afterIdx+1:]

	if afterIdx+paginator.Limit > len(posts) {
		return repo.listReturnValues(posts)
	}

	return repo.listReturnValues(posts[:paginator.Limit])
}

func (repo *inMemoryPostsRepository) listReturnValues(posts []domain.Post) ([]domain.Post, *domain.Cursor, error) {
	last := posts[len(posts)-1]
	return posts, &domain.Cursor{ID: last.ID, CreatedAt: last.CreatedAt}, nil
}

func (repo *inMemoryPostsRepository) FindByID(ctx context.Context, ID uuid.UUID) (*domain.Post, error) {
	repo.Lock()
	defer repo.Unlock()

	if post, ok := repo.store[ID]; ok {
		return &post, nil
	}

	return nil, domain.ErrPostNotFound
}

func (repo *inMemoryPostsRepository) FindBySlug(ctx context.Context, slug string) (*domain.Post, error) {
	repo.Lock()
	defer repo.Unlock()

	for _, post := range repo.store {
		if post.Slug == slug {
			return &post, nil
		}
	}

	return nil, domain.ErrPostNotFound
}

func (repo *inMemoryPostsRepository) Store(ctx context.Context, post *domain.Post) error {
	repo.Lock()
	defer repo.Unlock()

	id, err := uuid.NewV7()

	if err != nil {
		return err
	}
	post.ID = id
	repo.store[id] = *post

	return nil
}

func (repo *inMemoryPostsRepository) Update(ctx context.Context, post *domain.Post) error {
	_, err := repo.FindByID(ctx, post.ID)

	if err != nil {
		return err
	}
	copy := *post
	copy.UpdatedAt = time.Now()
	repo.store[copy.ID] = copy

	return nil
}

func (repo *inMemoryPostsRepository) Delete(ctx context.Context, post *domain.Post) error {
	post, err := repo.FindByID(ctx, post.ID)

	if err != nil {
		return err
	}

	delete(repo.store, post.ID)
	return nil
}

// FillWithRandomPosts is a test factory method which main purpose is
// to generate a number of random posts and store them in the repo.
func (repo *inMemoryPostsRepository) FillWithRandomPosts(userID uuid.UUID, length int) {
	ctx := context.Background()

	if userID == uuid.Nil {
		userID = uuid.New()
	}
	for i := range length {
		title := fmt.Sprintf("Post #%d", i+1)
		content := fmt.Sprintf("%s content", title)
		post := domain.NewPost(userID, title, content)

		repo.Store(ctx, post)
	}
}

func sortPosts(posts []domain.Post) {
	slices.SortFunc(posts, func(a domain.Post, b domain.Post) int {
		if a.CreatedAt.Before(b.CreatedAt) {
			return -1
		}
		if a.CreatedAt.After(b.CreatedAt) {
			return 1
		}

		return 0
	})
}
