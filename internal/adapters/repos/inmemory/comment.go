package inmemory

import (
	"context"
	"time"

	ports "post-and-comments/internal/inner-ports"
	"post-and-comments/internal/models"
)

type Comment struct {
	ID        models.ID
	Author    models.ID
	Published time.Time
	Text      models.CommentText
	Post      models.ID
	Parent    *Comment
}

type CommentStore struct {
	postComments   map[models.ID][]*Comment
	parentComments map[models.ID][]*Comment
}

func NewCommentStore() *CommentStore {
	return &CommentStore{
		postComments:   make(map[models.ID][]*Comment),
		parentComments: make(map[models.ID][]*Comment),
	}
}

func (c *CommentStore) Comments(ctx context.Context, postID models.ID, offset, limit, nestedLimit int) ([]*models.Comment, error) {
	comments, ok := c.postComments[postID]
	if !ok {
		return nil, nil
	}

	var rootComments []*models.Comment

	for _, comment := range comments {
		if comment.Parent == nil {
			com, err := c.toComment(ctx, comment)
			if err == nil {
				rootComments = append(rootComments, com)
			}
		}
	}

	offset, limit = min(offset, len(rootComments)), min(offset+limit, len(rootComments))
	rootComments = rootComments[offset:limit]
	c.buildCommentTree(ctx, rootComments, nestedLimit)

	return rootComments, nil
}

func (c *CommentStore) CommentsAfter(ctx context.Context, postID, afterCommentID models.ID, offset, limit, nestedLimit int) ([]*models.Comment, error) {
	var rootComments []*models.Comment

	comments, ok := c.parentComments[afterCommentID]
	if !ok {
		return nil, nil
	}

	for _, comment := range comments {
		com, err := c.toComment(ctx, comment)
		if err == nil {
			rootComments = append(rootComments, com)
		}
	}

	offset, limit = min(offset, len(rootComments)), min(offset+limit, len(rootComments))
	rootComments = rootComments[offset:limit]
	c.buildCommentTree(ctx, rootComments, nestedLimit)

	return rootComments, nil
}

func (c *CommentStore) AddComment(ctx context.Context, comment *models.Comment, parent models.ID) error {
	parents, ok := c.postComments[comment.ID]
	if !ok {
		return ports.ErrPostNotFound
	}

	var parentExists bool

	var parentComment *Comment

	for _, p := range parents {
		if p.ID == parent {
			parentExists = true
			parentComment = p

			break
		}
	}

	if !parentExists {
		return ports.ErrWrongParentComment
	}

	comment.ID = models.NewID()
	com := &Comment{
		ID:        comment.ID,
		Author:    comment.AuthorID,
		Published: comment.Published,
		Text:      comment.Text,
		Post:      comment.PostID,
		Parent:    parentComment,
	}
	c.postComments[comment.PostID] = append(c.postComments[comment.PostID], com)
	c.parentComments[parent] = append(c.parentComments[parent], com)

	return nil
}

var _ ports.CommentStore = (*CommentStore)(nil)

func (c *CommentStore) buildCommentTree(ctx context.Context, comments []*models.Comment, nestedLimit int) {
	if nestedLimit == 0 {
		return
	}

	for _, m := range comments {
		children := c.parentComments[m.ID]
		for _, child := range children {
			com, err := c.toComment(ctx, child)
			if err == nil {
				m.Replies = append(m.Replies, com)
			}
		}

		c.buildCommentTree(ctx, m.Replies, nestedLimit-1)
	}
}

func (c *CommentStore) toComment(ctx context.Context, comment *Comment) (*models.Comment, error) {
	return &models.Comment{
		ID:        comment.ID,
		AuthorID:  comment.Author,
		Published: comment.Published,
		Text:      comment.Text,
		PostID:    comment.Post,
		Replies:   nil,
	}, nil
}
