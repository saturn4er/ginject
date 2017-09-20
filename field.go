package inject

import (
	"reflect"
)

type fieldToInject struct {
	injected bool
	field    reflect.StructField
	value    reflect.Value
	name     string
}
