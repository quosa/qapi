package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/quosa/qapi/internal/bugs"
)

func TestV1GetAllBugsEmpty(t *testing.T) {

	bugStorage := bugs.NewBugService()
	v1handler := &v1APIHandler{bugStorage: bugStorage}
	mux := ConfigureV1APIRoutes(v1handler)
	rec := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/v1/bugs", nil)
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

func TestV1AddBug(t *testing.T) {
	bugStorage := bugs.NewBugService()
	v1handler := &v1APIHandler{bugStorage: bugStorage}
	mux := ConfigureV1APIRoutes(v1handler)
	rec := httptest.NewRecorder()

	req, err := http.NewRequest("POST", "/v1/bugs", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", rec.Code)
	}
	// if rec.Body.String() != "{New bug}" {
	// 	t.Errorf("expected new bug, got %s", rec.Body.String())
	// }
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
			name:       "get all bugs gives 200",
			path:       "/v1/bugs",
			wantStatus: http.StatusOK,
			wantBody:   "[{New bug}]", // 1 item initialized in the list
		},
		{
			name:       "get only bug gives 200",
			path:       "/v1/bugs/1", // ID counter starts at 1
			wantStatus: http.StatusOK,
			wantBody:   "{New bug}", // 1 item initialized in the list
		},
		{
			name:       "invalid bug id gives 400",
			path:       "/v1/bugs/ff3%293",
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "bug id above bounds gives 404",
			// ID counter starts at 1 and 1 item initialized in the list
			path:       "/v1/bugs/2",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "bug id too big for int gives 400",
			path:       "/v1/bugs/12345679012345678901234567890",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "bug id below bounds gives 404",
			path:       "/v1/bugs/-1",
			wantStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bugStorage := bugs.NewBugService()
			_, err := bugStorage.CreateBug(bugs.Bug{Title: "New bug"})
			if err != nil {
				t.Fatalf("could not create bug: %v", err)
			}
			v1handler := &v1APIHandler{bugStorage: bugStorage}

			mux := ConfigureV1APIRoutes(v1handler)
			rec := httptest.NewRecorder()

			req, err := http.NewRequest("GET", tt.path, nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			mux.ServeHTTP(rec, req)
			if rec.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, rec.Code)
			}
			// if tt.wantBody != "" && rec.Body.String() != tt.wantBody {
			// 	t.Errorf("expected '%s', got %s", tt.wantBody, rec.Body.String())
			// }
		})
	}
}
