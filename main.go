package main

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"strings"
	"sync"

	src "mangia_nastri/src"
)

var logger = src.SetupLogger()

type proxyHandler struct {
	mu sync.Mutex // guards n
	n  int
}

// `ComputeRequestHash` generates a SHA-256 hash for an HTTP request.
// It processes the request headers and body, and combines them with the request
// method and URL to create a unique content string. This content string is then
// hashed using SHA-256. The function logs the request details and the first ten
// characters of the hash for debugging purposes.
//
// Parameters
//   - r: the HTTP request to be hashed.
//
// Returns
//   A string representing the SHA-256 hash of the request content.

func (h *proxyHandler) ComputeRequestHash(r *http.Request) string {
	headers := src.ProcessHeaders(r.Header)
	body := src.ProcessBody(r.Body)
	url := r.URL.String()

	content := strings.Join([]string{r.Method, url, headers, body}, ", ")
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(content)))

	logger.Info("Request", "hash", hash[:10], "method", r.Method, "url", url, "headers", headers, "body", body)
	return hash
}

// `ServeHTTP` is the main entry point for the `proxyHandler` type. It is called
// when an HTTP request is received by the server. The function increments the
// request counter and calls `ComputeRequestHash` to generate a hash for the request.
//
// Parameters
//   - w: the HTTP response writer;
//   - r: the HTTP request to be processed.

func (h *proxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.n++
	h.ComputeRequestHash(r)
}

func main() {
	http.Handle("/", new(proxyHandler))

	logger.Print("\nMangia_nastri is ready to rock.\n")
	logger.Fatal(http.ListenAndServe(":8080", nil))
}
