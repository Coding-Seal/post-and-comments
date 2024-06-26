package models

import (
	"errors"
	"fmt"
	"time"
)

const MaxSymbols = 2000

var (
	ErrTooMuchSymbols = fmt.Errorf("too much symbols, limit %d", MaxSymbols)
	ErrCommentEmpty   = errors.New("empty comment")
)

type Comment struct {
	ID        ID
	AuthorID  ID
	Published time.Time
	Text      CommentText
	PostID    ID
	Replies   []*Comment
}

type CommentText string

func NewCommentText(text string) (CommentText, error) {
	if len(text) == 0 {
		return "", ErrCommentEmpty
	}

	if len(text) > MaxSymbols { // Maybe rewrite to account for UTF
		return "", ErrTooMuchSymbols
	}

	return CommentText(text), nil
}
