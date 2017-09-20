package inject

import (
	"reflect"
	"errors"
)

type factory struct {
	factory       interface{}
	factoryType   reflect.Type
	factoryValue  reflect.Value
	returnValType reflect.Type
}

func (f *factory) Call(module interface{}) (reflect.Value, error) {
	out := f.factoryValue.Call([]reflect.Value{reflect.ValueOf(module)})
	var err error
	errI := out[1].Interface()
	if errI != nil {
		err = errI.(error)
	}
	return out[0], err
}

func NewFactory(i interface{}) (*factory, error) {
	ft := reflect.TypeOf(i)
	if !isFactoryFunc(ft) {
		return nil, errors.New("bad factory type. Should be of func(interface{})(T, error)")
	}
	result := &factory{
		factory:       i,
		factoryType:   ft,
		factoryValue:  reflect.ValueOf(i),
		returnValType: ft.Out(0),
	}
	return result, nil
}
