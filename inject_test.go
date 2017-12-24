package ginject

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"errors"
)

type IModule3 interface {
	AndAnotherMethodFromModule3() error
}
type IModule2 interface {
	AndAnotherMethodFromModule2() error
}
type IModule1 interface {
	SomeMethodFromModule1() error
	SomeOtherMethodFromModule1() error
}

type Module1 struct {
	WeWantToPopulateThis     IModule2 `inject:""`
	WeAlsoWantToPopulateThis IModule3 `inject:""`
}

func (m *Module1) SomeMethodFromModule1() error {
	return nil
}

func (m *Module1) SomeOtherMethodFromModule1() error {
	return nil
}

type Module2 struct {
	onInjectedCalled     bool
	WeWantToPopulateThis IModule1 `inject:""`
}

func (m *Module2) AndAnotherMethodFromModule2() error {
	return nil
}
func (m *Module2) OnInjected() error {
	m.onInjectedCalled = true
	return nil
}

type Module3 struct {
	onInjectedCalled     bool
	WeWantToPopulateThis IModule1 `inject:""`
}

func (m *Module3) AndAnotherMethodFromModule3() error {
	return nil
}
func (m *Module3) OnInjected() error {
	m.onInjectedCalled = true
	return nil
}

func TestInject(t *testing.T) {
	Convey("Should test Graph initialization with bad modules", t, func() {
		graph := Graph{Debug: true}
		Convey("Should panic, if we are trying to add non-pointer module to graph", func() {
			So(func() {
				graph.ShouldAddModules(Module1{})
			}, ShouldPanic)
		})
		Convey("Should panic, if we are trying to add module, which is a pointer to nil", func() {
			So(func() {
				var m *Module1
				graph.ShouldAddModule(m)
			}, ShouldPanic)
		})
	})
	Convey("Should test Graph initialization with bad factories", t, func() {
		graph := Graph{Debug: true}
		Convey("Should panic, if we are trying to add bad factory", func() {
			So(func() {
				graph.ShouldAddFactory(Module1{})
			}, ShouldPanic)

			// Inputs
			So(func() {
				graph.ShouldAddFactory(func(m *Module2) (*Module1, error) {
					return new(Module1), nil
				})
			}, ShouldPanic)
			So(func() {
				graph.ShouldAddFactory(func(m, m1 interface{}) (*Module1, error) {
					return new(Module1), nil
				})
			}, ShouldPanic)
			// Outputs
			So(func() {
				graph.ShouldAddFactory(func() {})
			}, ShouldPanic)
			So(func() {
				graph.ShouldAddFactory(func(m interface{}) (*Module1, *Module1) {
					return new(Module1), nil
				})
			}, ShouldPanic)

			So(func() {
				graph.ShouldAddFactory(func() (bool, error, int) {
					return false, nil, 2
				})
			}, ShouldPanic)

		})
	})
	Convey("Should panic, if we can't find population for some field", t, func() {
		graph := Graph{Debug: true}
		var m1 = new(Module1)
		So(func() {
			graph.ShouldAddModule(m1)
		}, ShouldNotPanic)
		So(func() {
			graph.ShouldPopulate()
		}, ShouldPanic)
	})
	Convey("Should panic, if we found multiple  module injections", t, func() {
		graph := Graph{Debug: true}
		var m1 = new(Module1)
		var m2 = new(Module2)
		graph.ShouldAddModules(m1, m2, m2)
		So(func() {
			graph.ShouldPopulate()
		}, ShouldPanic)
	})
	Convey("Should panic, if we found multiple factory injections", t, func() {
		graph := Graph{Debug: true}
		var m1 = new(Module1)
		graph.ShouldAddModules(m1)
		var m2Factory = func() *Module2 {
			return new(Module2)
		}
		graph.ShouldAddFactory(m2Factory)
		graph.ShouldAddFactory(m2Factory)

		So(func() {
			graph.ShouldPopulate()
		}, ShouldPanic)
	})
	Convey("Should panic, if factory return error", t, func() {
		graph := Graph{Debug: true}
		var m1 = new(Module1)
		graph.ShouldAddModules(m1)
		var m2Factory = func() (*Module2, error) {
			return nil, errors.New("test error")
		}
		graph.ShouldAddFactory(m2Factory)

		So(func() {
			graph.ShouldPopulate()
		}, ShouldPanic)
	})
	Convey("Should test normal Graph population", t, func() {
		graph := Graph{Debug: true}
		var m1 = new(Module1)
		var m2 = new(Module2)
		var m3 = new(Module3)
		err := graph.AddModules(m1, m3)
		So(err, ShouldBeNil)
		graph.ShouldAddFactory(func(m interface{}) (*Module2, error) {
			So(m, ShouldEqual, m1)
			return m2, nil
		})
		So(func() {
			graph.ShouldPopulate()
		}, ShouldNotPanic)
		So(m1.WeWantToPopulateThis, ShouldEqual, m2)
		So(m1.WeAlsoWantToPopulateThis, ShouldEqual, m3)
		So(m3.onInjectedCalled, ShouldBeTrue)
	})
}
