package proxy

import (
	"mangia_nastri/commander"
	"mangia_nastri/conf"
	"mangia_nastri/logger"
	"testing"
)

var p = New(&conf.Proxy{
	Port:        8080,
	Name:        "test",
	Destination: "http://localhost:8080",
	DataSource: conf.DataSource{
		Type: "inMemory",
	},
}, logger.New("test"), make(chan commander.Action))

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
