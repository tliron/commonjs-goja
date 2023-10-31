package commonjs

import (
	contextpkg "context"

	"github.com/dop251/goja"
	"github.com/tliron/exturl"
)

//
// Bind
//

type Bind interface {
	Unbind() (any, *Context, error)
}

//
// EarlyBind
//

type EarlyBind struct {
	Value   any
	Context *Context
}

// ([Bind] interface)
func (self EarlyBind) Unbind() (any, *Context, error) {
	return self.Value, self.Context, nil
}

// ([CreateExtensionFunc] signature)
func CreateEarlyBindExtension(context *Context) goja.Value {
	return context.Environment.Runtime.ToValue(func(id string, exportName string) (value goja.Value, err error) {
		defer func() {
			if err_ := HandlePanic(recover()); err_ != nil {
				err = err_
			}
		}()

		if url, err := context.Resolve(contextpkg.TODO(), id, false); err == nil {
			childEnvironment := context.Environment.NewChild()
			childContext := childEnvironment.NewContext(url, nil)
			if exports, err := childEnvironment.cachedRequire(url, childContext); err == nil {
				var value any

				if exportName == "" {
					value = exports.Export()
				} else {
					value = exports.Get(exportName).Export() // Get can panic
				}

				return context.Environment.Runtime.ToValue(EarlyBind{
					Value:   value,
					Context: childContext,
				}), nil
			} else {
				return nil, err
			}
		} else {
			return nil, err
		}
	})
}

//
// LateBind
//

type LateBind struct {
	URL        exturl.URL
	ExportName string
	Context    *Context
}

// ([Bind] interface)
func (self LateBind) Unbind() (value any, context *Context, err error) {
	defer func() {
		if err_ := HandlePanic(recover()); err_ != nil {
			err = err_
		}
	}()

	childEnvironment := self.Context.Environment.NewChild()
	childContext := childEnvironment.NewContext(self.URL, nil)
	if exports, err := childEnvironment.cachedRequire(self.URL, childContext); err == nil {
		var value any

		if self.ExportName == "" {
			value = exports.Export()
		} else {
			value = exports.Get(self.ExportName).Export() // Get can panic
		}

		return value, childContext, nil
	} else {
		return nil, nil, err
	}
}

// ([CreateExtensionFunc] signature)
func CreateLateBindExtension(context *Context) goja.Value {
	return context.Environment.Runtime.ToValue(func(id string, exportName string) (goja.Value, error) {
		if url, err := context.Resolve(contextpkg.TODO(), id, false); err == nil {
			return context.Environment.Runtime.ToValue(LateBind{
				URL:        url,
				ExportName: exportName,
				Context:    context,
			}), nil
		} else {
			return nil, err
		}
	})
}
