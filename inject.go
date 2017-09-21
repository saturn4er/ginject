package inject

import (
	"reflect"

	"github.com/pkg/errors"
)

type OnInjecter interface {
	OnInjected() error
}
type Graph struct {
	modules   []*module
	factories []*factory
}

func (g *Graph) AddModules(modules ...interface{}) error {
	for _, m := range modules {
		err := g.AddModule(m)
		if err != nil {
			return errors.Wrapf(err, "failed to add module %v", m)
		}
	}
	return nil
}
func (g *Graph) AddModule(i interface{}) error {
	return g.AddNamedModule("", i)
}
func (g *Graph) AddNamedModule(name string, i interface{}) error {
	m, err := NewModule(name, i)
	if err != nil {
		return err
	}
	g.modules = append(g.modules, m)
	return nil
}
func (g *Graph) AddFactory(i interface{}) error {
	f, err := NewFactory(i)
	if err != nil {
		return err
	}
	g.factories = append(g.factories, f)
	return nil
}
func (g *Graph) Populate() error {
	for _, m := range g.modules {
		for _, field := range m.fieldsToInject {
			if field.injected {
				continue
			}
			var foundType reflect.Type
			var foundValue reflect.Value
			var found bool
			for _, m1 := range g.modules {
				if m1 == m {
					continue
				}
				if m1.ptrType.AssignableTo(field.field.Type) {
					if found {
						return errors.Errorf("found multiple injections for module %T field %s: %v, %v", m.module, field.field.Name, foundType, m1)
					}
					foundType, foundValue, found = m1.ptrType, m1.ptrValue, true
				}
			}
			for _, f := range g.factories {
				if f.returnValType.AssignableTo(field.field.Type) {
					if found {
						return errors.Errorf("found multiple injections for module %T field %s: %v, %v", m.module, field.field.Name, foundType, f.factoryType)
					}
					result, err := f.Call(m.module)
					if err != nil {
						return errors.Wrapf(err, "factory %v returns error", f.factoryType)
					}
					foundType, foundValue, found = f.factoryType, result, true
				}

			}
			if !found {
				return errors.Errorf("failed to find injection for module %T field %s", m.module, field.field.Name)
			}
			field.injected = true
			field.value.Set(foundValue)
		}
		if i, ok := m.module.(OnInjecter); ok {
			err := i.OnInjected()
			if err != nil {
				return errors.Wrapf(err, "failed to handle on OnInjected hook of module %T", m.module)
			}
		}
	}
	return nil
}
