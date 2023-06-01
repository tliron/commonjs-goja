package commonjs

import (
	"github.com/dop251/goja"
	"github.com/tliron/exturl"
)

//
// Context
//

type Context struct {
	Environment *Environment
	Parent      *Context
	Module      *Module
	Resolve     ResolveFunc
	Extensions  []goja.Value
}

func (self *Environment) NewContext(url exturl.URL, parent *Context) *Context {
	context := Context{
		Environment: self,
		Parent:      parent,
		Module:      self.NewModule(),
	}

	if url != nil {
		self.AddModule(url, context.Module)
	}

	context.Resolve = self.CreateResolver(url, &context)

	for _, extension := range self.Extensions {
		context.Extensions = append(context.Extensions, extension.Create(&context))
	}

	// See: https://nodejs.org/api/modules.html#modules_the_module_object
	// See: https://nodejs.org/api/modules.html#modules_require_id

	context.Module.Require = self.Runtime.ToValue(func(id string) (goja.Value, error) {
		return self.requireId(id, &context)
	}).(*goja.Object)

	context.Module.Require.Set("cache", self.Modules)

	context.Module.Require.Set("resolve", func(id string, options *goja.Object) (string, error) {
		// TODO: options?
		if url, err := context.Resolve(id, false); err == nil {
			return url.String(), nil
		} else {
			return "", err
		}
	})

	if parent != nil {
		context.Module.Require.Set("main", parent.Module)
		parent.Module.Children = append(parent.Module.Children, context.Module)
	} else {
		context.Module.Require.Set("main", nil)
	}

	return &context
}
