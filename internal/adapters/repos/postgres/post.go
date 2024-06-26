package postgres

import (
	"context"
	"database/sql"
	"errors"

	ports "post-and-comments/internal/inner-ports"
	"post-and-comments/internal/models"
)

type PostStore struct {
	db *sql.DB
}

func (p *PostStore) DisableComments(ctx context.Context, postID models.ID) error {
	_, err := p.db.ExecContext(ctx, "UPDATE posts SET  disable_comments = false WHERE id = $1", postID)
	if err != nil {
		return errors.Join(ports.ErrInternal, err)
	}

	return nil
}

func (p *PostStore) Posts(ctx context.Context, offset, limit int) ([]*models.Post, error) {
	rows, err := p.db.QueryContext(ctx, "SELECT * FROM posts LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, errors.Join(ports.ErrInternal, err)
	}

	var posts []*models.Post

	for rows.Next() {
		post := new(models.Post)

		err := rows.Scan(&post.ID, &post.AuthorID, &post.Text, &post.Published, &post.DisableComments)
		if err != nil {
			return nil, errors.Join(ports.ErrInternal, err)
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (p *PostStore) Post(ctx context.Context, postID models.ID) (*models.Post, error) {
	row := p.db.QueryRowContext(ctx, "SELECT * FROM posts WHERE id = $1", postID)

	var post models.Post

	err := row.Scan(&post.ID, &post.AuthorID, &post.Text, &post.Published, &post.DisableComments)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(ports.ErrPostNotFound, err)
		}

		return nil, errors.Join(ports.ErrInternal, err)
	}

	return &post, nil
}

func (p *PostStore) AddPost(ctx context.Context, post *models.Post) error {
	post.ID = models.NewID()

	_, err := p.db.ExecContext(ctx, "INSERT INTO posts (id, author_id, post_text, published, disable_comments) VALUES ($1, $2, $3, $4, $5)", post.ID, post.AuthorID, post.Text, post.Published, post.DisableComments)
	if err != nil {
		return errors.Join(ports.ErrInternal, err)
	}

	return nil
}

var _ ports.PostStore = (*PostStore)(nil)

func NewPostStore(db *sql.DB) *PostStore {
	return &PostStore{
		db: db,
	}
}
