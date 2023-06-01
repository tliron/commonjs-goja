package commonjs

import (
	"github.com/dop251/goja"
)

type JavaScriptFunc = func(goja.FunctionCall) goja.Value

func Call(runtime *goja.Runtime, function JavaScriptFunc, arguments ...any) any {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("%s", r)
		}
	}()

	arguments_ := make([]goja.Value, len(arguments))
	for index, argument := range arguments {
		arguments_[index] = runtime.ToValue(argument)
	}

	return function(goja.FunctionCall{
		This:      nil,
		Arguments: arguments_,
	}).Export()
}
