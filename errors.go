package ginject

import "fmt"

type ErrOnInjectedHookError struct {
	Module interface{}
	Err    error
}

func (e ErrOnInjectedHookError) Error() string {
	return fmt.Sprintf("OnInjected hook of module %T returns error: %v", e.Module, e.Err)
}

type ErrNamedModuleAlreadyExists struct {
	Name string
}

func (e ErrNamedModuleAlreadyExists) Error() string {
	return fmt.Sprintf("module with name %s already exists", e.Name)
}

type ErrFactoryReturnError struct {
	Factory interface{}
	Err     error
}

func (e ErrFactoryReturnError) Error() string {
	return fmt.Sprintf("factory %T returns error: %v", e.Factory, e.Err)
}

type ErrMultipleInjectionsFound struct {
	Module     interface{}
	Field      string
	Injection1 interface{}
	Injection2 interface{}
}

func (e ErrMultipleInjectionsFound) Error() string {
	return fmt.Sprintf("found multiple injections for %T.%s: %#v, %#v", e.Module, e.Field, e.Injection1, e.Injection2)
}

type ErrCantFindInjection struct {
	Module interface{}
	Field  string
	Name   string
}

func (e ErrCantFindInjection) Error() string {
	if e.Name == "" {
		return fmt.Sprintf("can't find injection for %T.%s", e.Module, e.Field)
	} else {
		return fmt.Sprintf("can't find injection for %T.%s (name: %s)", e.Module, e.Field, e.Name)
	}
}
