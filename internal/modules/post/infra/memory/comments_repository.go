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

type commentStore map[uuid.UUID]domain.Comment

type inMemoryCommentsRepository struct {
	store commentStore
	sync.Mutex
}

func NewInMemoryCommentsRepository(initialStore commentStore) *inMemoryCommentsRepository {
	if initialStore == nil {
		initialStore = make(commentStore)
	}

	return &inMemoryCommentsRepository{
		store: initialStore,
	}
}

func (repo *inMemoryCommentsRepository) Store(ctx context.Context, comment *domain.Comment) error {
	repo.Lock()
	defer repo.Unlock()

	id, err := uuid.NewV7()

	if err != nil {
		return err
	}

	comment.ID = id
	repo.store[id] = *comment

	return nil
}

func (repo *inMemoryCommentsRepository) Update(ctx context.Context, comment *domain.Comment) error {
	_, err := repo.Find(ctx, comment.ID)

	if err != nil {
		return err
	}
	copy := *comment
	copy.UpdatedAt = time.Now()
	repo.store[copy.ID] = copy

	return nil
}

func (repo *inMemoryCommentsRepository) ListAll(ctx context.Context, postID uuid.UUID) ([]domain.Comment, error) {
	comments := filterComments(slices.Collect(maps.Values(repo.store)), func(c domain.Comment) bool {
		return c.PostID == postID
	})
	sortComments(comments)

	return comments, nil
}

func (repo *inMemoryCommentsRepository) List(ctx context.Context, postID uuid.UUID, paginator domain.Paginator) ([]domain.Comment, *domain.Cursor, error) {
	allComments, err := repo.ListAll(ctx, postID)

	if err != nil {
		return []domain.Comment{}, nil, err
	}

	if len(allComments) == 0 {
		return allComments, nil, nil
	}

	if paginator.After == nil {
		return repo.listReturnValues(allComments[:paginator.Limit])
	}

	afterIdx := slices.IndexFunc(allComments, func(c domain.Comment) bool {
		return c.ID == paginator.After.ID
	})

	if afterIdx == -1 {
		comments := allComments[:domain.DefaultPaginatorLimit]
		return repo.listReturnValues(comments)
	}
	comments := allComments[afterIdx + 1:]

	if afterIdx + paginator.Limit > len(comments) {
		return repo.listReturnValues(comments)
	}

	return repo.listReturnValues(comments[:paginator.Limit])
}

func (repo *inMemoryCommentsRepository) listReturnValues(comments []domain.Comment) ([]domain.Comment, *domain.Cursor, error) {
	last := comments[len(comments)-1]
	return comments, &domain.Cursor{ ID: last.ID, CreatedAt: last.CreatedAt }, nil
}

func (repo *inMemoryCommentsRepository) Find(ctx context.Context, ID uuid.UUID) (*domain.Comment, error) {
	repo.Lock()
	defer repo.Unlock()

	if comment, ok := repo.store[ID]; ok {
		return &comment, nil
	}

	return nil, domain.ErrCommentNotFound
}

func (repo *inMemoryCommentsRepository) Delete(ctx context.Context, comment *domain.Comment) error {
	_, err := repo.Find(ctx, comment.ID)

	if err != nil {
		return err
	}
	repo.Lock()
	defer repo.Unlock()

	delete(repo.store, comment.ID)

	return nil
}

// Well, it's like the inMemoryPostsRepository FillWithRandoms method
// but for inMemoryCommentsRepository
func (repo *inMemoryCommentsRepository) FillWithRandomComments(postID, authorID uuid.UUID, length int) {
	ctx := context.Background()

	for i := range length {
		if postID == uuid.Nil {
			postID = uuid.New()
		}

		if authorID == uuid.Nil {
			authorID = uuid.New()
		}

		content := fmt.Sprintf("Post #%s comment #%d content", postID.String(), i + 1)
		repo.Store(ctx, domain.NewComment(postID, authorID, content))
	}
}


func sortComments(comments []domain.Comment) {
	slices.SortFunc(comments, func(a domain.Comment, b domain.Comment) int {
		if a.CreatedAt.Before(b.CreatedAt) { return - 1 }
		if a.CreatedAt.After(b.CreatedAt) { return 1 }
		
		return 0
	})
}

func filterComments(comments []domain.Comment, filterFunc func(domain.Comment) bool) []domain.Comment {
	var filteredComments []domain.Comment

	for _, comment := range comments {
		if filterFunc(comment) {
			filteredComments = append(filteredComments, comment)
		}
	}

	return filteredComments
}
