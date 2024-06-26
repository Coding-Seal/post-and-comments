package outports

import (
	"context"
	"errors"

	"post-and-comments/internal/models"
)

var ErrUserNotFound = errors.New("user not found")

type UserStore interface {
	User(ctx context.Context, id models.ID) (*models.User, error)
	AddUser(ctx context.Context, user *models.User) error
}
