package models

import (
	"errors"
	"time"
)

var ErrPostEmpty = errors.New("empty post")

type Post struct {
	ID              ID
	AuthorID        ID
	Text            PostText
	Published       time.Time
	DisableComments bool
}

type PostText string

func NewPostText(text string) (PostText, error) {
	if len(text) == 0 {
		return "", ErrPostEmpty
	}

	return PostText(text), nil
}
