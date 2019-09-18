package value_test

import (
	"github.com/advanderveer/jqp/value"
	"reflect"
	"strings"
	"testing"
)

func TestToNative(t *testing.T) {
	v := value.ToNative(value.String("abc"))
	if !reflect.DeepEqual(v, "abc") {
		t.Fatalf("unexpected to value result, got: %#v", v)
	}

	v = value.ToNative(value.Float(1.5))
	if !reflect.DeepEqual(v, 1.5) {
		t.Fatalf("unexpected to value result, got: %#v", v)
	}

	v = value.ToNative(value.Int(2))
	if !reflect.DeepEqual(v, 2) {
		t.Fatalf("unexpected to value result, got: %#v", v)
	}

	v = value.ToNative(value.Array{value.Int(3), value.String("abc")})
	if !reflect.DeepEqual(v, []interface{}{3, "abc"}) {
		t.Fatalf("unexpected to value result, got: %#v", v)
	}

	v = value.ToNative(value.Map{"a": value.Int(3), "b": value.String("abc")})
	if !reflect.DeepEqual(v, map[string]interface{}{"a": 3, "b": "abc"}) {
		t.Fatalf("unexpected to value result, got: %#v", v)
	}

	v = value.ToNative(value.Func(func(args ...value.Value) value.Value {
		return value.String(strings.ToUpper(string(args[0].(value.String))))
	}))
	if !reflect.DeepEqual(v.(func(...interface{}) interface{})("foo"), "FOO") {
		t.Fatalf("unexpected to value result, got: %#v", v)
	}
}

func TestFromNativeWithoutPorting(t *testing.T) {
	v := value.FromNative("abc", false)
	if v.String() != `abc` {
		t.Fatal("unexpected to value result, got: " + v.String())
	}

	v = value.FromNative(10, false)
	if v.String() != `10` {
		t.Fatal("unexpected to value result, got: " + v.String())
	}

	v = value.FromNative(1.0, false)
	if v.String() != `1E+00` {
		t.Fatal("unexpected to value result, got: " + v.String())
	}

	v = value.FromNative([]interface{}{1, "abc"}, false)
	if v.String() != `[1, abc]` {
		t.Fatal("unexpected to value result, got: " + v.String())
	}

	v = value.FromNative(map[string]interface{}{"foo": 1, "bar": "abc"}, false)
	if v.String() != `{bar:abc, foo:1}` {
		t.Fatal("unexpected to value result, got: " + v.String())
	}

	v = value.FromNative(value.String("foo"), false)
	if v.String() != `foo` {
		t.Fatal("unexpected to value result, got: " + v.String())
	}

	v = value.FromNative(func(args ...interface{}) interface{} { return 100 }, false)
	if v.String() != `func()` {
		t.Fatal("unexpected to value result, got: " + v.String())
	}
}

func TestFromNativeWithPorting(t *testing.T) {
	v := value.FromNative(map[string]interface{}{"foo": 1, "bar": "abc"}, true)
	if v.String() != `{ port }` {
		t.Fatal("expected port value")
	}
}
