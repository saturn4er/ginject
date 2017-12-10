package inject

import (
	"reflect"
)

func isStructPtr(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct
}

var errorInterface = reflect.TypeOf((*error)(nil)).Elem()
var sliceOfEmptyInterface []interface{}
var emptyInterfaceType = reflect.TypeOf(sliceOfEmptyInterface).Elem()

func isFactoryFunc(t reflect.Type) bool {
	if t.Kind() != reflect.Func {
		return false
	}
	switch t.NumIn() {
	case 0:
	case 1:
		if !t.In(0).Implements(emptyInterfaceType) {
			return false
		}
	default:
		return false
	}
	switch t.NumOut() {
	case 1:
	case 2:
		if !t.Out(1).Implements(errorInterface) {
			return false
		}
	case 0:
		return false
	default:
		return false
	}
	return true
}
