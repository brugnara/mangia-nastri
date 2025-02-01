package proxy

import (
	"mangia_nastri/conf"
	"mangia_nastri/datasources"
	"mangia_nastri/datasources/inMemory"
	"mangia_nastri/logger"
	"net/http"
	"sync"
)

type proxyHandler struct {
	mu         sync.Mutex // guards n
	n          int
	config     *conf.Proxy
	dataSource datasources.DataSource
	log        logger.Logger
}

// ServeHTTP is the main entry point for the `proxyHandler` type. It is called
// when an HTTP request is received by the server. The function increments the
// request counter and calls `ComputeRequestHash` to generate a hash for the request.
//
// Parameters
//   - w: the HTTP response writer;
//   - r: the HTTP request to be processed.
func (p *proxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.n++
	hash := p.computeRequestHash(r)

	value, err := p.dataSource.Get(hash)

	if err == nil {
		p.log.Info("Request already processed", "hash", hash, "value", value)
		return
	}

	err = p.dataSource.Set(hash, "ciao")
	if err != nil {
		p.log.Error("Error setting value", "hash", hash, "error", err)
	}
	p.log.Info("Request processed", "hash", hash)

}

func New(config *conf.Proxy, log logger.Logger) (proxy *proxyHandler) {
	proxy = &proxyHandler{}

	proxy.log = log.CloneWithPrefix("proxy:" + config.Name)

	proxy.config = config

	switch config.DataSource.Type {
	case "inMemory":
		proxy.dataSource = inMemory.New(&proxy.log)
	default:
		proxy.log.Fatalf("Unknown data source: %v", config.DataSource)
	}

	return
}
