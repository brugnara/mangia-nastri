package proxy

import (
	"encoding/json"
	"io"
	"slices"
)

// ProcessHeaders reads the raw headers of an HTTP request and processes them
// into a single string, sorting all key-value pairs alphabetically.
//
// Parameters
//   - rawHeaders: a map containing the raw headers of an HTTP request.
//
// Returns
//
//	A string containing the processed and sorted headers.
func (p *proxyHandler) ProcessHeaders(rawHeaders map[string][]string, ignore []string) string {
	// this may lead to a memory leak, check
	headersMap := make(map[string]interface{}, len(rawHeaders))

	for k, v := range rawHeaders {
		if slices.Contains(ignore, k) {
			log.Info("Ignoring header", "header", k)
			continue
		}
		headersMap[k] = v
	}

	return p.stringifyObject(
		(headersMap),
	)
}

// ProcessBody reads the raw body of an HTTP request and processes it into a
// single string, sorting all key-value pairs alphabetically.
//
// Parameters
//   - rawBody: an object containing the raw body of an HTTP request.
//
// Returns
//
//	A string containing the processed and sorted body.
func (p *proxyHandler) ProcessBody(rawBody io.ReadCloser) string {
	defer rawBody.Close()

	var body map[string]interface{}
	if err := json.NewDecoder(rawBody).Decode(&body); err != nil {
		log.Error("Failed to decode body", "error", err)

		return ""
	}

	return p.stringifyObject(
		(body),
	)
}

// stringifyObject converts a generic object to its JSON string representation.
// If the marshalling process fails, it logs the error and terminates the program.
//
// Parameters
//   - obj: the generic object to be converted to a JSON string.
//
// Returns
//
//	A string containing the JSON representation of the input object.
func (p *proxyHandler) stringifyObject(obj map[string]interface{}) string {
	bytes, err := json.Marshal(obj)
	if err != nil {
		log.Error("Failed to marshal sorted body", "error", err)

		return ""
	}

	return string(bytes)
}
