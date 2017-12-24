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
	a.Log("test")
}