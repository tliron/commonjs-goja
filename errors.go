package commonjs

import (
	"errors"

	"github.com/dop251/goja"
	"github.com/tliron/kutil/util"
)

func UnwrapJavaScriptException(err error) error {
	if exception, ok := err.(*goja.Exception); ok {
		switch exported := exception.Value().Export().(type) {
		case error:
			return UnwrapJavaScriptException(exported)

		case map[string]any:
			if value, ok := exported["value"]; ok {
				switch value_ := value.(type) {
				case error:
					return UnwrapJavaScriptException(value_)

				default:
					return errors.New(util.ToString(value))
				}
			}

		default:
			// This will work only if it's the *current* Go error in the runtime
			return UnwrapJavaScriptException(exception)
		}
	}

	return err
}

// Call with a recover() value. If it's an error, then it will
// be unwrapped and returned.
//
// Otherwise, will re-panic the value.
//
// This function is useful for cases in which Goja indicates errors
// by panicking instead of returning an error.
//
// Usage:
//
//	func MyFunc() (err error) {
//		defer func() {
//			if err_ := HandleJavaScriptPanic(recover()); err_ != nil {
//				err = err_
//			}
//		}()
//		// do something that can panic
//	}
func HandleJavaScriptPanic(r any) error {
	if r == nil {
		return nil
	}

	if err, ok := r.(error); ok {
		return UnwrapJavaScriptException(err)
	} else {
		panic(r)
	}
}
