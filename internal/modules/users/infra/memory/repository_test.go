package memory

import (
	"comu/internal/modules/users/domain"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryRepositoryStoreMethod(t *testing.T) {

	t.Run("repo.Store should successfully add a new user entry", func(t *testing.T) {
		repo := NewInMemoryRepository(nil)
		newUser := domain.NewUser("John Doe", "johndoe@gmail.com", "uUu4MRrxPBhVmqUnb6m1upSh")
		_assert := assert.New(t)

		result := repo.Store(context.Background(),newUser)

		if _assert.Nil(result) {
			_assert.NotEmpty(newUser.ID.Value())
			_assert.Equal(newUser.Name, "John Doe")
			_assert.Equal(newUser.Email, "johndoe@gmail.com")
			_assert.NotNil(newUser.CreatedAt)
			_assert.NotNil(newUser.UpdatedAt)
		}
	})

	t.Run("repo.Store should fail and return ErrUserEmailTaken", func(t *testing.T) {
		id := uuid.New()
		existingUser := domain.NewUser("Markus Lacharette", "marcuschevalier@gmail.com", "a0XYkQKxOUQ1zwY2iIB5y4Ci")
		existingUser.SetID(id)

		repo := NewInMemoryRepository(userStore{id: *existingUser})
		newUser := domain.NewUser("Marcus Lacharette", "marcuschevalier@gmail.com", "GyEqfiokhRGb4Ci19")

		result := repo.Store(context.Background(), newUser)

		assert.Equal(t, domain.ErrUserEmailTaken, result)
	})
}

func TestInMemoryRepositoryFindMethod(t *testing.T) {

	t.Run("repo.Find should successfully return the user with the given id", func(t *testing.T) {
		id := uuid.New()
		user := domain.NewUser("Marmoush Dia", "marmoush125@gmail.com", "kQKxOUQ1zwY2iIB5y4")
		user.SetID(id)

		repo := NewInMemoryRepository(userStore{id: *user})
		_assert := assert.New(t)

		result, err := repo.FindByID(context.Background(), id)

		if _assert.Nil(err) {
			_assert.Equal(user.Name, result.Name)
			_assert.Equal(user.Email, result.Email)
		}
	})

	t.Run("repo.Find should fail and return ErrUserNotFound", func(t *testing.T) {
		repo := NewInMemoryRepository(nil)
		id := uuid.New()

		_, err := repo.FindByID(context.Background(), id)

		assert.Equal(t, domain.ErrUserNotFound, err)
	})
}

func TestInMemoryRepositoryUpdateMethod(t *testing.T) {

	t.Run("repo.Update should successfully update a given user information", func(t *testing.T) {
		id := uuid.New()
		subjectUser := domain.NewUser("John Doe", "johndoe@gmail.com", "uUu4MRrxPBhVmqUnb6m1upSh")
		subjectUser.SetID(id)

		repo := NewInMemoryRepository(userStore{id: *subjectUser})

		subjectUser.Name = "Johnathan ThePopFan"
		subjectUser.Email = "popfan989@gmail.com"

		_assert := assert.New(t)
		expected := repo.Update(context.Background(), subjectUser)

		if _assert.Nil(expected) {
			_assert.Equal(subjectUser.Name, "Johnathan ThePopFan")
			_assert.Equal(subjectUser.Email, "popfan989@gmail.com")
			_assert.NotEqual(subjectUser.UpdatedAt, subjectUser.CreatedAt)
		}
	})

	t.Run("repo.Update should fail and return ErrUserNotFound", func(t *testing.T) {
		repo := NewInMemoryRepository(nil)
		user := domain.NewUser("Frank Deschamps", "frank450@gmail.com", "qiEA80snbhpuDFleTv5fpE")
		user.SetID(uuid.New())

		user.Name = "Franklin Deschamps"

		result := repo.Update(context.Background(), user)

		assert.Equal(t, domain.ErrUserNotFound, result)
	})

	t.Run("repo.Update should fail and return ErrUserEmailTaken", func(t *testing.T) {
		existingUser := domain.NewUser("Marcus Voilier", "marcuschevalier@gmail.com", "a0XYkQKxOUQ1zwY2iIB5y4Ci")
		subjectUser := domain.NewUser("Frank Deschamps", "frank450@gmail.com", "qiEA80snbhpuDFleTv5fpE")
		existingUserId := uuid.New()
		subjectUserId := uuid.New()

		existingUser.SetID(existingUserId)
		subjectUser.SetID(subjectUserId)

		repo := NewInMemoryRepository(userStore{
			existingUserId: *existingUser,
			subjectUserId:  *subjectUser,
		})

		subjectUser.Email = "marcuschevalier@gmail.com"

		result := repo.Update(context.Background(), subjectUser)

		assert.Equal(t, domain.ErrUserEmailTaken, result)
	})

}

func TestInMemoryRepositoryDeleteMethod(t *testing.T) {

	t.Run("repo.Delete should successfully delete a given user", func(t *testing.T) {
		subjectUser := domain.NewUser("Michael Johnson", "michiavel002@gmail.com", "o7RwVfoIHkAWUpnUE7j")
		id := uuid.New()
		subjectUser.SetID(id)

		repo := NewInMemoryRepository(userStore{id: *subjectUser})
		ctx := context.Background()

		result := repo.Delete(ctx, subjectUser)
		_assert := assert.New(t)

		if _assert.Nil(result) {
			_, err := repo.FindByID(ctx, id)
			_assert.Equal(domain.ErrUserNotFound, err)
		}
	})

	t.Run("repo.Delete should fail and return ErrUserNotFound", func(t *testing.T) {
		subjectUser := domain.NewUser("Michael Johnson", "michiavel002@gmail.com", "o7RwVfoIHkAWUpnUE7j")
		id := uuid.New()
		subjectUser.SetID(id)

		repo := NewInMemoryRepository(nil)

		result := repo.Delete(context.Background(), subjectUser)

		assert.Equal(t, domain.ErrUserNotFound, result)
	})
}
