package testinggo

import (
	"context"
	"io"

	"github.com/google/uuid"
)

type Store interface {
	io.Closer
	Create(context.Context, *TaskItem) error
	Delete(context.Context, uuid.UUID) error
	Get(context.Context, uuid.UUID) (*TaskItem, error)
}
