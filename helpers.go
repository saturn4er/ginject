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
	if t.NumIn() != 1 || !t.In(0).Implements(emptyInterfaceType) {
		return false
	}

	if t.NumOut() < 2 || !t.Out(1).Implements(errorInterface) {
		return false
	}
	return true
}
