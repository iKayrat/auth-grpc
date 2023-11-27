package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/iKayrat/auth-grpc/internal/model"
)

type Store interface {
	Create(ctx context.Context, user model.User) (model.User, error)
	GetAll(ctx context.Context) ([]model.User, error)
	GetById(ctx context.Context, id uuid.UUID) (model.User, bool)
	GetByUsername(ctx context.Context, username string) (model.User, error)
	Update(ctx context.Context, updateUser model.User) (model.User, error)
	Delete(ctx context.Context, id string) error
}
