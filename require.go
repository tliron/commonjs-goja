package commonjs

import (
	"github.com/dop251/goja"
)

func (self *Environment) NewRequire(jsContext *Context) *goja.Object {
	require := func(id string) (goja.Value, error) {
		context, cancelContext := self.NewTimeoutContext()
		defer cancelContext()

		return self.require(context, id, jsContext)
	}

	resolve := func(id string, options *goja.Object) (string, error) {
		// TODO: support resolve options

		context, cancelContext := self.NewTimeoutContext()
		defer cancelContext()

		if url, err := jsContext.Resolve(context, id, false); err == nil {
			return url.String(), nil
		} else {
			return "", err
		}
	}

	requireObject := self.Runtime.ToValue(require).(*goja.Object)

	requireObject.Set("cache", self.Modules)

	requireObject.Set("resolve", resolve)

	if jsContext.Parent != nil {
		requireObject.Set("main", jsContext.Parent.Module)
	} else {
		requireObject.Set("main", nil)
	}

	return requireObject
}
