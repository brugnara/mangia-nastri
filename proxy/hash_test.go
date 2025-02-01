package proxy

import (
	"mangia_nastri/conf"
	"testing"
)

var p = New(&conf.Config{})

func TestHashReturnsEmptyStringWhenInputIsInvalidJSON(t *testing.T) {
	if p.hash("") != "" {
		t.Fail()
	}

	if p.hash("a") != "" {
		t.Fail()
	}
}

func TestDifferentHashesForDifferentInputs(t *testing.T) {
	if p.hash(`{"a": "ciao"}`) == p.hash(`{"a": "hello"}`) {
		t.Fail()
	}
}

func TestDifferentJSONShapeSameContentSameHash(t *testing.T) {
	if p.hash(`{"a": 1, "b": 2}`) != p.hash(`{"b": 2, "a": 1}`) {
		t.Fail()
	}
}

func TestComplexJSONSameContentSameHash(t *testing.T) {
	if p.hash(`{"a": 1, "b": {"c": 3}}`) != p.hash(`{"b": {"c": 3}, "a": 1}`) {
		t.Fail()
	}
}

func TestJSONWithArrayContentDifferentOrderDifferentHash(t *testing.T) {
	if p.hash(`{"a": 1, "b": [1, 2, 3]}`) == p.hash(`{"b": [3, 1, 2], "a": 1}`) {
		t.Fail()
	}
}
func TestJSONWithArrayContentDifferentOrderSameHash(t *testing.T) {
	if p.hash(`{"a": 1, "b": [1, 2, 3]}`) != p.hash(`{"b": [1, 2, 3], "a": 1}`) {
		t.Fail()
	}
}
