package outports

import (
	"context"
	"errors"

	"post-and-comments/internal/models"
)

var (
	ErrPostNotFound = errors.New("post not found")
	ErrInternal     = errors.New("internal error")
)

type PostStore interface {
	DisableComments(ctx context.Context, postID models.ID) error
	Posts(ctx context.Context, offset, limit int) ([]*models.Post, error)
	Post(ctx context.Context, postID models.ID) (*models.Post, error)
	AddPost(ctx context.Context, post *models.Post) error
}
