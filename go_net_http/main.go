package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/quosa/qapi/internal/bugs"
)

type v1APIHandler struct {
	bugStorage bugs.IBugService
}

// Utility function to read JSON request body (safely)
// and unmarshal it into target struct
func readJSONBody(r *http.Request, target interface{}) error {
	// TODO: check content-encoding
	// Body can be empty
	if r.Body == nil {
		return errors.New("empty body")
	}
	// Body can be too long
	// Body can be malformed JSON
	// Body can be valid JSON but not the expected structure
	err := json.NewDecoder(r.Body).Decode(&target)
	if err != nil {
		fmt.Printf("JSON body parser, ERROR: %+v\n", err)
		return err
	}
	return nil
}

// Utility function to return JSON response with HTTP status code
func returnJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func (h *v1APIHandler) getAllBugs(w http.ResponseWriter, r *http.Request) {
	allBugs, err := h.bugStorage.GetAllBugs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	returnJSON(w, http.StatusOK, allBugs)
}

func (h *v1APIHandler) addBug(w http.ResponseWriter, r *http.Request) {
	bug := bugs.Bug{}
	err := readJSONBody(r, &bug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// TODO: validate incoming bug data

	created_bug, err := h.bugStorage.CreateBug(bug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	returnJSON(w, http.StatusCreated, created_bug)
}

func (h *v1APIHandler) getSingleBug(w http.ResponseWriter, r *http.Request) {
	i, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid task ID")
		return
	}
	if i < 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bug, err := h.bugStorage.GetBugByID(uint64(i))
	if err != nil {
		if errors.Is(err, bugs.ErrorInvalidInput) {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else if errors.Is(err, bugs.ErrorNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	returnJSON(w, http.StatusOK, bug)
}

func ConfigureV1APIRoutes(v1 *v1APIHandler) *http.ServeMux {
	var mux = http.NewServeMux()

	// v1 API route config
	mux.HandleFunc("GET /v1/bugs", v1.getAllBugs)
	mux.HandleFunc("POST /v1/bugs", v1.addBug)
	mux.HandleFunc("GET /v1/bugs/{id}", v1.getSingleBug)

	// NOTE: {$} to match only root (others get 404)
	mux.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Welcome to QAPI Bugs!")
	})
	return mux
}
func main() {
	bugStorage := bugs.NewBugService()
	var v1APIHandler = &v1APIHandler{bugStorage: bugStorage}
	mux := ConfigureV1APIRoutes(v1APIHandler)
	fmt.Println("QAPI bug server is starting...")

	http.ListenAndServe(":8080", mux)
}
