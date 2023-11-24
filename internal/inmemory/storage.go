package storage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/iKayrat/auth-grpc/internal/model"
)

type UserStorage struct {
	Users map[string]model.User
	mu    sync.RWMutex
}

func New() Store {
	return &UserStorage{
		Users: make(map[string]model.User),
	}
}

func NewUS() *UserStorage {
	return &UserStorage{
		Users: make(map[string]model.User),
	}
}

func InitStorage() *UserStorage {
	storage := NewUS()

	adminUUID, _ := uuid.Parse("e257152c-43e2-4fe2-8cd3-0618006890c5")

	// initial users
	initialUsers := []model.User{
		{
			ID:       adminUUID,
			Email:    "admin@admin",
			Username: "admin",
			Password: "admin",
			Admin:    true,
		},
		{
			ID:       uuid.New(),
			Email:    "a@yandex.ru",
			Username: "Alisa",
			Password: "yandex123",
			Admin:    false,
		},
		{
			ID:       uuid.New(),
			Email:    "m@mail.ru",
			Username: "Marusya",
			Password: "vcontact21",
			Admin:    false,
		},
	}

	storage.mu.Lock()
	defer storage.mu.Unlock()
	for _, user := range initialUsers {
		id := user.ID.String()
		storage.Users[id] = user
	}

	return storage
}

func (us *UserStorage) Create(ctx context.Context, user model.User) (model.User, error) {
	us.mu.Lock()
	defer us.mu.Unlock()

	id, _ := uuid.NewUUID()

	idd := id.String()
	log.Println("id", idd)

	user.ID = id
	us.Users[idd] = user

	log.Println("users[map]", us.Users)

	return us.Users[idd], nil
}

func (us *UserStorage) GetAll(ctx context.Context) []model.User {
	us.mu.RLock()
	defer us.mu.RUnlock()

	users := make([]model.User, len(us.Users))
	for _, u := range us.Users {
		users = append(users, u)
	}

	return users
}

func (us *UserStorage) GetById(ctx context.Context, id uuid.UUID) (model.User, bool) {
	us.mu.RLock()
	defer us.mu.RUnlock()

	user, exists := us.Users[id.String()]

	return user, exists
}

func (us *UserStorage) GetByUsername(ctx context.Context, username string) (model.User, error) {
	us.mu.RLock()
	defer us.mu.RUnlock()

	for _, v := range us.Users {
		if v.Username == username {
			return v, nil
		}
	}

	return model.User{}, errors.New("user not found")
}

func (us *UserStorage) Update(ctx context.Context, updateUser model.User) (model.User, error) {
	us.mu.Lock()
	defer us.mu.Unlock()

	_, exists := us.Users[updateUser.ID.String()]
	fmt.Println("us.Users[updateUser.ID.String()]:", us.Users[updateUser.ID.String()])
	if !exists {
		return model.User{}, fmt.Errorf("user with ID %s not found", updateUser.ID.String())
	}

	us.Users[updateUser.ID.String()] = updateUser

	return us.Users[updateUser.ID.String()], nil
}

func (us *UserStorage) Delete(ctx context.Context, id uuid.UUID) error {
	us.mu.Lock()
	defer us.mu.Unlock()

	_, exists := us.Users[id.String()]
	if !exists {
		return fmt.Errorf("user with ID %s not found", id.String())
	}

	return nil
}
