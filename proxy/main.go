package proxy

import (
	"mangia_nastri/conf"
	"mangia_nastri/logger"
	"net/http"
	"sync"
)

var log = logger.New("proxy")

// ServeHTTP is the main entry point for the `proxyHandler` type. It is called
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

	h.computeRequestHash(r)
}

type proxyHandler struct {
	mu     sync.Mutex // guards n
	n      int
	config *conf.Config
}

func New(config *conf.Config) (proxy *proxyHandler) {
	proxy = &proxyHandler{}

	proxy.config = config

	return
}
