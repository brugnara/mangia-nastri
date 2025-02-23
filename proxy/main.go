package proxy

import (
	"bytes"
	"io"
	"mangia_nastri/commander"
	"mangia_nastri/conf"
	"mangia_nastri/datasources"
	"mangia_nastri/datasources/inMemory"
	"mangia_nastri/datasources/redis"
	"mangia_nastri/datasources/sqlite3"
	"mangia_nastri/logger"
	"net/http"
	"sync"
)

type proxyHandler struct {
	mu          sync.Mutex // protects handledReqs
	handledReqs int
	config      *conf.Proxy
	dataSource  datasources.DataSource
	log         logger.Logger
	client      *http.Client

	Action chan commander.Action
	Ready  chan bool
	Name   string
}

func (p *proxyHandler) proxy(r *http.Request) (payload datasources.Payload, err error) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Restore body

	proxyReq, err := http.NewRequest(r.Method, p.config.Destination+r.URL.Path, r.Body)

	if err != nil {
		return
	}

	// Create the payload
	payload = datasources.Payload{
		Request: datasources.Request{
			Header: make(http.Header),
			URL:    r.URL.String(),
			Body:   string(bodyBytes),
		},
		Response: datasources.Response{
			Header: make(http.Header),
		},
	}

	// Copy headers
	for name, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(name, value)
			// some logics to ignore headers edventually goes here
			payload.Request.Header.Add(name, value)
		}
	}

	// Perform the request
	resp, err := p.client.Do(proxyReq)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for name, values := range resp.Header {
		for _, value := range values {
			payload.Response.Header.Add(name, value)
		}
	}

	// Write the status code
	payload.Response.Status = resp.StatusCode

	// Copy response body to a tmp
	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	// store response data
	payload.Response.Body = string(bodyBytes)

	return
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

	p.handledReqs++
	hash := p.computeRequestHash(r)
	log := p.log.CloneWithPrefix(hash.String())

	payload, err := p.dataSource.Get(hash)

	if err == nil {
		log.Info("Request got from cache", "payload", payload)
	} else {
		log.Info("Processing request", "hash", hash)
		// do the request to the destination, store in value the result
		payload, err = p.proxy(r)
		if err != nil {
			log.Error("Error processing request", "error", err)
			http.Error(w, "Error processing request", http.StatusInternalServerError)
			return
		}

		if p.config.DoRecord {
			log.Info("Recording request")
			err = p.dataSource.Set(hash, payload)
			if err != nil {
				log.Error("Error setting value", "error", err)
			}
		} else {
			log.Info("Not recording request")
		}
	}

	// Write the response
	for k, v := range payload.Response.Header {
		w.Header().Set(k, v[0])
	}

	w.WriteHeader(payload.Response.Status)

	_, err = w.Write([]byte(payload.Response.Body))

	if err != nil {
		log.Error("Error writing response", "error", err)
		http.Error(w, "Error writing response", http.StatusInternalServerError)
	}

	p.log.Info("Request processed", "hash", hash)
}

func New(config *conf.Proxy, log logger.Logger) (proxy *proxyHandler) {
	proxy = &proxyHandler{
		log:    log.CloneWithPrefix("proxy:" + config.Name),
		config: config,
		Name:   config.Name,
		client: &http.Client{},
		Action: make(chan commander.Action),
		Ready:  make(chan bool),
	}

	proxy.log.Info("Creating proxy", "name", config.Name, "destination", config.Destination)

	go func() {
		for a := range proxy.Action {
			switch a {
			case commander.DO_RECORD:
				proxy.log.Info("Setting change: recording requests")
				proxy.config.DoRecord = true
			case commander.DO_NOT_RECORD:
				proxy.log.Info("Setting change: not recording requests")
				proxy.config.DoRecord = false
			}
		}
	}()

	switch config.DataSource.Type {
	case "inMemory":
		proxy.dataSource = inMemory.New(&proxy.log)
	case "redis":
		proxy.dataSource = redis.New(&proxy.log, config.DataSource.URI)
	case "sqlite3":
		proxy.dataSource = sqlite3.New(&proxy.log, config.DataSource.URI)
	default:
		proxy.log.Fatalf("Unknown data source: %v", config.DataSource)

		return
	}

	go func() {
		log.Info("Proxy is waiting for dataSource to be ready")
		<-proxy.dataSource.Ready()

		log.Infof("Proxy is ready to serve requests to %s on localhost:%d. record=%t", config.Destination, config.Port, config.DoRecord)
		proxy.Ready <- true
	}()

	return
}
