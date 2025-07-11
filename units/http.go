package main

import (
	"fmt"
	"net/http"
)

// Simple counter server.
type Counter struct {
	n int
}

func (ctr *Counter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctr.n++
	fmt.Fprintf(w, "counter = %d\n", ctr.n)
}

func main() {
	ctr := new(Counter)
	http.Handle("/counter", ctr)
}
