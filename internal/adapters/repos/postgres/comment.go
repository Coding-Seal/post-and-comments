package postgres

import (
	"context"
	"database/sql"
	"errors"
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
	Parent    models.ID
}

type CommentStore struct {
	db *sql.DB
}

func NewCommentStore(db *sql.DB) *CommentStore {
	return &CommentStore{
		db: db,
	}
}

func (c *CommentStore) Comments(ctx context.Context, postID models.ID, offset, limit, nestedLimit int) ([]*models.Comment, error) {
	var comments []Comment

	rows, err := c.db.QueryContext(ctx, "SELECT * FROM comments WHERE post_id = $1 LIMIT $2 OFFSET $3", postID, limit, offset)
	if err != nil {
		return nil, errors.Join(ports.ErrInternal, err)
	}

	for rows.Next() {
		var com Comment
		if err := rows.Scan(&com.ID, &com.Author, &com.Published, &com.Text, &com.Post, &com.Parent); err != nil {
			return nil, errors.Join(ports.ErrInternal, err)
		}

		comments = append(comments, com)
	}

	rootComments := make([]*models.Comment, 0, len(comments))
	for _, com := range comments {
		rootComments = append(rootComments, toComment(com))
	}

	c.buildCommentTree(ctx, rootComments, nestedLimit)

	return rootComments, nil
}

func (c *CommentStore) CommentsAfter(ctx context.Context, postID, afterCommentID models.ID, offset, limit, nestedLimit int) ([]*models.Comment, error) {
	var comments []Comment

	rows, err := c.db.QueryContext(ctx, "SELECT * FROM comments WHERE parent_id = $1 LIMIT $2 OFFSET $3", afterCommentID, limit, offset)
	if err != nil {
		return nil, errors.Join(ports.ErrInternal, err)
	}

	for rows.Next() {
		var com Comment
		if err := rows.Scan(&com.ID, &com.Author, &com.Published, &com.Text, &com.Post, &com.Parent); err != nil {
			return nil, errors.Join(ports.ErrInternal, err)
		}

		comments = append(comments, com)
	}

	offset, limit = min(offset, len(comments)), min(offset+limit, len(comments))

	comments = comments[offset:limit]
	rootComments := make([]*models.Comment, 0, len(comments))

	for _, com := range comments {
		rootComments = append(rootComments, toComment(com))
	}

	c.buildCommentTree(ctx, rootComments, nestedLimit)

	return rootComments, nil
}

func (c *CommentStore) AddComment(ctx context.Context, comment *models.Comment, parent models.ID) error {
	comment.ID = models.NewID()
	com := &Comment{
		ID:        comment.ID,
		Author:    comment.AuthorID,
		Published: comment.Published,
		Text:      comment.Text,
		Post:      comment.PostID,
		Parent:    parent,
	}

	_, err := c.db.ExecContext(ctx, "INSERT INTO comments (id, author_id, published, text, post_id, parent_id) VALUES ($1, $2, $3, $4, $5, $6)", com.ID, com.Author, com.Published, com.Text, com.Post, com.Parent)
	if err != nil {
		return errors.Join(err, ports.ErrInternal)
	}

	return nil
}

var _ ports.CommentStore = (*CommentStore)(nil)

func (c *CommentStore) buildCommentTree(ctx context.Context, comments []*models.Comment, nestedLimit int) {
	if nestedLimit == 0 {
		return
	}

	for _, comment := range comments {
		children, err := c.getChildren(ctx, comment.ID)
		if err != nil {
			return
		}

		for _, child := range children {
			comment.Replies = append(comment.Replies, toComment(child))
		}

		c.buildCommentTree(ctx, comment.Replies, nestedLimit-1)
	}
}

func (c *CommentStore) getChildren(ctx context.Context, parentID models.ID) ([]Comment, error) {
	var comments []Comment

	rows, err := c.db.QueryContext(ctx, "SELECT * FROM comments WHERE parent_id = $1", parentID)
	if err != nil {
		return nil, errors.Join(ports.ErrInternal, err)
	}

	for rows.Next() {
		var com Comment
		if err := rows.Scan(&com.ID, &com.Author, &com.Published, &com.Text, &com.Post, &com.Parent); err != nil {
			return nil, errors.Join(ports.ErrInternal, err)
		}

		comments = append(comments, com)
	}

	return comments, nil
}

func toComment(comment Comment) *models.Comment {
	return &models.Comment{
		ID:        comment.ID,
		Published: comment.Published,
		Text:      comment.Text,
		PostID:    comment.Post,
		Replies:   nil,
	}
}
