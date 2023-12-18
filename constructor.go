package commonjs

import (
	"github.com/dop251/goja"
)

type JavaScriptConstructorFunc = func(constructor goja.ConstructorCall) *goja.Object

func NewConstructor(runtime *goja.Runtime, f func(constructor goja.ConstructorCall) (any, error)) JavaScriptConstructorFunc {
	return func(constructor goja.ConstructorCall) *goja.Object {
		if r, err := f(constructor); err == nil {
			return runtime.ToValue(r).ToObject(runtime)
		} else {
			panic(runtime.NewGoError(err))
		}
	}
}
