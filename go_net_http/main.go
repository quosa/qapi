package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/quosa/qapi/internal/bugs"
)

type v1APIHandler struct {
	bugStorage bugs.IBugService
}

func (h *v1APIHandler) getAllBugs(w http.ResponseWriter, r *http.Request) {
	allBugs, err := h.bugStorage.GetAllBugs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%+v", allBugs)
}

func (h *v1APIHandler) addBug(w http.ResponseWriter, r *http.Request) {
	bug, err := h.bugStorage.CreateBug(bugs.Bug{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Status:      r.FormValue("status"),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%+v", bug)
}

func (h *v1APIHandler) getSingleBug(w http.ResponseWriter, r *http.Request) {
	i, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid task ID")
		return
	}
	bug, err := h.bugStorage.GetBugByID(uint64(i))
	// TODO: handle not found error properly
	if err != nil {
		if errors.Is(err, bugs.ErrorNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%+v", bug)
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
