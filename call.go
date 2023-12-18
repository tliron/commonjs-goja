package commonjs

import (
	"fmt"

	"github.com/dop251/goja"
)

// This is the returned type when calling Export() on a [goja.FunctionCall].
type ExportedJavaScriptFunc = func(functionCall goja.FunctionCall) goja.Value

// The function argument can be a [ExportedJavaScriptFunc] or a [goja.Value] representing a function.
func Call(runtime *goja.Runtime, function any, this any, arguments ...any) (value any, err error) {
	switch function_ := function.(type) {
	case goja.Value:
		return CallValue(runtime, function_, this, arguments...)

	case ExportedJavaScriptFunc:
		return CallExported(runtime, function_, this, arguments...)

	default:
		return nil, fmt.Errorf("not a function: %v", function)
	}
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

func CallValue(runtime *goja.Runtime, function goja.Value, this any, arguments ...any) (any, error) {
	if function_, ok := goja.AssertFunction(function); ok {
		if r, err := function_(runtime.ToValue(this), toValues(runtime, arguments)...); err == nil {
			return r.Export(), nil
		} else {
			return nil, UnwrapJavaScriptException(err)
		}
	} else {
		return nil, fmt.Errorf("not a function: %v", function)
	}
}

func CallExported(runtime *goja.Runtime, function ExportedJavaScriptFunc, this any, arguments ...any) (value any, err error) {
	defer func() {
		if err_ := HandleJavaScriptPanic(recover()); err_ != nil {
			err = err_
		}
	}()

	functionCall := goja.FunctionCall{
		This:      runtime.ToValue(this),
		Arguments: toValues(runtime, arguments),
	}

	return function(functionCall).Export(), nil // can panic
}
