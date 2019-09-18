package jqp_test

import (
	"testing"

	"github.com/advanderveer/jqp"
)

func TestSingleFieldUnmarshal(t *testing.T) {
	src := "foo"
	type A struct {
		Bar string `jqp:"$"`
	}

	v := A{}
	if err := jqp.Unmarshal(src, &v); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if v.Bar != "foo" {
		t.Fatalf("unmarshal didn't yield correct value, got: %v", v.Bar)
	}
}

func TestNestedStructUnmarshal(t *testing.T) {
	src := map[string]interface{}{
		"b": map[string]interface{}{
			"c": map[string]interface{}{
				"bar": "foo",
			}},
	}

	type A struct {
		B struct {
			C struct {
				Bar string `jqp:"$.bar"`
			} `jqp:"$.c"`
		} `jqp:"$.b"`
	}

	v := A{}
	if err := jqp.Unmarshal(src, &v); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if v.B.C.Bar != "foo" {
		t.Fatalf("unmarshal didn't yield correct value, got: %v", v.B.C.Bar)
	}
}
