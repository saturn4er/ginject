package ginject

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"errors"
)

type IModule4 interface {
	AndAnotherMethodFromModule4() error
}
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
	WeWantToPopulateThis IModule2 `inject:""`
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
	M41                  IModule4 `inject:"module4 1"`
	M42                  IModule4 `inject:"module4 2"`
}

func (m *Module3) AndAnotherMethodFromModule3() error {
	return nil
}
func (m *Module3) OnInjected() error {
	m.onInjectedCalled = true
	return nil
}

type Module4 struct {
	Version string
}

func (m *Module4) AndAnotherMethodFromModule4() error {
	return nil
}

type OnInjectedErrModule struct {
	WeWantToPopulateThis IModule1 `inject:""`
	onInjectedErr        error
}

func (o *OnInjectedErrModule) OnInjected() error {
	return o.onInjectedErr
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
		var m21 = new(Module2)
		var m22 = new(Module2)
		graph.ShouldAddModules(m1, m21, m22)
		err := graph.Populate()
		So(err, ShouldResemble, ErrMultipleInjectionsFound{Module: m1, Field: "WeWantToPopulateThis", Injection1: m21, Injection2: m22}, )
	})
	Convey("Should panic, if we found multiple factory injections", t, func() {
		graph := Graph{Debug: true}
		var m1 = new(Module1)
		var m2 = new(Module2)
		var m3 = new(Module3)
		graph.ShouldAddModules(m1, m3)
		var m2Factory = func() *Module2 {
			return m2
		}
		var m2Factory2 = func() *Module2 {
			return m2
		}
		graph.ShouldAddFactory(m2Factory)
		graph.ShouldAddFactory(m2Factory2)
		err := graph.Populate()
		e := err.(ErrMultipleInjectionsFound)
		So(e.Module, ShouldEqual, m1)
		So(e.Field, ShouldEqual, "WeWantToPopulateThis")
		So(e.Injection1, ShouldEqual, m2Factory)
		So(e.Injection2, ShouldEqual, m2Factory2)
	})
	Convey("Should panic, if we found multiple injections(module+factory)", t, func() {
		graph := Graph{Debug: true}
		var m1 = new(Module1)
		var m2 = new(Module2)
		var m3 = new(Module3)
		graph.ShouldAddModules(m1, m2, m3)
		var m2Factory = func() *Module2 {
			return m2
		}
		graph.ShouldAddFactory(m2Factory)
		err := graph.Populate()
		e := err.(ErrMultipleInjectionsFound)
		So(e.Module, ShouldEqual, m1)
		So(e.Field, ShouldEqual, "WeWantToPopulateThis")
		So(e.Injection1, ShouldEqual, m2)
		So(e.Injection2, ShouldEqual, m2Factory)
	})
	Convey("Should return error, if factory return error", t, func() {

		m1 := new(Module1)
		factoryErr := errors.New("test error")
		graph := Graph{Debug: true}
		graph.ShouldAddModules(m1)
		var m2Factory = func() (*Module2, error) {
			return nil, factoryErr
		}
		graph.ShouldAddFactory(m2Factory)
		err := graph.Populate().(ErrFactoryReturnError)
		So(err.Factory, ShouldEqual, m2Factory)
		So(err.Err, ShouldEqual, factoryErr)
	})
	Convey("Shouldn't add two modules with same name", t, func() {
		graph := Graph{Debug: true}
		var m1 = new(Module1)
		err := graph.AddNamedModule("1", m1)
		So(err, ShouldBeNil)
		err = graph.AddNamedModule("1", m1)
		So(err, ShouldResemble, ErrNamedModuleAlreadyExists{Name: "1"})
	})
	Convey("Should return error, if no named module found", t, func() {
		graph := Graph{Debug: true}
		var m1 = new(Module1)
		var m2 = new(Module2)
		var m3 = new(Module3)
		err := graph.AddModules(m1, m2, m3)
		So(err, ShouldBeNil)
		err = graph.Populate()
		So(err, ShouldResemble, ErrCantFindInjection{Module: m3, Field: "M41", Name: "module4 1"})
	})
	Convey("Should return error on Populate, if OnInjected hook return error", t, func() {
		graph := Graph{Debug: true}
		var m1 = new(Module1)
		var m2 = new(Module2)
		var m3 = new(OnInjectedErrModule)
		m3.onInjectedErr = errors.New("test error")
		err := graph.AddModules(m1, m2, m3)
		So(err, ShouldBeNil)
		err = graph.Populate()
		So(err, ShouldResemble, ErrOnInjectedHookError{Module: m3, Err: m3.onInjectedErr})
	})
	Convey("Should test normal Graph population", t, func() {
		graph := Graph{Debug: true}
		var m1 = new(Module1)
		var m2 = new(Module2)
		var m3 = new(Module3)
		var m41 = &Module4{"1"}
		var m42 = &Module4{"2"}
		err := graph.AddModules(m1, m3)
		So(err, ShouldBeNil)
		err = graph.AddNamedModule("module4 1", m41)
		So(err, ShouldBeNil)
		err = graph.AddNamedModule("module4 2", m42)
		So(err, ShouldBeNil)

		graph.ShouldAddFactory(func(m interface{}) (*Module2, error) {
			So(m, ShouldEqual, m1)
			return m2, nil
		})
		So(func() {
			graph.ShouldPopulate()
		}, ShouldNotPanic)
		So(m1.WeWantToPopulateThis, ShouldEqual, m2)
		So(m3.M41, ShouldEqual, m41)
		So(m3.M42, ShouldEqual, m42)
		So(m3.onInjectedCalled, ShouldBeTrue)
	})
}
