package value

import (
	"testing"
)

func TestFloatShrinking(t *testing.T) {
	if (Float(1.0)).shrink().String() != "1" {
		t.Fatal("expected shrink to cause int")
	}

	if (Float(1.5)).shrink().String() != "1.5E+00" {
		t.Fatalf("expected shrink not to cause shrink, got: %s", (Float(1.5)).shrink())
	}
}
