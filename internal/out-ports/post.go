package out_ports

import (
	"context"

	"post-and-comments/internal/models"
)

type PostService interface {
	DisableComments(ctx context.Context, postID models.ID) error
	Posts(ctx context.Context, offset, limit int) ([]*models.Post, error)
	Post(ctx context.Context, postID models.ID) (*models.Post, error)
	AddPost(ctx context.Context, post *models.Post) error
}
