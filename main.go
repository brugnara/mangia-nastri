package main

import (
	"fmt"
	"sync"
	"strings"
	"net/http"
	"encoding/json"
	"crypto/sha256"

	src "mangia_nastri/src"
)

var Logger = src.SetupLogger()

type proxyHandler struct {
	mu sync.Mutex // guards n
	n  int
}

func (h *proxyHandler) ComputeRequestHash(r *http.Request) string {
	var headers string = src.ProcessHeaders(r.Header)
	var body map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		Logger.Error("Failed to parse request body", "error", err)
	} else {
		Logger.Debug("Parsed request body", "body", body)
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		Logger.Error("Failed to encode body to JSON", "error", err)
	} else {
		Logger.Debug("Encoded body to JSON", "body", string(bodyBytes))
	}

	url := r.URL.String()

	content := strings.Join([]string{r.Method, url, headers, string(bodyBytes)}, ", ")

	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(content)))

	Logger.Info("Request", "method", r.Method, "url", url, "headers", headers, "body", string(bodyBytes), "hash", hash[:10])

	// return "ciao lorenzino"
	return hash
}

func (h *proxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.n++
	hash := h.ComputeRequestHash(r)
	Logger.Print(hash)
}

func main() {
	http.Handle("/", new(proxyHandler))
	Logger.Print("\nMangia-nastri is ready to rock.\n")
	Logger.Fatal(http.ListenAndServe(":8080", nil))
}
