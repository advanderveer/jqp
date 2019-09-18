package value

import (
	"testing"
)

func TestArrayString(t *testing.T) {
	if (Array{}).String() != "[]" {
		t.Fatal("expected empty array to string like this")
	}

	if (Array{Int(10)}).String() != "[10]" {
		t.Fatal("expected one element array to string like this")
	}

	if (Array{Int(10), String("foo")}).String() != "[10, foo]" {
		t.Fatal("expected one element array to string like this")
	}
}
