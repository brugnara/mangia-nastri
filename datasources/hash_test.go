package datasources

import (
	"testing"
)

func TestHashReturnsEmptyStringWhenInputIsInvalidJSON(t *testing.T) {
	if ComputeHash("") != "" {
		t.Fail()
	}

	if ComputeHash("a") != "" {
		t.Fail()
	}
}

func TestDifferentHashesForDifferentInputs(t *testing.T) {
	if ComputeHash(`{"a": "ciao"}`) == ComputeHash(`{"a": "hello"}`) {
		t.Fail()
	}
}

func TestDifferentJSONShapeSameContentSameHash(t *testing.T) {
	if ComputeHash(`{"a": 1, "b": 2}`) != ComputeHash(`{"b": 2, "a": 1}`) {
		t.Fail()
	}
}

func TestComplexJSONSameContentSameHash(t *testing.T) {
	if ComputeHash(`{"a": 1, "b": {"c": 3}}`) != ComputeHash(`{"b": {"c": 3}, "a": 1}`) {
		t.Fail()
	}
}

func TestJSONWithArrayContentDifferentOrderDifferentHash(t *testing.T) {
	if ComputeHash(`{"a": 1, "b": [1, 2, 3]}`) == ComputeHash(`{"b": [3, 1, 2], "a": 1}`) {
		t.Fail()
	}
}
func TestJSONWithArrayContentDifferentOrderSameHash(t *testing.T) {
	if ComputeHash(`{"a": 1, "b": [1, 2, 3]}`) != ComputeHash(`{"b": [1, 2, 3], "a": 1}`) {
		t.Fail()
	}
}
