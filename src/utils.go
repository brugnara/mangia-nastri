package src

import (
  "sort"
  "strings"
)

func ProcessHeaders(headers map[string][]string) string {
	array := make([]string, 0, len(headers))

	for k := range headers {
		array = append(array, k + ": " + headers[k][0])
	}

	sort.Strings(array)

	return strings.Join(array, ", ")
}
