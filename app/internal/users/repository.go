package users

import (
	"context"
	"model"

	"github.com/google/uuid"
)

type repository interface {
	createUser(ctx context.Context, user *model.User) error
	getUser(ctx context.Context, id uuid.UUID) (*model.User, error)
	updateUser(ctx context.Context, user *model.User) error
	deleteUser(ctx context.Context, id uuid.UUID) error
	listUsers(ctx context.Context, filter UserFilter) ([]*model.User, int64, error)
	getUserByEmail(ctx context.Context, email string) (*model.User, error)
	getUserByUsername(ctx context.Context, username string) (*model.User, error)
}

type UserFilter struct {
	Username string
	Email    string
	Status   model.StatusEnum
	Page     int
	PageSize int
}
