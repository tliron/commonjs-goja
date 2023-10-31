package commonjs

import (
	"github.com/dop251/goja"
)

func ToObject(runtime *goja.Runtime, value any) *goja.Object {
	return ToObjectFromValue(runtime, runtime.ToValue(value))
}

func ToObjectFromValue(runtime *goja.Runtime, value goja.Value) *goja.Object {
	object := runtime.NewObject()
	value_ := value.ToObject(runtime)
	for _, key := range value_.Keys() {
		object.Set(key, value_.Get(key))
	}
	return object
}
