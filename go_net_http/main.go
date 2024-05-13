package main

import (
	"fmt"
	"net/http"
)

type v1APIHandler struct{}

func (h *v1APIHandler) getAllTasks(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "[]")
}

func main() {
	var mux = http.NewServeMux()
	var v1 = &v1APIHandler{}

	mux.HandleFunc("GET /v1/tasks/", v1.getAllTasks)
	// NOTE: {$} to match only root (others get 404)
	mux.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Welcome to QAPI!")
	})

	fmt.Println("QAPI Server is starting...")

	http.ListenAndServe(":8080", mux)
}

/*

// Welcome message only in root / path
curl -i http://localhost:8080/
HTTP/1.1 200 OK
Date: Mon, 13 May 2024 17:58:23 GMT
Content-Length: 16
Content-Type: text/plain; charset=utf-8

Welcome to QAPI!

// Empty array in /v1/tasks/ path (wrong content type still)
curl -i http://localhost:8080/v1/tasks/
HTTP/1.1 200 OK
Date: Mon, 13 May 2024 17:58:37 GMT
Content-Length: 2
Content-Type: text/plain; charset=utf-8

[]

// 404 for non-existing v2 path
curl -i http://localhost:8080/v2/tasks
HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Mon, 13 May 2024 17:58:53 GMT
Content-Length: 19

404 page not found

// 404 for nonsensical paths
curl -i http://localhost:8080/foo
HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Mon, 13 May 2024 17:59:00 GMT
Content-Length: 19

404 page not found
*/
