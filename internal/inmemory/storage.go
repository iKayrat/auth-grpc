package storage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"sync"

	"github.com/google/uuid"
	"github.com/iKayrat/auth-grpc/internal/model"
)

type UserStorage struct {
	Users map[string]model.User
	mu    sync.RWMutex
}

func New(us *UserStorage) Store {
	return &UserStorage{
		Users: make(map[string]model.User),
	}
}

// NewUS returns pointer to UserStorage
func NewUS() *UserStorage {
	return &UserStorage{
		Users: make(map[string]model.User, 0),
	}
}

func InitStorage() *UserStorage {
	storage := NewUS()

	storage.mu.Lock()
	defer storage.mu.Unlock()

	adminUUID, _ := uuid.Parse("e257152c-43e2-4fe2-8cd3-0618006890c5")

	uid1, _ := uuid.Parse("1c539b0a-7d71-418b-9136-dbcbccd7908b")
	uid2 := uuid.New()

	// initial users
	initialUsers := []model.User{
		{
			ID:       adminUUID,
			Email:    "admin@admin",
			Username: "a",
			Password: "12",
			Admin:    true,
		},
		{
			ID:       uid1,
			Email:    "a@yandex.ru",
			Username: "Alisa",
			Password: "yandex123",
			Admin:    false,
		},
		{
			ID:       uid2,
			Email:    "m@mail.ru",
			Username: "Marusya",
			Password: "vcontact21",
			Admin:    false,
		},
	}

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

func (us *UserStorage) GetAll(ctx context.Context) ([]model.User, error) {
	us.mu.RLock()
	defer us.mu.RUnlock()
	if len(us.Users) == 0 {
		return nil, errors.New(UserNotFound)
	}

	users := make([]model.User, 0)
	for _, u := range us.Users {
		users = append(users, model.User{
			ID:       u.ID,
			Email:    u.Email,
			Username: u.Username,
		})
	}

	//sort slice
	sort.Slice(users, func(i, j int) bool {
		return users[i].ID.ID() < users[j].ID.ID()
	})

	return users, nil
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

	return model.User{}, errors.New(UserNotFound)
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

func (us *UserStorage) Delete(ctx context.Context, id string) error {
	us.mu.Lock()
	defer us.mu.Unlock()

	_, exists := us.Users[id]
	if !exists {
		return fmt.Errorf(UserWithIdNotFound, id)
	}

	delete(us.Users, id)

	return nil
}
