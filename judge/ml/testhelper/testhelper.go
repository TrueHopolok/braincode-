package testhelper

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func DiffValues(t *testing.T, lhs, rhs any) {
	t.Helper()

	var buf []string
	diffValues(func(s string) {
		buf = append(buf, s)
	}, "root", reflect.ValueOf(lhs), reflect.ValueOf(rhs))

	if buf != nil {
		t.Error("lhs != rhs:\n\n" + strings.Join(buf, "\n\n"))
	}
}

func diffValues(pushError func(string), path string, lhs, rhs reflect.Value) {
	if lhs.Type() != rhs.Type() {
		var l, r string
		if lhs.Type().String() == rhs.Type().String() {
			l = lhs.Type().PkgPath()
			r = rhs.Type().PkgPath()
		} else {
			l = lhs.Type().String()
			r = rhs.Type().String()
		}

		pushError(fmt.Sprintf("%s: type mismatch:\nlhs: %v\nrhs: %v", path, l, r))
		return
	}

	switch lhs.Kind() {
	case reflect.Interface, reflect.Pointer:
		if lhs.IsNil() != rhs.IsNil() {
			pushError(fmt.Sprintf("%s: nilness mismatch:\nlhs: %#v\nrhs: %#v", path, lhs, rhs))
			return
		}

		if lhs.IsNil() {
			return
		}

		var newPath string
		if lhs.Kind() == reflect.Pointer {
			newPath = "(*" + path + ")"
		} else {
			newPath = path + ".(" + lhs.Elem().Type().String() + ")"
		}
		diffValues(pushError, newPath, lhs.Elem(), rhs.Elem())

	case reflect.Array, reflect.Slice:
		if lhs.Len() != rhs.Len() {
			pushError(fmt.Sprintf("%s: len mismatch:\nlen(lhs) = %v (%#v)\nlen(rhs) = %v (%#v)",
				path, lhs.Len(), lhs, rhs.Len(), rhs))
			return
		}

		if lhs.Kind() == reflect.Slice && lhs.IsNil() != rhs.IsNil() {
			pushError(fmt.Sprintf("%s: nilness mismatch:\nlhs: %#v\nrhs: %#v", path, lhs, rhs))
			return
		}

		for i := range lhs.Len() {
			diffValues(pushError, fmt.Sprintf("%s[%d]", path, i), lhs.Index(i), rhs.Index(i))
		}

	case reflect.Map:
		if lhs.IsNil() != rhs.IsNil() {
			pushError(fmt.Sprintf("%s: nilness mismatch:\nlhs: %#v\nrhs: %#v", path, lhs, rhs))
			return
		}

		if lUnique := uniqueKeys(lhs, rhs); len(lUnique) > 0 {
			var acc []string
			for _, v := range lUnique {
				acc = append(acc, fmt.Sprintf("%#v", v))
			}
			pushError(fmt.Sprintf("%s: lhs has extra keys:\n%v", path, strings.Join(acc, ", ")))
		}
		if rUnique := uniqueKeys(rhs, lhs); len(rUnique) > 0 {
			var acc []string
			for _, v := range rUnique {
				acc = append(acc, fmt.Sprintf("%#v", v))
			}
			pushError(fmt.Sprintf("%s: rhs has extra keys:\n%v", path, strings.Join(acc, ", ")))
		}

		for iter := lhs.MapRange(); iter.Next(); {
			lv := iter.Value()
			rv := rhs.MapIndex(iter.Key())
			if !rv.IsValid() {
				continue
			}
			diffValues(pushError, fmt.Sprintf("%s[%#v]", path, iter.Key()), lv, rv)
		}

	case reflect.Struct:
		for i := range lhs.NumField() {
			lv := lhs.Field(i)
			rv := rhs.Field(i)

			diffValues(pushError, fmt.Sprintf("%s.%s", path, lhs.Type().Field(i).Name), lv, rv)
		}

	default:
		if !lhs.Equal(rhs) {
			pushError(fmt.Sprintf("%s: value mismatch:\nlhs: %#v\nrhs: %#v", path, lhs, rhs))
		}
	}
}

func uniqueKeys(l, r reflect.Value) []reflect.Value {
	iter := l.MapRange()
	var res []reflect.Value
	for iter.Next() {
		if !r.MapIndex(iter.Key()).IsValid() {
			res = append(res, iter.Key())
		}
	}
	return res
}
