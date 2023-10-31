package commonjs

import (
	"fmt"

	"github.com/dop251/goja"
)

type JavaScriptFunc = func(goja.FunctionCall) goja.Value

type PanicHandlerFunc = func(r any)

func Call(panicHandler PanicHandlerFunc, runtime *goja.Runtime, function JavaScriptFunc, this any, arguments ...any) (value any, err error) {
	defer func() {
		if r := recover(); r != nil {
			if err_, ok := r.(error); ok {
				err = err_
			}
			if panicHandler != nil {
				panicHandler(r)
			}
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
	}).Export(), nil
}

func GetAndCall(panicHandler PanicHandlerFunc, runtime *goja.Runtime, object *goja.Object, name string, this any, arguments ...any) (value any, err error) {
	defer func() {
		if r := recover(); r != nil {
			if err_, ok := r.(error); ok {
				err = err_
			}
			if panicHandler != nil {
				panicHandler(r)
			}
		}
	}()

	if function := object.Get(name); function != nil {
		if function_, ok := function.Export().(JavaScriptFunc); ok {
			return Call(panicHandler, runtime, function_, this, arguments...)
		} else {
			return "", fmt.Errorf("\"%s\" is not a function", name)
		}
	} else {
		return "", fmt.Errorf("\"%s\" not found", name)
	}
}
