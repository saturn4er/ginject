package main

import (
	"github.com/saturn4er/ginject"
)

type IModule2 interface {
	AndAnotherMethodFromModule2() error
}
type IModule1 interface {
	SomeMethodFromModule1() error
	SomeOtherMethodFromModule1() error
}

type Module1 struct {
	WeWantToPopulateThis IModule2 `inject:""`
}

func (m *Module1) SomeMethodFromModule1() error {
	return nil
}

func (m *Module1) SomeOtherMethodFromModule1() error {
	return nil
}

type Module2 struct {
	WeWantToPopulateThis IModule1 `inject:""`
}

func (m *Module2) AndAnotherMethodFromModule2() error {
	return nil
}
func (m *Module2) OnInjected() error {
	return nil
}

func Module1Factory() *Module1 {
	return new(Module1)
}
func main() {
	var m2 Module2
	graph := ginject.Graph{Debug: true}
	graph.ShouldAddModules(&m2)
	graph.ShouldAddFactory(Module1Factory)
	err := graph.Populate()
	if err != nil {
		panic(err)
	}
}
