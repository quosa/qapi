package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test GET /v1/tasks
func TestV1GetAllTasks(t *testing.T) {

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

func TestV1GetSingleTask(t *testing.T) {
	var tasks = []Task{{Description: "New task"}}
	mux := GetRoutes(tasks)
	rec := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/v1/tasks/0", nil)
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

func TestV1GetAllTaskWithOneEntry(t *testing.T) {
	var tasks = []Task{{Description: "New task"}}
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
	if rec.Body.String() != "[{New task}]" {
		t.Errorf("expected new task, got %s", rec.Body.String())
	}
}

func TestV1GetSingleTaskOutOfBounds(t *testing.T) {
	var tasks = []Task{{Description: "New task"}}
	mux := GetRoutes(tasks)
	rec := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/v1/tasks/1", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rec.Code)
	}
	if rec.Body.String() != "Task not found" {
		t.Errorf("expected 'Task not found', got %s", rec.Body.String())
	}
}

func TestV1GetSingleTaskInvalidID(t *testing.T) {
	var tasks = []Task{{Description: "New task"}}
	mux := GetRoutes(tasks)
	rec := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/v1/tasks/foobar", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
	if rec.Body.String() != "Invalid task ID" {
		t.Errorf("expected 'Invalid task ID', got %s", rec.Body.String())
	}
}
