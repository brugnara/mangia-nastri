package src

import (
	"encoding/json"
	"io"
	"os"
	"sort"
)

var logger = SetupLogger()

// `sortNestedObject` function accepts a generic mapping with string keys and
// potentially nested values. It returns a new mapping with all items sorted
// alphabetically by key.
// If a value in the mapping is nested, the function is called recursively.
//
// Parameters
//   - obj: the possibly nested object to be sorted.
//
// Returns
//   A new map[string]interface{} with all keys sorted alphabetically,
//   and any nested maps also sorted recursively.

func sortNestedObject(obj map[string]interface{}) map[string]interface{} {
	keys := make([]string, 0, len(obj))
	for k := range obj {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	sortedObj := make(map[string]interface{}, len(obj))
	for _, k := range keys {
		if nestedObj, ok := obj[k].(map[string]interface{}); ok {
			sortedObj[k] = sortNestedObject(nestedObj)
		} else {
			sortedObj[k] = obj[k]
		}
	}

	return sortedObj
}

// `stringifyObject` converts a generic object to its JSON string representation.
// If the marshalling process fails, it logs the error and terminates the program.
//
// Parameters
//   - obj: the generic object to be converted to a JSON string.
//
// Returns
//   A string containing the JSON representation of the input object.

func stringifyObject(obj map[string]interface{}) string {
	bytes, err := json.Marshal(obj)
	if err != nil {
		logger.Error("Failed to marshal sorted body", "error", err)
		os.Exit(1)
	}

	return string(bytes)
}

// `ProcessHeaders` reads the raw headers of an HTTP request and processes them
// into a single string, sorting all key-value pairs alphabetically.
//
// Parameters
//   - rawHeaders: a map containing the raw headers of an HTTP request.
//
// Returns
//   A string containing the processed and sorted headers.

func ProcessHeaders(rawHeaders map[string][]string) string {
	headersMap := make(map[string]interface{}, len(rawHeaders))
	for k, v := range rawHeaders {
		headersMap[k] = v
	}

	return stringifyObject(
		sortNestedObject(headersMap),
	)
}

// `ProcessBody` reads the raw body of an HTTP request and processes it into a
// single string, sorting all key-value pairs alphabetically.
//
// Parameters
//   - rawBody: an object containing the raw body of an HTTP request.
//
// Returns
//   A string containing the processed and sorted body.

func ProcessBody(rawBody io.ReadCloser) string {
	defer rawBody.Close()

	var body map[string]interface{}
	if err := json.NewDecoder(rawBody).Decode(&body); err != nil {
		logger.Error("Failed to decode body", "error", err)
		os.Exit(1)
	}

	return stringifyObject(
		sortNestedObject(body),
	)
}
