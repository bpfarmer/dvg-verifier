package main

import (
	"fmt"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
)

func verifyReq(w http.ResponseWriter, r *http.Request) {
	req := strings.Split(r.URL.Path[1:], "/")
	addr := req[1]
	val := req[2]
	if len(addr) > 0 && len(val) > 0 {

	}
	fmt.Fprintf(w, "Request: %s!", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/verify/", verifyReq)
	http.ListenAndServe(":4000", nil)
}
