package services

import (
	ports "post-and-comments/internal/inner-ports"
	outports "post-and-comments/internal/out-ports"
)

type Comment struct {
	ports.CommentStore
}

var _ outports.CommentService = (*Comment)(nil)
