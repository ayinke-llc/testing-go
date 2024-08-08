package testinggo

import (
	"context"
	"errors"
	"io"

	"github.com/google/uuid"
)

var ErrItemNotFound = errors.New("task item not found")

type Store interface {
	io.Closer
	Create(context.Context, *TaskItem) error
	Delete(context.Context, uuid.UUID) error
	Get(context.Context, uuid.UUID) (*TaskItem, error)
}
