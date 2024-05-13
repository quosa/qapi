package main

import (
	"fmt"
	"net/http"
	"strconv"
)

type Task struct {
	Description string
}
type v1APIHandler struct {
	tasks []Task
}

func (h *v1APIHandler) getAllTasks(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", h.tasks)
}

func (h *v1APIHandler) addTask(w http.ResponseWriter, r *http.Request) {
	h.tasks = append(h.tasks, Task{Description: "New task"})
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", h.tasks[len(h.tasks)-1])
}

func (h *v1APIHandler) getSingleTask(w http.ResponseWriter, r *http.Request) {
	i, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid task ID")
		return
	}
	if i < 0 || i >= len(h.tasks) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Task not found")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", h.tasks[i])
}
func GetRoutes(tasks []Task) *http.ServeMux {
	var mux = http.NewServeMux()
	var v1 = &v1APIHandler{tasks: tasks}

	// v1 API route config
	mux.HandleFunc("GET /v1/tasks", v1.getAllTasks)
	mux.HandleFunc("POST /v1/tasks", v1.addTask)
	mux.HandleFunc("GET /v1/tasks/{id}", v1.getSingleTask)

	// NOTE: {$} to match only root (others get 404)
	mux.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Welcome to QAPI!")
	})
	return mux
}
func main() {
	var tasks = []Task{}
	mux := GetRoutes(tasks)
	fmt.Println("QAPI Server is starting...")

	http.ListenAndServe(":8080", mux)
}

/*
See unit tests in main_test.go for v1 tests.

curl -i http://localhost:8080/
HTTP/1.1 200 OK
Date: Mon, 13 May 2024 18:12:36 GMT
Content-Length: 16
Content-Type: text/plain; charset=utf-8

Welcome to QAPI!

curl -i http://localhost:8080/gibberish
HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Mon, 13 May 2024 18:12:47 GMT
Content-Length: 19

404 page not found
*/
