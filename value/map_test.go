package value

import (
	"testing"
)

func TestMapString(t *testing.T) {
	if (Map{}).String() != "{}" {
		t.Fatal("expected empty Map to string like this")
	}

	if (Map{"foo": Int(10)}).String() != "{foo:10}" {
		t.Fatal("expected one element Map to string like this")
	}

	if (Map{"foo": Int(10), "bar": String("foo")}).String() != "{bar:foo, foo:10}" {
		t.Fatal("expected two element Map to string like this")
	}
}
