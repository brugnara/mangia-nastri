package main

import (
	"os"
	"fmt"
	"sort"
	"sync"
	"strings"
	"net/http"
	"encoding/json"
	"crypto/sha256"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

type proxyHandler struct {
	mu sync.Mutex // guards n
	n  int
}

func setupLogger() log.Logger {
	styles := log.DefaultStyles()
	styles.Levels[log.InfoLevel] = lipgloss.NewStyle().
		SetString("MangiaNastri").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("201")).
		Foreground(lipgloss.Color("#FFFFFF"))

	styles.Keys["err"] = lipgloss.NewStyle().Foreground(lipgloss.Color("204"))
	styles.Values["err"] = lipgloss.NewStyle().Bold(true)

	logger := log.New(os.Stderr)
	logger.SetStyles(styles)

	return *logger
}

var logger = setupLogger()

// Actual logic

func (h *proxyHandler) ComputeRequestHash(r *http.Request) string {
	headerArray := make([]string, 0, len(r.Header))

	for k := range r.Header {
		headerArray = append(headerArray, k + ": " + r.Header[k][0])
	}

	sort.Strings(headerArray)
	headers := strings.Join(headerArray, ", ")

	var body map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		logger.Error("Failed to parse request body", "error", err)
	} else {
		logger.Debug("Parsed request body", "body", body)
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		logger.Error("Failed to encode body to JSON", "error", err)
	} else {
		logger.Debug("Encoded body to JSON", "body", string(bodyBytes))
	}

	url := r.URL.String()

	content := strings.Join([]string{r.Method, url, headers, string(bodyBytes)}, ", ")

	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(content)))

	logger.Info("Request", "method", r.Method, "url", url, "headers", headers, "body", string(bodyBytes), "hash", hash[:10])

	// return "ciao lorenzino"
	return hash
}

func (h *proxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.n++
	hash := h.ComputeRequestHash(r)
	logger.Print(hash)
	// logger.Printf("count is %d", h.n)
}

func main() {
	http.Handle("/", new(proxyHandler))
	logger.Print("\nMangia-nastri is ready to rock.\n")
	logger.Fatal(http.ListenAndServe(":8080", nil))
}
