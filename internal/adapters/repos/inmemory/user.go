package inmemory

import (
	"context"

	ports "post-and-comments/internal/inner-ports"
	"post-and-comments/internal/models"
)

type UserStore struct {
	table map[models.ID]*models.User
}

func NewUserStore() *UserStore {
	return &UserStore{
		table: make(map[models.ID]*models.User),
	}
}

func (s *UserStore) User(ctx context.Context, id models.ID) (*models.User, error) {
	if u, ok := s.table[id]; ok {
		return u, nil
	}

	return nil, ports.ErrUserNotFound
}

func (s *UserStore) AddUser(ctx context.Context, user *models.User) error {
	user.ID = models.NewID()
	s.table[user.ID] = user

	return nil
}
