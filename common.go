package commonjs

import (
	"github.com/dop251/goja"
	"github.com/tliron/commonlog"
)

var log = commonlog.GetLogger("commonjs-goja")

// Call with a recover() value. If it's an error, then it will
// be returned. A [*goja.Exception] will be unwrapped and returned.
//
// Otherwise, will re-panic the value.
//
// This function is useful for cases in which Goja indicates errors
// by panicing instead of returning an error.
//
// Usage:
//
//	func MyFunc() (err error) {
//		defer func() {
//			if err_ := HandlePanic(recover()); err_ != nil {
//				err = err_
//			}
//		}()
//		// do something that can panic
//	}
func HandlePanic(r any) error {
	if r == nil {
		return nil
	}

	if exception, ok := r.(*goja.Exception); ok {
		if err := exception.Unwrap(); err != nil {
			return err
		} else {
			return exception
		}
	} else if err, ok := r.(error); ok {
		return err
	} else {
		panic(r)
	}
}
