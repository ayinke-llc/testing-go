package server

import (
	testinggo "ayinke-llc/gophercrunch/testing-go"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type createTaskItemRequest struct {
	Title       string
	Description string
}

func createTaskItem(store testinggo.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req = new(createTaskItemRequest)

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"message" : "invalid request body"}`))
			return
		}

		if strings.TrimSpace(req.Title) == "" || strings.TrimSpace(req.Description) == "" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"message" : "please provide a valid title or Description"}`))
			return
		}

		item := &testinggo.TaskItem{
			Title:       req.Title,
			ID:          uuid.New(),
			Description: req.Description,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := store.Create(r.Context(), item); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"message" : "could not create task item"}`))
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"id" :"%s", "message" : "task created successfully"}`, item.ID)))
	}
}

func fetchTaskItem(store testinggo.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		idFromRequest := chi.URLParam(r, "id")

		id, err := uuid.Parse(idFromRequest)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"message" : "invalid task item ID"}`))
			return
		}

		if id == uuid.Nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"message" : "uuid cannot be an empty one"}`))
			return
		}

		item, err := store.Get(r.Context(), id)
		if errors.Is(err, testinggo.ErrItemNotFound) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"message" : "item not found"}`))
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"message" : "an error occurred"}`))
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(item)
	}
}
