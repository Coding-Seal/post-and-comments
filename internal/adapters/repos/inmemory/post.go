package inmemory

import (
	"context"
	"slices"
	"time"

	ports "post-and-comments/internal/inner-ports"
	"post-and-comments/internal/models"
)

type PostStore struct {
	posts map[models.ID]*models.Post
}

func (p *PostStore) DisableComments(ctx context.Context, postID models.ID) error {
	if _, exists := p.posts[postID]; !exists {
		return ports.ErrPostNotFound
	}

	p.posts[postID].DisableComments = true

	return nil
}

func (p *PostStore) Posts(ctx context.Context, offset, limit int) ([]*models.Post, error) {
	posts := make([]*models.Post, 0, len(p.posts))
	for _, post := range p.posts {
		posts = append(posts, post)
	}

	slices.SortFunc(posts, func(lhs, rhs *models.Post) int {
		if lhs.ID.String() < rhs.ID.String() {
			return -1
		} else if lhs.ID.String() > rhs.ID.String() {
			return 1
		}

		return 0
	})

	offset, limit = min(offset, len(posts)), min(offset+limit, len(posts))

	return posts[offset:limit], nil
}

func (p *PostStore) Post(ctx context.Context, postID models.ID) (*models.Post, error) {
	post, exists := p.posts[postID]
	if !exists {
		return nil, ports.ErrPostNotFound
	}

	return post, nil
}

func (p *PostStore) AddPost(ctx context.Context, post *models.Post) error {
	id := models.NewID()
	post.ID = id
	post.Published = time.Now()
	p.posts[id] = post

	return nil
}

var _ ports.PostStore = (*PostStore)(nil)

func NewPostStore() *PostStore {
	return &PostStore{
		posts: make(map[models.ID]*models.Post),
	}
}
