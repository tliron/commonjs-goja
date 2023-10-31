package commonjs

import (
	"fmt"

	"github.com/dop251/goja"
)

// This is returned type when calling Export() on a [goja.FunctionCall].
type JavaScriptFunc = func(goja.FunctionCall) goja.Value

func Call(runtime *goja.Runtime, function JavaScriptFunc, this any, arguments ...any) (value any, err error) {
	defer func() {
		if err_ := HandlePanic(recover()); err_ != nil {
			err = err_
		}
	}()

	var this_ goja.Value
	if this != nil {
		this_ = runtime.ToValue(this)
	}

	arguments_ := make([]goja.Value, len(arguments))
	for index, argument := range arguments {
		arguments_[index] = runtime.ToValue(argument)
	}

	return function(goja.FunctionCall{
		This:      this_,
		Arguments: arguments_,
	}).Export(), nil // can panic
}

func GetAndCall(runtime *goja.Runtime, object *goja.Object, name string, this any, arguments ...any) (value any, err error) {
	defer func() {
		if err_ := HandlePanic(recover()); err_ != nil {
			err = err_
		}
	}()

	if function := object.Get(name); function != nil { // Get can panic
		if function_, ok := function.Export().(JavaScriptFunc); ok {
			return Call(runtime, function_, this, arguments...)
		} else {
			return "", fmt.Errorf("\"%s\" is not a function", name)
		}
	} else {
		return "", fmt.Errorf("\"%s\" not found", name)
	}
}
