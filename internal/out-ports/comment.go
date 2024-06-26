package out_ports

import (
	"context"

	"post-and-comments/internal/models"
)

type CommentService interface {
	Comments(ctx context.Context, postID models.ID, offset, limit, nestedLimit int) ([]*models.Comment, error)
	CommentsAfter(ctx context.Context, postID, afterCommentID models.ID, offset, limit, nestedLimit int) ([]*models.Comment, error)
	AddComment(ctx context.Context, comment *models.Comment, parent models.ID) error
}
