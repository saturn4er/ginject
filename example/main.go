package main

import (
	"github.com/saturn4er/ginject"
	"fmt"
)

type IModule2 interface {
	AndAnotherMethodFromModule2() error
}
type IModule1 interface {
	SomeMethodFromModule1() error
	SomeOtherMethodFromModule1() error
}

type Module1 struct {
	Env                  string
	WeWantToPopulateThis IModule2 `inject:""`
}

func (m *Module1) SomeMethodFromModule1() error {
	return nil
}

func (m *Module1) SomeOtherMethodFromModule1() error {
	return nil
}

type Module2 struct {
	SqlM1 IModule1 `inject:"m1 sql"`
	PGM1  IModule1 `inject:"m1 pg"`
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
	var m1sql = Module1{Env: "sql"}
	var m1pg = Module1{Env: "pg"}
	var m2 Module2
	graph := ginject.Graph{Debug: true}
	graph.AddNamedModule("m1 sql", &m1sql)
	graph.AddNamedModule("m1 pg", &m1pg)
	graph.ShouldAddModules(&m2)
	graph.ShouldAddFactory(Module1Factory)
	err := graph.Populate()
	if err != nil {
		panic(err)
	}
	fmt.Println(m2.PGM1)
	fmt.Println(m2.SqlM1)
}
