package proxy

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"mangia_nastri/src"
	"net/http"
	"strings"
)

// computeRequestHash generates a SHA-256 hash for an HTTP request.
// It processes the request headers and body, and combines them with the request
// method and URL to create a unique content string. This content string is then
// hashed using SHA-256. The function logs the request details and the first ten
// characters of the hash for debugging purposes.
//
// Parameters
//   - r: the HTTP request to be hashed.
//
// Returns
//
//	A string representing the SHA-256 hash of the request content.
func (h *proxyHandler) computeRequestHash(r *http.Request) Hash {
	headers := src.ProcessHeaders(r.Header, h.config.Ignore.Headers)
	body := src.ProcessBody(r.Body)
	url := r.URL.String()

	content := strings.Join([]string{r.Method, url, headers, body}, ", ")
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(content)))

	log.Info("Request", "hash", hash[:10], "method", r.Method, "url", url, "headers", headers, "body", body)

	return h.hash(hash)
}

func (h *proxyHandler) hash(doc string) Hash {
	// Dumb af, but it's a cheap way to specific the most generic thing
	// you can :-/
	var v interface{}
	err := json.Unmarshal([]byte(doc), &v) // NB: You should handle errors :-/

	if err != nil {
		log.Error("Failed to marshal sorted body", "error", err)
		return Hash("")
	}

	cdoc, _ := json.Marshal(v)
	sum := sha256.Sum256(cdoc)
	return Hash(hex.EncodeToString(sum[0:]))
}
