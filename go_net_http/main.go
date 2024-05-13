package main

import (
	"fmt"
	"net/http"
)

func main() {
	var mux = http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Welcome to QAPI!")
	})

	fmt.Println("QAPI Server is starting...")

	http.ListenAndServe(":8080", mux)
}

/*
curl -i http://localhost:8080/
HTTP/1.1 200 OK
Date: Mon, 13 May 2024 17:25:21 GMT
Content-Length: 16
Content-Type: text/plain; charset=utf-8

Welcome to QAPI!
*/
