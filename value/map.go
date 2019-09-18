package value

import (
	"sort"
	"strings"
)

// Map maps strings to values
type Map map[string]Value

var _ Value = Map{}

func (m Map) String() string {
	vals := make([]string, 0, len(m))
	for k, v := range m {
		vals = append(vals, k+":"+v.String())
	}

	sort.Strings(vals)
	return "{" + strings.Join(vals, ", ") + "}"
}

func (m Map) Eval(ctx Context) Value {
	return m
}

func (m Map) whichType() valueType {
	return mapType
}

func (m Map) toType(which valueType) Value {
	switch which {
	case mapType:
		return m
	default:
		panic("type coversion from '" + m.whichType().String() + "' to '" + which.String() + "' not implemented")
	}
}
