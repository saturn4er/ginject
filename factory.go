package inject

import (
	"errors"
	"reflect"
)

type factory struct {
	factory       interface{}
	factoryType   reflect.Type
	factoryValue  reflect.Value
	returnValType reflect.Type
}

func (f *factory) Call(module interface{}) (reflect.Value, error) {
	var in []reflect.Value
	if f.factoryType.NumIn() > 0 {
		in = []reflect.Value{reflect.ValueOf(module)}
	}
	out := f.factoryValue.Call(in)
	var err error
	if f.factoryType.NumOut() > 1 {
		errI := out[1].Interface()
		if errI != nil {
			err = errI.(error)
		}
	}
	return out[0], err
}

func NewFactory(i interface{}) (*factory, error) {
	ft := reflect.TypeOf(i)
	if !isFactoryFunc(ft) {
		return nil, errors.New("bad factory type. Should be one of: (func() T), (func() T, error) (func(interface{}) T), (func(interface{})(T, error))")
	}
	result := &factory{
		factory:       i,
		factoryType:   ft,
		factoryValue:  reflect.ValueOf(i),
		returnValType: ft.Out(0),
	}
	return result, nil
}
