package main

import (
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

func verifyReq(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Request: %s!", r.URL.Path[2:])
}

func main() {

	http.HandleFunc("/verify", verifyReq)
	http.ListenAndServe(":4000", nil)
}
