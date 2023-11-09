package commonjs

import (
	"fmt"

	"github.com/dop251/goja"
)

// This is the returned type when calling Export() on a [goja.FunctionCall].
type JavaScriptFunc = func(functionCall goja.FunctionCall) goja.Value

// The function argument can be a [JavaScriptFunc] or a [goja.Value] representing a function.
func Call(runtime *goja.Runtime, function any, this any, arguments ...any) (value any, err error) {
	defer func() {
		if err_ := HandleJavaScriptPanic(recover()); err_ != nil {
			err = err_
		}
	}()

	switch function_ := function.(type) {
	case JavaScriptFunc:
		functionCall := goja.FunctionCall{
			This:      runtime.ToValue(this),
			Arguments: toValues(runtime, arguments),
		}

		return function_(functionCall).Export(), nil // can panic

	case goja.Value:
		if function__, ok := goja.AssertFunction(function_); ok {
			if r, err := function__(runtime.ToValue(this), toValues(runtime, arguments)...); err == nil {
				return r.Export(), nil
			} else {
				return nil, UnwrapJavaScriptException(err)
			}
		}
	}

	return nil, fmt.Errorf("not a function: %v", function)
}

func GetAndCall(runtime *goja.Runtime, object *goja.Object, name string, this any, arguments ...any) (value any, err error) {
	defer func() {
		if err_ := HandleJavaScriptPanic(recover()); err_ != nil {
			err = err_
		}
	}()

	if function := object.Get(name); function != nil { // Get can panic
		return Call(runtime, function, this, arguments...)
	} else {
		return "", fmt.Errorf("\"%s\" not found", name)
	}
}

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

// Utils

func toValues(runtime *goja.Runtime, values []any) []goja.Value {
	if length := len(values); length > 0 {
		values_ := make([]goja.Value, length)
		for index, argument := range values {
			values_[index] = runtime.ToValue(argument)
		}
		return values_
	} else {
		return nil
	}
}
