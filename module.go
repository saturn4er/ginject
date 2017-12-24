package ginject

import (
	"reflect"

	"github.com/pkg/errors"
)

type module struct {
	name     string
	module   interface{}
	ptrType  reflect.Type
	ptrValue reflect.Value
	t        reflect.Type
	val      reflect.Value
	fields   []*moduleField
}

type moduleField struct {
	module   *module
	injected bool
	field    reflect.StructField
	value    reflect.Value
	name     string
}

func NewModule(name string, m interface{}) (*module, error) {
	ptrType := reflect.TypeOf(m)
	if !isStructPtr(ptrType) {
		return nil, errors.New("module should be pointer to structure")
	}
	ptrValue := reflect.ValueOf(m)
	if ptrValue.IsNil() {
		return nil, errors.New("passed pointer to nil")
	}

	result := &module{
		module:   m,
		name:     name,
		ptrType:  ptrType,
		ptrValue: ptrValue,
		val:      ptrValue.Elem(),
		t:        ptrType.Elem(),
	}
	for i := 0; i < result.t.NumField(); i++ {
		fieldVal := result.val.Field(i)
		field := result.t.Field(i)
		if tag, ok := field.Tag.Lookup("inject"); ok {
			result.fields = append(result.fields, &moduleField{
				module: result,
				field:  field,
				value:  fieldVal,
				name:   tag,
			})
		}
	}
	return result, nil
}
