package ginject

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
)


type OnInjecter interface {
	OnInjected() error
}
type Graph struct {
	Debug     bool
	modules   []*module
	factories []*factory
	names     map[string]bool
}

func (g *Graph) ShouldAddModules(modules ...interface{}) {
	err := g.AddModules(modules...)
	if err != nil {
		panic(err)
	}
}
func (g *Graph) AddModules(modules ...interface{}) error {
	for _, m := range modules {
		err := g.AddModule(m)
		if err != nil {
			return errors.Wrapf(err, "failed to add module %T", m)
		}
	}
	return nil
}
func (g *Graph) ShouldAddModule(module interface{}) {
	err := g.AddModule(module)
	if err != nil {
		panic(err)
	}
}
func (g *Graph) AddModule(i interface{}) error {
	return g.AddNamedModule("", i)
}
func (g *Graph) AddNamedModule(name string, i interface{}) error {
	if name != "" {
		if g.names == nil {
			g.names = make(map[string]bool)
		}
		if _, ok := g.names[name]; ok {
			return ErrNamedModuleAlreadyExists{Name: name}
		}
	}
	m, err := NewModule(name, i)
	if err != nil {
		return err
	}
	if name != "" {
		g.names[name] = true
	}
	g.modules = append(g.modules, m)
	return nil
}
func (g *Graph) ShouldAddFactory(factory interface{}) {
	err := g.AddFactory(factory)
	if err != nil {
		panic(err)
	}
}
func (g *Graph) AddFactory(i interface{}) error {
	f, err := NewFactory(i)
	if err != nil {
		return err
	}
	g.factories = append(g.factories, f)
	return nil
}
func (g *Graph) ShouldPopulate() {
	err := g.Populate()
	if err != nil {
		panic(err)
	}
}
func (g *Graph) Populate() error {
	for _, m := range g.modules {
		if g.Debug {
			fmt.Printf("Populating module %T:\n", m.module)
		}
		for _, field := range m.fields {
			if g.Debug {
				fmt.Printf("\tPopulating field %s: \n", field.field.Name)
			}
			var err error
			if field.name == "" {
				err = g.populateField(field)
			} else {
				err = g.populateNamedField(field)
			}

			if err != nil {
				return err
			}
		}
		if i, ok := m.module.(OnInjecter); ok {
			if g.Debug {
				fmt.Printf("\tHandling OnInjected method of %v: ", m.t.Name())
			}
			err := i.OnInjected()
			if err != nil {
				if g.Debug {
					fmt.Println(err)
				}
				return ErrOnInjectedHookError{Module: m.module, Err: err}
			} else if g.Debug {
				fmt.Println("✓")
			}
		} else if g.Debug {
			if g.Debug {
				fmt.Println("\tOnInjected method wasn't found")
			}
		}
	}
	return nil
}
func (g *Graph) populateNamedField(field *moduleField) error {
	for _, m1 := range g.modules {
		if m1 == field.module {
			if g.Debug {
				fmt.Printf("\t\tSkipping %T: can't inject to same module\n", m1.module)
			}
			continue
		}
		if m1.ptrType.AssignableTo(field.field.Type) && m1.name == field.name {
			if g.Debug {
				fmt.Printf("\t\tFound named module %T with name '%s'\n", m1.module, m1.name)
			}
			field.injected = true
			field.value.Set(m1.ptrValue)
			if g.Debug {
				fmt.Printf("\t\tInjected %v\n", m1.ptrValue.Type().String())
			}
			return nil
		} else if g.Debug {
			fmt.Printf("\t\tSkipping %T: %v is not assignable to %v\n", m1.module, m1.ptrType, field.field.Type)
		}
	}
	if g.Debug {
		fmt.Println("\t\tInjection not found :(")
	}
	return ErrCantFindInjection{Module: field.module.module, Field: field.field.Name, Name: field.name}
}
func (g *Graph) populateField(field *moduleField) (err error) {
	var found bool
	if g.Debug {
		fmt.Println("\t\tLooking for assignable module")
	}
	var val reflect.Value
	for _, m1 := range g.modules {
		if m1 == field.module {
			if g.Debug {
				fmt.Printf("\t\t\tSkipping %T: can't inject to same module\n", m1.module)
			}
			continue
		}
		if m1.ptrType.AssignableTo(field.field.Type) {
			if g.Debug {
				fmt.Printf("\t\t\tFound %T\n", m1.module)
			}
			if found {
				if g.Debug {
					fmt.Println("\t\t\tMultiple assignments found :(")
				}
				return ErrMultipleInjectionsFound{Module: field.module.module, Field: field.field.Name, Injection1: val.Interface(), Injection2: m1.module}
			}
			val, found = m1.ptrValue, true
		} else if g.Debug {
			fmt.Printf("\t\t\tSkipping %T: %v is not assignable to %v\n", m1.module, m1.ptrType, field.field.Type)
		}

	}
	if g.Debug {
		fmt.Println("\t\tLooking for assignable factory value")
	}
	var factoryFound interface{}
	for _, f := range g.factories {
		if f.returnValType.AssignableTo(field.field.Type) {
			if g.Debug {
				fmt.Printf("\t\t\tFound value %v of function %T result\n", f.returnValType.String(), f.factory)
			}
			if found {
				if g.Debug {
					fmt.Println("\t\t\tMultiple assignments found :(")
				}
				if factoryFound != nil {
					return ErrMultipleInjectionsFound{Module: field.module.module, Field: field.field.Name, Injection1: factoryFound, Injection2: f.factory}
				} else {
					return ErrMultipleInjectionsFound{Module: field.module.module, Field: field.field.Name, Injection1: val.Interface(), Injection2: f.factory}
				}
			}
			v, err := f.Call(field.module.module)
			if err != nil {
				if g.Debug {
					fmt.Printf("\t\t\tFailed to execute factory function: %v :(\n", err)
				}
				return ErrFactoryReturnError{Factory: f.factory, Err: err}
			}
			factoryFound = f.factory
			val, found = v, true
		}
	}
	if !found {
		if g.Debug {
			fmt.Println("\t\tInjection not found :(")
		}
		return ErrCantFindInjection{Module: field.module.module, Field: field.field.Name, Name: field.name}
	}
	field.injected = true
	field.value.Set(val)
	if g.Debug {
		fmt.Printf("\t\tInjected %v\n", val.Type().String())
	}
	return
}
