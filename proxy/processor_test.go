package proxy

import (
	"testing"
)

var p = New(nil)

func TestStringifyWorksAsExpected(t *testing.T) {
	if p.stringifyObject(map[string]interface{}{
		"b": 2,
		"a": 1,
	}) != `{"a":1,"b":2}` {
		t.Fail()
	}
}

func TestProcessHeadersIgnoresHeaders(t *testing.T) {
	if p.ProcessHeaders(map[string][]string{
		"Authorization": {"Bearer", "123"},
		"Host":          {"localhost"},
	}, []string{"Host"}) != `{"Authorization":["Bearer","123"]}` {
		t.Fail()
	}
}
