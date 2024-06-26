package models

import (
	"github.com/google/uuid"
)

type ID uuid.UUID

func ParseID(id string) (ID, error) {
	uid, err := uuid.Parse(id)
	return ID(uid), err
}

func NewID() ID {
	return ID(uuid.New())
}

func (id ID) String() string {
	return uuid.UUID(id).String()
}
