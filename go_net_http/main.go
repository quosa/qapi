package main

import (
	"fmt"
	"net/http"
)

func main() {
	var qapi_server http.Server

	fmt.Println("QAPI Server is starting...")

	qapi_server.Addr = ":8080"
	qapi_server.ListenAndServe()
}

/*
curl -i http://localhost:8080/
HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Mon, 13 May 2024 17:16:04 GMT
Content-Length: 19

404 page not found
*/
