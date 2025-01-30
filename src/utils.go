package src

import (
	"encoding/json"
	"io"
	"mangia_nastri/logger"
	"slices"
)

var log = logger.New("utils")

// stringifyObject converts a generic object to its JSON string representation.
// If the marshalling process fails, it logs the error and terminates the program.
//
// Parameters
//   - obj: the generic object to be converted to a JSON string.
//
// Returns
//
//	A string containing the JSON representation of the input object.
func stringifyObject(obj map[string]interface{}) string {
	bytes, err := json.Marshal(obj)
	if err != nil {
		log.Error("Failed to marshal sorted body", "error", err)

		return ""
	}

	return string(bytes)
}

// ProcessHeaders reads the raw headers of an HTTP request and processes them
// into a single string, sorting all key-value pairs alphabetically.
//
// Parameters
//   - rawHeaders: a map containing the raw headers of an HTTP request.
//
// Returns
//   A string containing the processed and sorted headers.

func ProcessHeaders(rawHeaders map[string][]string, ignore []string) string {
	headersMap := make(map[string]interface{}, len(rawHeaders))
	for k, v := range rawHeaders {
		if slices.Contains(ignore, k) {
			log.Info("Ignoring header", "header", k)
			continue
		}
		headersMap[k] = v
	}

	return stringifyObject(
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
func ProcessBody(rawBody io.ReadCloser) string {
	defer rawBody.Close()

	var body map[string]interface{}
	if err := json.NewDecoder(rawBody).Decode(&body); err != nil {
		log.Error("Failed to decode body", "error", err)

		return ""
	}

	return stringifyObject(
		(body),
	)
}
