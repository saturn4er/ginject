## ginject

`ginject` is a library, that will help you inject dependencies to your modules


### Example

```go
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

func A() *Module1 {
	return new(Module1)
}
func main() {
	var m2 Module2
	graph := ginject.Graph{}
	graph.ShouldAddModules(&m2)
	graph.ShouldAddFactory(A)
	err := graph.Populate()
	if err != nil {
		panic(err)
	}
}

```

### Factories

Factories can provide modules specific dependency. For example, if you need logger for module A which will print output in format `[module:a] log-message`, you can add logs factory to your graph:

```go
package main
import (
	"github.com/saturn4er/ginject"
	"reflect"
	"fmt"
)
type A struct{
	Log func(msg string) `inject:""`
}
func main(){
	var a = new(A)
	graph := ginject.Graph{}
	graph.ShouldAddModule(a)
	graph.ShouldAddFactory(func(m interface{}) func(msg string){
		module:= reflect.TypeOf(m).Elem().Name()
		return func(msg string){
			fmt.Printf("[module:%s] %s\n", module, msg)
		}
	})
	graph.ShouldPopulate()
	a.Log("test") // Should output "[module:A] test"
}
```
#### Factories types

Valid values of factories is:
- func() T
- func() (T, error)
- func(m interface{}) T
- func(m interface{}) )T, error)