package memory

import (
	"comu/internal/modules/users/domain"
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type userStore map[uuid.UUID]domain.User

type inMemoryRepository struct {
	users userStore
	sync.Mutex
}

func NewInMemoryRepository(initialStore userStore) *inMemoryRepository {
	if initialStore == nil {
		initialStore = make(userStore)
	}

	return &inMemoryRepository{
		users: initialStore,
	}
}

func (repo *inMemoryRepository) FindByID(ctx context.Context, ID uuid.UUID) (*domain.User, error) {
	repo.Lock()
	defer repo.Unlock()

	user, ok := repo.users[ID]

	if !ok {
		return nil, domain.ErrUserNotFound
	}

	return &user, nil
}

func (repo *inMemoryRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	repo.Lock()
	defer repo.Unlock()

	for _, user := range repo.users {
		if user.Email == email {
			return &user, nil
		}
	}

	return nil, domain.ErrUserNotFound
}

func (repo *inMemoryRepository) Store(ctx context.Context, user *domain.User) error {

	if repo.emailIsTaken(user.Email) {
		return domain.ErrUserEmailTaken
	}

	repo.Lock()
	defer repo.Unlock()

	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	user.ID = id

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	repo.users[id] = *user

	return nil
}

func (repo *inMemoryRepository) Update(ctx context.Context, user *domain.User) error {

	savedUser, err := repo.FindByID(ctx, user.ID)

	if err != nil {
		return err
	}

	if savedUser.Email != user.Email && repo.emailIsTaken(user.Email) {
		return domain.ErrUserEmailTaken
	}
	repo.Lock()
	defer repo.Unlock()

	if user.CreatedAt != savedUser.CreatedAt {
		user.CreatedAt = savedUser.CreatedAt
	}
	user.UpdatedAt = time.Now()
	repo.users[user.ID] = *user

	return nil
}

func (repo *inMemoryRepository) Delete(ctx context.Context, user *domain.User) error {
	_, err := repo.FindByID(ctx, user.ID)

	if err != nil {
		return err
	}

	repo.Lock()
	defer repo.Unlock()

	delete(repo.users, user.ID)

	return nil
}

func (repo *inMemoryRepository) emailIsTaken(email string) bool {
	repo.Lock()
	defer repo.Unlock()

	for _, user := range repo.users {
		if user.Email == email {
			return true
		}
	}

	return false
}
