package postgres

import (
	"context"
	"database/sql"
	"errors"

	ports "post-and-comments/internal/inner-ports"
	"post-and-comments/internal/models"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (s *UserStore) User(ctx context.Context, id models.ID) (*models.User, error) {
	row := s.db.QueryRowContext(ctx, "SELECT * FROM users WHERE id = $1", id)

	var user models.User

	err := row.Scan(&user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrUserNotFound
		}

		return nil, errors.Join(ports.ErrInternal, err)
	}

	return &user, nil
}

func (s *UserStore) AddUser(ctx context.Context, user *models.User) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO users (id) VALUES ($1)", user.ID)
	if err != nil {
		return errors.Join(ports.ErrInternal, err)
	}

	return nil
}
