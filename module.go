package inject

import (
	"reflect"

	"github.com/pkg/errors"
)

type module struct {
	name           string
	module         interface{}
	ptrType        reflect.Type
	ptrValue       reflect.Value
	moduleType     reflect.Type
	moduleValue    reflect.Value
	fieldsToInject []*fieldToInject
}

func (m *module) IsNamed() bool {
	return m.name != ""
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
		module:      m,
		name:        name,
		ptrType:     ptrType,
		ptrValue:    ptrValue,
		moduleValue: ptrValue.Elem(),
		moduleType:  ptrType.Elem(),
	}
	for i := 0; i < result.moduleType.NumField(); i++ {
		fieldVal := result.moduleValue.Field(i)
		field := result.moduleType.Field(i)
		tag, ok := field.Tag.Lookup("inject")
		if !ok {
			continue
		}
		result.fieldsToInject = append(result.fieldsToInject, &fieldToInject{
			field: field,
			value: fieldVal,
			name:  tag,
		})

	}
	return result, nil
}
