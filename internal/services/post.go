package services

import (
	ports "post-and-comments/internal/inner-ports"
	outports "post-and-comments/internal/out-ports"
)

type Post struct {
	ports.PostStore
}

var _ outports.PostService = (*Post)(nil)
