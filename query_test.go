package jqp_test

import (
	"encoding/json"
	"reflect"
	"strconv"
	"testing"

	"github.com/advanderveer/jqp"
)

func TestNativeQuery(t *testing.T) {
	for i, c := range []struct {
		v      interface{}
		query  string
		result interface{}
	}{
		{func(args ...interface{}) interface{} { return "bar" }, `$()`, "bar"},
		{map[string]interface{}{
			"foo": func(args ...interface{}) interface{} { return 1.5 },
		}, `$.foo()`, 1.5},

		{map[string]interface{}{
			"foo": func(args ...interface{}) interface{} { return []interface{}{"a", 1.5, 5} },
		}, `$.foo()`, []interface{}{"a", 1.5, 5}},

		{map[string]interface{}{
			"foo": func(args ...interface{}) interface{} { return map[string]interface{}{"a": "a", "b": 1.5, "c": 5} },
		}, `$.foo()`, map[string]interface{}{"a": "a", "b": 1.5, "c": 5}},

		{[]interface{}{map[string]interface{}{
			"foo": func(args ...interface{}) interface{} { return "bar" },
		}}, `$[0].foo()`, "bar"},

		{[]interface{}{map[string]interface{}{
			"foo": func(args ...interface{}) interface{} { return func(args ...interface{}) interface{} { return "bar" } },
		}}, `$[0].foo()()`, "bar"},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			res := jqp.Query(c.query, c.v)
			if !reflect.DeepEqual(res, c.result) {
				t.Fatalf("query '%s' on '%#v' gave '%#v', expected: '%s'", c.query, c.v, res, c.result)
			}
		})
	}
}

func TestJSONQuery(t *testing.T) {
	for i, c := range []struct {
		json   string
		query  string
		result interface{}
	}{
		{`{"foo": "bar"}`, "$.foo", "bar"},
		{`{"foo": [3,4]}`, "$.foo[0]", 3.0},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var v interface{}
			err := json.Unmarshal([]byte(c.json), &v)
			if err != nil {
				t.Fatal(err)
			}

			res := jqp.Query(c.query, v)
			if !reflect.DeepEqual(res, c.result) {
				t.Fatalf("query '%s' on '%s' gave '%#v', expected: '%s'", c.query, c.json, res, c.result)
			}
		})
	}
}
