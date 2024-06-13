package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

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
	if strings.TrimSpace(rec.Body.String()) != "[]" {
		t.Errorf("expected empty list, got %s", rec.Body.String())
	}
}

func TestV1AddBug(t *testing.T) {
	// Arrange
	bugStorage := bugs.NewBugService()
	v1handler := &v1APIHandler{bugStorage: bugStorage}
	mux := ConfigureV1APIRoutes(v1handler)
	rec := httptest.NewRecorder()

	new_bug := bugs.Bug{
		Title:       "New bug",
		Description: "New description",
	}
	body_bytes, _ := json.Marshal(new_bug)
	req, err := http.NewRequest("POST", "/v1/bugs", bytes.NewReader(body_bytes))
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	// Act
	mux.ServeHTTP(rec, req)

	// Assert
	if rec.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", rec.Code)
	}
	created_bug := bugs.Bug{}
	err = json.NewDecoder(rec.Body).Decode(&created_bug)
	if err != nil {
		t.Errorf("could not decode response: %v", err)
	}
	if created_bug.Title != new_bug.Title {
		t.Errorf("expected title %s, got %s", new_bug.Title, created_bug.Title)
	}
	if created_bug.Description != new_bug.Description {
		t.Errorf("expected description %s, got %s", new_bug.Description, created_bug.Description)
	}
}

func TestV1API(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantBody   interface{}
	}{
		{
			name:       "root gives 200",
			path:       "/",
			wantStatus: http.StatusOK,
			wantBody:   "Welcome to QAPI Bugs!",
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
			wantBody: []bugs.Bug{
				{ID: 1, Title: "New bug"},
			}, // 1 item initialized in the list
		},
		{
			name:       "get only bug gives 200",
			path:       "/v1/bugs/1", // ID counter starts at 1
			wantStatus: http.StatusOK,
			wantBody:   bugs.Bug{ID: 1, Title: "New bug"}, // 1 item initialized in the list
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
			wantStatus: http.StatusBadRequest,
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

			// This became overly complicated due to go reflection and static typing
			// Consider chopping this test back up per API path and verb...
			if tt.wantBody != nil {
				var gotBody interface{}
				switch reflect.TypeOf(tt.wantBody) {
				case reflect.TypeOf([]bugs.Bug{}):
					bugList := []bugs.Bug{}
					err = json.NewDecoder(rec.Body).Decode(&bugList)
					gotBody = bugList
				case reflect.TypeOf(bugs.Bug{}):
					bug := bugs.Bug{}
					err = json.NewDecoder(rec.Body).Decode(&bug)
					gotBody = bug
				case reflect.TypeOf(""):
					gotBody = rec.Body.String()
				default:
					t.Errorf("unexpected want body type %v", reflect.TypeOf(tt.wantBody))
				}
				if err != nil {
					t.Errorf("could not decode response: %v", err)
				}
				opts := cmpopts.IgnoreFields(bugs.Bug{}, "CreatedAt", "UpdatedAt")
				if diff := cmp.Diff(tt.wantBody, gotBody, opts); diff != "" {
					t.Errorf("case %s mismatch (-want +got):\n%s", tt.name, diff)
				}
			}
		})
	}
}
