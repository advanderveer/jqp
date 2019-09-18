package jqp

import (
	"errors"
	"reflect"
)

func unmarshal(src interface{}, rv reflect.Value) (err error) {
	typ := rv.Type()
	if typ.Kind() != reflect.Struct {
		return errors.New("jqp/unmarshal: value must be pointer to a struct")
	}

	for i := 0; i < typ.NumField(); i++ {
		sf := typ.Field(i)
		q, hasTag := sf.Tag.Lookup("jqp")
		if !hasTag {
			continue //no tag
		}

		fv := rv.Field(i)
		if !fv.CanSet() {
			return errors.New("jqp/unmarshal: field '" + sf.Name + "' cannot be set, must be exported")
		}

		// @TODo have and handle query errors
		src = Query(q, src)
		qv := reflect.ValueOf(src)

		switch fv.Kind() {
		case reflect.Struct:
			err = unmarshal(src, fv)
			if err != nil {
				return
			}
		case reflect.Slice: //@TODO how about arrays?

			//

			srcs, ok := src.([]interface{})
			if !ok {
				return errors.New("jqp/unmarshal: query resulted in a '" + qv.Type().String() + "' but the field '" + sf.Name + "' is of type " + fv.Kind().String())
			}

			fv.Set(reflect.MakeSlice(fv.Type(), len(srcs), len(srcs)))
			for j := 0; j < len(srcs); j++ {
				err = unmarshal(srcs[i], fv.Index(j))
				if err != nil {
					return
				}
			}

			// _ = srcs
			// panic("not implemented")

		// @TODO handle slice of: structs, basic types
		// @TODO handle other kinds: https://godoc.org/reflect#Kind
		default:
			if !qv.Type().AssignableTo(fv.Type()) {
				return errors.New("jqp/unmarshal: query resulted in a '" + qv.Type().String() + "' but it is not assignable to field '" + sf.Name + "' of type " + fv.Kind().String())
			}

			// finally, set the field
			fv.Set(qv)
		}
	}

	return
}

// Unmarshal will read data from 'src' into the value pointed to by 'v' using
// the queries in its field tags.
func Unmarshal(src interface{}, v interface{}) (err error) {

	// as seen on: https://golang.org/src/encoding/json/decode.go?s=4043:4091#L170
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("jqp/unmarshal: value must be a pointer and not nil")
	}

	return unmarshal(src, rv.Elem())
}
