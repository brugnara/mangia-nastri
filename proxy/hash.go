package proxy

import (
	"crypto/sha256"
	"fmt"
	"mangia_nastri/datasources"
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
func (p *proxyHandler) computeRequestHash(r *http.Request) datasources.Hash {
	headers := p.ProcessHeaders(r.Header, p.config.Ignore.Headers)
	body := p.ProcessBody(r.Body)
	url := r.URL.String()

	content := strings.Join([]string{r.Method, url, headers, body}, ", ")
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(content)))

	p.log.Info("Request", "hash", hash[:10], "method", r.Method, "url", url, "headers", headers, "body", body)

	return datasources.Hash(hash)
}
