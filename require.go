package commonjs

import (
	"github.com/dop251/goja"
)

func (self *Context) NewRequire() *goja.Object {
	require := func(id string) (*goja.Object, error) {
		context, cancelContext := self.Environment.NewTimeoutContext()
		defer cancelContext()

		if object, _, err := self.ResolveAndRequire(context, id, false, false, nil); err == nil {
			return object, nil
		} else {
			return nil, err
		}
	}

	resolve := func(id string, options *goja.Object) (string, error) {
		// TODO: support CommonJS resolve options

		context, cancelContext := self.Environment.NewTimeoutContext()
		defer cancelContext()

		if url, err := self.Resolve(context, id, false); err == nil {
			return url.String(), nil
		} else {
			return "", err
		}
	}

	requireObject := self.Environment.Runtime.ToValue(require).(*goja.Object)

	requireObject.Set("cache", self.Environment.Modules)

	requireObject.Set("resolve", resolve)

	if self.Parent != nil {
		requireObject.Set("main", self.Parent.Module)
	} else {
		requireObject.Set("main", nil)
	}

	return requireObject
}
