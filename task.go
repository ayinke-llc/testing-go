package testinggo

import (
	"time"

	"github.com/google/uuid"
)

type TaskItem struct {
	ID          uuid.UUID
	Title       string
	Description string

	CreatedAt time.Time
	UpdatedAt time.Time
}
