package out_ports

import (
	"context"

	"post-and-comments/internal/models"
)

type UserService interface {
	User(ctx context.Context, id models.ID) (*models.User, error)
	AddUser(ctx context.Context, user *models.User) error
}
