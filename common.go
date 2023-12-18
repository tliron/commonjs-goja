package commonjs

import (
	"github.com/dop251/goja"
	"github.com/tliron/commonlog"
)

var log = commonlog.GetLogger("commonjs")

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
