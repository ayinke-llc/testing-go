package memory

import (
	testinggo "ayinke-llc/gophercrunch/testing-go"
	"context"
	"sync"

	"github.com/google/uuid"
)

type MemoryStore struct {
	rw sync.RWMutex

	items map[uuid.UUID]*testinggo.TaskItem
}

func New() *MemoryStore {
	return &MemoryStore{
		rw:    sync.RWMutex{},
		items: make(map[uuid.UUID]*testinggo.TaskItem),
	}
}

func (m *MemoryStore) Close() error {

	m.rw.Lock()
	defer m.rw.Unlock()

	m.items = make(map[uuid.UUID]*testinggo.TaskItem)
	return nil
}

func (m *MemoryStore) Create(_ context.Context,
	item *testinggo.TaskItem) error {

	m.rw.Lock()
	defer m.rw.Unlock()

	m.items[item.ID] = item
	return nil
}

func (m *MemoryStore) Delete(_ context.Context,
	id uuid.UUID) error {

	m.rw.Lock()
	defer m.rw.Unlock()

	delete(m.items, id)

	return nil
}

func (m *MemoryStore) Get(_ context.Context, id uuid.UUID) (*testinggo.TaskItem, error) {

	m.rw.RLock()
	defer m.rw.RUnlock()

	item, ok := m.items[id]
	if !ok {
		return nil, testinggo.ErrItemNotFound
	}

	return item, nil
}
