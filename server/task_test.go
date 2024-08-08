package server

import (
	testinggo "ayinke-llc/gophercrunch/testing-go"
	"ayinke-llc/gophercrunch/testing-go/mocks"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	gomock "github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestFetchTaskItem(t *testing.T) {

	tt := []struct {
		mockFn             func(store *mocks.MockStore)
		name               string
		expectedStatusCode int
		id                 uuid.UUID
	}{
		{
			name:               "id from url is invalid",
			expectedStatusCode: http.StatusBadRequest,
			id:                 uuid.Nil,
			mockFn: func(store *mocks.MockStore) {

			},
		},
		{
			name:               "item not found",
			expectedStatusCode: http.StatusNotFound,
			id:                 uuid.New(),
			mockFn: func(store *mocks.MockStore) {
				store.EXPECT().Get(gomock.Any(), gomock.Any()).Times(1).
					Return(nil, testinggo.ErrItemNotFound)
			},
		},
		{
			name:               "unknown error when fetching an item",
			expectedStatusCode: http.StatusInternalServerError,
			id:                 uuid.New(),
			mockFn: func(store *mocks.MockStore) {
				store.EXPECT().Get(gomock.Any(), gomock.Any()).Times(1).
					Return(nil, errors.New("error not found"))
			},
		},
		{
			name:               "item fetched successfully",
			expectedStatusCode: http.StatusOK,
			id:                 uuid.New(),
			mockFn: func(store *mocks.MockStore) {
				store.EXPECT().
					Get(gomock.Any(), gomock.Any()).Times(1).
					Return(&testinggo.TaskItem{
						ID: uuid.New(),
					}, nil)
			},
		},
	}

	for _, v := range tt {

		t.Run(v.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mocks.NewMockStore(ctrl)

			v.mockFn(store)

			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))

			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("id", v.id.String())

			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

			fetchTaskItem(store).ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
		})
	}
}

func TestCreateTaskItem(t *testing.T) {

	tt := []struct {
		req                *createTaskItemRequest
		shouldFail         bool
		mockFn             func(store *mocks.MockStore)
		name               string
		expectedStatusCode int
	}{
		{
			name: "title not present",
			mockFn: func(store *mocks.MockStore) {
				store.EXPECT().Create(gomock.Any(), gomock.Any()).
					Times(0).Return(nil)
			},
			expectedStatusCode: http.StatusBadRequest,
			req: &createTaskItemRequest{
				Description: "this is my Description",
			},
			shouldFail: true,
		},
		{
			name: "description not present",
			mockFn: func(store *mocks.MockStore) {
				store.EXPECT().Create(gomock.Any(), gomock.Any()).
					Times(0).Return(nil)
			},
			expectedStatusCode: http.StatusBadRequest,
			req: &createTaskItemRequest{
				Title: "this is my title",
			},
			shouldFail: true,
		},
		{
			name: "item creation fails",
			mockFn: func(store *mocks.MockStore) {
				store.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("could not create an item"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: &createTaskItemRequest{
				Title:       "this is my title",
				Description: "this is my description",
			},
			shouldFail: true,
		},
		{
			name: "item creation succeeds",
			mockFn: func(store *mocks.MockStore) {
				store.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			req: &createTaskItemRequest{
				Title:       "this is my title",
				Description: "this is my description",
			},
		},
	}

	for _, v := range tt {

		t.Run(v.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mocks.NewMockStore(ctrl)

			v.mockFn(store)

			rr := httptest.NewRecorder()

			var b = bytes.NewBuffer(nil)

			err := json.NewEncoder(b).Encode(v.req)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/", b)

			createTaskItem(store).ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
		})
	}
}
