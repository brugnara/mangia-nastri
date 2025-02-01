package proxy

import (
	"testing"
)

func TestStringifyWorksAsExpected(t *testing.T) {
	if p.stringifyObject(map[string]interface{}{
		"b": 2,
		"a": 1,
	}) != `{"a":1,"b":2}` {
		t.Fail()
	}
}
