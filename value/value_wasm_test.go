// +build wasm

package value_test

import (
	"syscall/js"
	"testing"

	"github.com/advanderveer/jqp"
	"github.com/advanderveer/jqp/value"
)

func TestJavaScriptQuery(t *testing.T) {
	v := jqp.Query(`$.window.location.href`, value.FromJS(js.Global()))
	if v.String() == "" {
		t.Fatal("expected to js query to return something")
	}

	// should correctly convert js value to an int (shrinking from float)
	// and execute the addition
	v = jqp.Query(`$+11`, value.FromJS(js.ValueOf(10)))
	if v.String() != "21" {
		t.Fatalf("unexpected query result, got: %s", v)
	}

	v = jqp.Query(`$[0]`, value.FromJS(js.ValueOf([]interface{}{100})))
	if v.String() != "100" {
		t.Fatalf("unexpected query result, got: %s", v)
	}
}

func TestJavascriptCalling(t *testing.T) {
	v := jqp.Query(`$.JSON.parse('{"foo": "bar"}').foo`, value.FromJS(js.Global()))
	if v.String() != "bar" {
		t.Fatalf("unexpected query result, got: %s", v)
	}
}
