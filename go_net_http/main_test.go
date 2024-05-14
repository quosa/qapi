package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestV1GetAllTasksEmpty(t *testing.T) {

	var tasks = []Task{}
	mux := GetRoutes(tasks)
	rec := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/v1/tasks", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
	if rec.Body.String() != "[]" {
		t.Errorf("expected empty list, got %s", rec.Body.String())
	}
}

func TestV1AddTask(t *testing.T) {
	var tasks = []Task{}
	mux := GetRoutes(tasks)
	rec := httptest.NewRecorder()

	req, err := http.NewRequest("POST", "/v1/tasks", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
	if rec.Body.String() != "{New task}" {
		t.Errorf("expected new task, got %s", rec.Body.String())
	}
}

func TestV1API(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "root gives 200",
			path:       "/",
			wantStatus: http.StatusOK,
			wantBody:   "Welcome to QAPI!",
		},
		{
			name:       "non-existing route gives 404",
			path:       "/foo",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "v1 api root gives 404",
			path:       "/v1",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "get all tasks gives 200",
			path:       "/v1/tasks",
			wantStatus: http.StatusOK,
			wantBody:   "[{New task}]", // 1 item initialized in the list
		},
		{
			name:       "get only task gives 200",
			path:       "/v1/tasks/0",
			wantStatus: http.StatusOK,
			wantBody:   "{New task}", // 1 item initialized in the list
		},
		{
			name:       "invalid task id gives 400",
			path:       "/v1/tasks/ff3%293",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "task id above bounds gives 404",
			path:       "/v1/tasks/1", // 1 item initialized in the list
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "task id too big for int gives 400",
			path:       "/v1/tasks/12345679012345678901234567890",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "task id below bounds gives 404",
			path:       "/v1/tasks/-1",
			wantStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tasks = []Task{{Description: "New task"}}
			mux := GetRoutes(tasks)
			rec := httptest.NewRecorder()

			req, err := http.NewRequest("GET", tt.path, nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			mux.ServeHTTP(rec, req)
			if rec.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, rec.Code)
			}
			if tt.wantBody != "" && rec.Body.String() != tt.wantBody {
				t.Errorf("expected '%s', got %s", tt.wantBody, rec.Body.String())
			}
		})
	}
}
