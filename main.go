package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type proxyHandler struct {
	mu sync.Mutex // guards n
	n  int
}

func (h *proxyHandler) computeRequestHash(r *http.Request) string {
	fmt.Println(r.Header)
	fmt.Println(r.Method)

	return "ciao lorenzino"
}

func (h *proxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.n++
	fmt.Fprintf(w, "count is %d\n", h.n)

	fmt.Fprintf(w, r.RequestURI)

	fmt.Println(h.computeRequestHash(r))
}

func main() {
	http.Handle("/", new(proxyHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
